from __future__ import annotations

from typing import AsyncIterable, NamedTuple
from abc import ABC, abstractmethod

import asyncio
import logging
import pickle
from datetime import datetime

import grpc.aio as grpc
from google.protobuf import timestamp_pb2
from grpc import StatusCode
import zstandard as zstd

import evochi.v1.evochi_pb2 as v1
import evochi.v1.evochi_pb2_grpc as client


class Eval(NamedTuple):
    """The evaluated slice and its rewards.

    Note that the len(rewards) must be equal to the length of the slice, which
    is `slice.stop - slice.start`.
    """

    slice: slice
    rewards: list[float]

    @staticmethod
    def from_flat(slices: list[slice], rewards: list[float]) -> list[Eval]:
        """Constructs a list of Eval objects from a list of slices and a list of
        rewards."""
        total_width = sum(sl.stop - sl.start for sl in slices)
        assert total_width == len(rewards), "Number of rewards must match the total width of the slices"
        evals: list[Eval] = []
        offset = 0
        for sl in slices:
            evals.append(
                Eval(
                    slice=slice(sl.start, sl.stop),
                    rewards=rewards[offset : offset + (sl.stop - sl.start)],
                )
            )
            offset += sl.stop - sl.start
        return evals


class Worker[S](ABC):
    def __init__(self, channel: grpc.Channel, cores: int) -> None:
        """Initializes the client to interact with the server via the given channel.

        Args:
            channel: The gRPC channel to use for communication.
            cores: The maximum number of cores to use for evaluation.
        """
        self._channel = channel
        self._cores = cores
        self._stub = client.EvochiServiceStub(channel)
        self._heartbeat_seq_id = 0
        self._closed = False
        self._token: str | None = None
        self._id: str | None = None
        self._pop_size: int | None = None
        self._max_epochs: int | None = None
        self._heartbeat_interval: int = 0
        self._current_state: S | None = None

    @property
    def cores(self) -> int:
        """Returns the number of cores the worker is using."""
        return self._cores

    @property
    def id(self) -> str:
        if self._id is None:
            raise RuntimeError("Worker has not been initialized yet")
        return self._id

    @property
    def population_size(self) -> int:
        if self._pop_size is None:
            raise RuntimeError("Worker has not been initialized yet")
        return self._pop_size

    @property
    def max_epochs(self) -> int:
        if self._max_epochs is None:
            raise RuntimeError("Worker has not been initialized yet")
        return self._max_epochs

    @property
    def state(self) -> S:
        if self._current_state is None:
            raise RuntimeError("Worker has not been initialized yet")
        return self._current_state

    @abstractmethod
    def initialize(self) -> S:
        """Lifecycle hook that is called for the first worker to initialize the state."""

    @abstractmethod
    def evaluate(self, epoch: int, slices: list[slice]) -> list[Eval]:
        """Lifecycle hook that is called when the worker should perform an evaluation step.

        Args:
            epoch: The current epoch / generation.
            slices: List of population slices to evaluate. Usually a single slice.
                The total number of individuals across all slices will be <= self.cores.

        Returns:
            List of Eval objects, one per slice, containing the rewards for each individual in that slice.
            The number of rewards in each Eval must exactly match the slice size (slice.stop - slice.start).
        """

    @abstractmethod
    def optimize(self, epoch: int, rewards: list[float]) -> S:
        """Lifecycle hook that is called when the worker should perform an optimization step."""

    def on_stop(self, cancel: bool) -> None:
        """Hook that is called when the worker is requested to stop.

        Args:
            cancel: Whether the worker's connection has been cancelled.
        """

    def on_state_change(self, state: S) -> None:
        """Hook that is called when the worker's state changes."""

    async def close(self) -> None:
        """Closes the client's channel."""
        if self._closed:
            raise RuntimeError("Worker is already closed")
        logging.debug("Closing worker")
        self._closed = True
        await self._channel.close()

    async def start(self) -> None:
        logging.debug("Starting worker with %d cores, waiting to be ready...", self._cores)
        await self._channel.channel_ready()
        logging.debug("Worker ready, starting...")
        await self._handle_events()

    async def _keep_alive(self) -> None:
        """Sends a heartbeat to the server periodically."""
        while not self._closed:
            self._heartbeat_seq_id += 1
            timestamp = timestamp_pb2.Timestamp()
            timestamp.FromDatetime(datetime.now())
            await self._heartbeat(
                v1.HeartbeatRequest(
                    seq_id=self._heartbeat_seq_id,
                    timestamp=timestamp,
                )
            )
            await asyncio.sleep(self._heartbeat_interval)

    async def _handle_events(self) -> None:
        """Event loop that listens for events from the server and dispatches them to the appropriate handlers."""
        iterator = self._subscribe(v1.SubscribeRequest(cores=self._cores))
        try:
            async for event in iterator:
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
                self.on_stop(True)
            else:
                raise e

    def _handle_hello_event(self, event: v1.HelloEvent) -> None:
        logging.debug(
            "Received hello event with id %s and token %s",
            event.id,
            event.token,
        )
        self._id = event.id
        self._token = event.token
        self._heartbeat_interval = event.heartbeat_interval
        self._pop_size = event.population_size
        self._max_epochs = event.max_epochs

        self._current_state = self._decompress_state(event.state) if event.state else None
        if self._current_state is not None:
            self.on_state_change(self._current_state)

        asyncio.create_task(self._keep_alive())

    async def _handle_init_event(self, event: v1.InitializeEvent) -> None:
        logging.debug("Received init event with task id %s", event.task_id)
        state = self.initialize()

        self._current_state = state
        self.on_state_change(self._current_state)

        await self._finish_initialization(
            v1.FinishInitializationRequest(
                task_id=event.task_id,
                state=self._compressed_state(),
            )
        )

    async def _handle_eval_event(self, event: v1.EvaluateEvent) -> None:
        logging.debug("Received eval event with task id %s", event.task_id)
        evals = self.evaluate(event.epoch, [slice(sl.start, sl.end) for sl in event.slices])
        await self._finish_evaluation(
            v1.FinishEvaluationRequest(
                task_id=event.task_id,
                evaluations=[
                    v1.Evaluation(
                        slice=v1.Slice(start=e.slice.start, end=e.slice.stop),
                        rewards=e.rewards,
                    )
                    for e in evals
                ],
            )
        )

    async def _handle_optimize_event(self, event: v1.OptimizeEvent) -> None:
        logging.debug("Received optimize event with task id %s", event.task_id)
        optimized = self.optimize(event.epoch, list(event.rewards))

        self._current_state = optimized
        self.on_state_change(optimized)

        await self._finish_optimization(v1.FinishOptimizationRequest(task_id=event.task_id))

    async def _handle_share_state_event(self, event: v1.ShareStateEvent) -> None:
        logging.debug("Received share state event with task id %s", event.task_id)
        assert self._current_state is not None
        await self._finish_share_state(
            v1.FinishShareStateRequest(
                task_id=event.task_id,
                state=self._compressed_state(),
            )
        )

    async def _handle_stop_event(self, event: v1.StopEvent) -> None:
        logging.debug("Received stop event with task id %s", event.task_id)
        self.on_stop(False)

    def _subscribe(self, request: v1.SubscribeRequest) -> AsyncIterable[v1.SubscribeResponse]:
        logging.debug("Subscribing to events")
        return self._stub.Subscribe(request)

    async def _heartbeat(self, request: v1.HeartbeatRequest) -> v1.HeartbeatResponse:
        logging.debug("Sending heartbeat")
        return await self._stub.Heartbeat(request, metadata=self._metadata())

    async def _finish_evaluation(self, request: v1.FinishEvaluationRequest) -> v1.FinishEvaluationResponse:
        logging.debug("Sending finish evaluation")
        return await self._stub.FinishEvaluation(request, metadata=self._metadata())

    async def _finish_optimization(self, request: v1.FinishOptimizationRequest) -> v1.FinishOptimizationResponse:
        logging.debug("Sending finish optimization")
        return await self._stub.FinishOptimization(request, metadata=self._metadata())

    async def _finish_initialization(self, request: v1.FinishInitializationRequest) -> v1.FinishInitializationResponse:
        logging.debug("Sending finish initialization")
        return await self._stub.FinishInitialization(request, metadata=self._metadata())

    async def _finish_share_state(self, request: v1.FinishShareStateRequest) -> v1.FinishShareStateResponse:
        logging.debug("Sending finish share state")
        return await self._stub.FinishShareState(request, metadata=self._metadata())

    def _compressed_state(self) -> bytes:
        """Returns the current state compressed using zstandard."""
        return zstd.compress(pickle.dumps(self._current_state))

    @staticmethod
    def _decompress_state(state: bytes) -> S:
        """Returns the current state decompressed using zstandard."""
        return pickle.loads(zstd.decompress(state))

    def _metadata(self) -> list[tuple[str, str]]:
        """Returns the metadata to use for requests to the server."""
        if self._token is None:
            raise RuntimeError("Client does not have a token yet")
        return [("authorization", f"Bearer {self._token}")]


__all__ = ["Worker", "Eval"]
