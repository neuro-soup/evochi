from typing import Any, AsyncIterable, Callable, NamedTuple

import logging
import asyncio
import pickle

from grpc import StatusCode
import grpc.aio as grpc
import zstandard as zstd

import evochi.v1.evochi_pb2 as v1
import evochi.v1.evochi_pb2_grpc as client


class Eval(NamedTuple):
    slice: slice
    rewards: list[float]


type EvalFn[S] = Callable[[Worker[S], int, list[slice]], list[Eval]]
type InitFn[S] = Callable[[Worker[S]], Any]
type OptimizeFn[S] = Callable[[Worker[S], int, list[float]], Any]
type StopFn[S] = Callable[[Worker[S]], bool]


class Worker[S]:
    def __init__(
        self,
        channel: grpc.Channel,
        cores: int,
        evaluate: EvalFn[S],
        initialize: InitFn[S],
        optimize: OptimizeFn[S],
        stop: StopFn[S] | None = None,
    ) -> None:
        """Initializes the client to interact with the server via the given channel.

        Args:
            channel: The gRPC channel to use for communication.
            cores: The number of cores to use for evaluation.
            evaluate: A function that evaluates slices of the population upon request.
            initialize: A function that initializes the state upon request.
            optimize: A function that makes optimization steps upon request.
            stop: A function that is called when the worker should stop.
        """
        self._channel = channel
        self._cores = cores
        self._evaluate = evaluate
        self._initialize = initialize
        self._optimize = optimize
        self._stop = stop
        self._stub = client.EvochiServiceStub(channel)
        self._heartbeat_seq_id = 0
        self._closed = False
        self._token: str | None = None
        self._pop_size: int | None = None
        self._heartbeat_interval: int = 0
        self._current_state: S | None = None

    @property
    def cores(self) -> int:
        return self._cores

    @property
    def pop_size(self) -> int:
        if self._pop_size is None:
            raise RuntimeError("Worker has not been initialized yet")
        return self._pop_size

    @property
    def state(self) -> S:
        if self._current_state is None:
            raise RuntimeError("Worker has not been initialized yet")
        return self._current_state

    async def close(self) -> None:
        """Closes the client's channel."""
        if self._closed:
            raise RuntimeError("Worker is already closed")
        logging.debug("Closing worker")
        self._closed = True
        await self._channel.close()

    async def start(self) -> None:
        """Starts the worker."""
        logging.debug(
            "Starting worker with %d cores, waiting to be ready...", self._cores
        )
        await self._channel.channel_ready()
        logging.debug("Worker ready, starting...")
        await self._handle_events()

    async def _keep_alive(self) -> None:
        """Sends a heartbeat to the server periodically."""
        while not self._closed:
            self._heartbeat_seq_id += 1
            await self._heartbeat(v1.HeartbeatRequest(seq_id=self._heartbeat_seq_id))
            await asyncio.sleep(self._heartbeat_interval)

    async def _handle_events(self) -> None:
        """Handles events from the server."""
        iter = self._subscribe(v1.SubscribeRequest(cores=self._cores))
        try:
            async for event in iter:
                match event.type:
                    case v1.EVENT_TYPE_HELLO:
                        self._handle_hello_event(event.hello)
                    case v1.EVENT_TYPE_INITIALIZE:
                        await self._handle_init_event(event.initialize)
                    case v1.EVENT_TYPE_EVALUATE:
                        await self._handle_eval_event(event.evaluate)
                    case v1.EVENT_TYPE_SHARE_STATE:
                        await self._handle_share_state_event(event.share_state)
                    case v1.EVENT_TYPE_OPTIMIZE:
                        await self._handle_optimize_event(event.optimize)
                    case v1.EVENT_TYPE_STOP:
                        await self._handle_stop_event(event.stop)
                    case _:
                        logging.warning("Received unknown event type %s", event.type)
        except grpc.AioRpcError as e:
            if e.code() == StatusCode.CANCELLED:
                logging.debug("Worker cancelled")
            else:
                raise e

    def _handle_hello_event(self, event: v1.HelloEvent) -> None:
        """Handles a hello event from the server."""
        logging.debug(
            "Received hello event with id %s and token %s", event.id, event.token
        )
        self._token = event.token
        self._heartbeat_interval = event.heartbeat_interval
        self._pop_size = event.population_size
        self._current_state = (
            self._decompress_state(event.state) if event.state else None
        )
        asyncio.create_task(self._keep_alive())  # TODO: is this correct?

    async def _handle_init_event(self, event: v1.InitializeEvent) -> None:
        """Handles an init event from the server."""
        logging.debug("Received init event with task id %s", event.task_id)
        state = self._initialize(self)
        self._current_state = state
        await self._finish_initialization(
            v1.FinishInitializationRequest(
                task_id=event.task_id,
                state=self._compressed_state(),
            )
        )

    async def _handle_eval_event(self, event: v1.EvaluateEvent) -> None:
        """Handles an eval event from the server."""
        logging.debug("Received eval event with task id %s", event.task_id)
        evals = self._evaluate(
            self,
            event.epoch,
            [slice(sl.start, sl.end) for sl in event.slices],
        )
        await self._finish_evaluation(
            v1.FinishEvaluationRequest(
                task_id=event.task_id,
                evaluations=[
                    v1.Evaluation(slice=sl, rewards=eval.rewards)
                    for sl, eval in zip(event.slices, evals)
                ],
            )
        )

    async def _handle_optimize_event(self, event: v1.OptimizeEvent) -> None:
        """Handles an optimize event from the server."""
        logging.debug("Received optimize event with task id %s", event.task_id)
        optimized = self._optimize(self, event.epoch, list(event.rewards))
        self._current_state = optimized
        await self._finish_optimization(
            v1.FinishOptimizationRequest(task_id=event.task_id)
        )

    async def _handle_share_state_event(self, event: v1.ShareStateEvent) -> None:
        """Handles a share state event from the server."""
        logging.debug("Received share state event with task id %s", event.task_id)
        assert self._current_state is not None
        await self._finish_share_state(
            v1.FinishShareStateRequest(
                task_id=event.task_id,
                state=self._compressed_state(),
            )
        )

    async def _handle_stop_event(self, event: v1.StopEvent) -> None:
        """Handles a stop event from the server."""
        logging.debug("Received stop event with task id %s", event.task_id)
        if self._stop is not None:
            self._stop(self)

    def _subscribe(
        self,
        request: v1.SubscribeRequest,
    ) -> AsyncIterable[v1.SubscribeResponse]:
        """Sends a subscribe request to the server and returns an iterable of responses.

        Args:
            request: The subscribe request to send to the server.

        Returns:
            An iterable of subscribe responses received from the server.
        """
        logging.debug("Subscribing to events")
        return self._stub.Subscribe(request)

    async def _heartbeat(self, request: v1.HeartbeatRequest) -> v1.HeartbeatResponse:
        """Sends a heartbeat request to the server and returns the response.

        Args:
            request: The heartbeat request to send to the server.

        Returns:
            The heartbeat response received from the server.
        """
        logging.debug("Sending heartbeat")
        return await self._stub.Heartbeat(request, metadata=self._metadata())

    async def _finish_evaluation(
        self,
        request: v1.FinishEvaluationRequest,
    ) -> v1.FinishEvaluationResponse:
        """Sends a finish evaluation request to the server and returns the response.

        Args:
            request: The finish evaluation request to send to the server.

        Returns:
            The finish evaluation response received from the server.
        """
        logging.debug("Sending finish evaluation")
        return await self._stub.FinishEvaluation(request, metadata=self._metadata())

    async def _finish_optimization(
        self,
        request: v1.FinishOptimizationRequest,
    ) -> v1.FinishOptimizationResponse:
        """Sends a finish optimization request to the server and returns the response.

        Args:
            request: The finish optimization request to send to the server.

        Returns:
            The finish optimization response received from the server.
        """
        logging.debug("Sending finish optimization")
        return await self._stub.FinishOptimization(request, metadata=self._metadata())

    async def _finish_initialization(
        self,
        request: v1.FinishInitializationRequest,
    ) -> v1.FinishInitializationResponse:
        """Sends a finish initialization request to the server and returns the response.

        Args:
            request: The finish initialization request to send to the server.

        Returns:
            The finish initialization response received from the server.
        """
        logging.debug("Sending finish initialization")
        return await self._stub.FinishInitialization(request, metadata=self._metadata())

    async def _finish_share_state(
        self,
        request: v1.FinishShareStateRequest,
    ) -> v1.FinishShareStateResponse:
        """Sends a finish share state request to the server and returns the response.

        Args:
            request: The finish share state request to send to the server.

        Returns:
            The finish share state response received from the server.
        """
        logging.debug("Sending finish share state")
        return await self._stub.FinishShareState(request, metadata=self._metadata())

    def _compressed_state(self) -> bytes:
        """Returns the current state compressed using zstandard."""
        return zstd.compress(pickle.dumps(self._current_state))

    def _decompress_state(self, state: bytes) -> S:
        """Returns the current state decompressed using zstandard."""
        return pickle.loads(zstd.decompress(state))

    def _metadata(self) -> list[tuple[str, str]]:
        """Returns the metadata to use for requests to the server."""
        if self._token is None:
            raise RuntimeError("Client does not have a token yet")
        return [("authorization", f"Bearer {self._token}")]


__all__ = ["Worker"]
