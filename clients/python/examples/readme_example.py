from dataclasses import dataclass
import random

import grpc.aio as grpc

import evochi.v1 as evochi


@dataclass
class State:
    # Shared information across all workers is stored here. The state is centrally
    # and synchronously updated at the end of each epoch.

    # The state might also contain configuration options for other workers, such
    # as `seed`, `learning_rate`, etc.

    # IMPORTANT: The state must be serializable using `pickle`. The state is only
    # sent once per epoch by a single worker and whenever a new worker joins
    # the training. However, the received state must be loaded using `pickle`
    # on each worker, which means that data structures, such as `torch.Tensor`,
    # that are stored on a GPU device must be moved to the CPU when initializing
    # and optimizing the state so that non-GPU workers can deserialize the
    # state.
    seed: float


class AwesomeWorker(evochi.Worker[State]):
    def __init__(self, channel: grpc.Channel, cores: int) -> None:
        super().__init__(channel, cores)

    def initialize(self) -> State:
        # This method is called on the first worker to join the training. Since
        # the server doesn't know anything about the state of the workers, the
        # first worker is responsible for initializing the state, which is then
        # broadcasted to all subsequent workers.
        # TODO: initialize state parameters of the model
        return State(seed=42)

    def evaluate(self, epoch: int, slices: list[slice]) -> list[evochi.Eval]:
        # This method is called whenever the server requests an evaluation step
        # for the current worker. The given slices represent the index ranges of
        # the population to be evaluated.
        #
        # Here, you can perform (an arbitrary number of) environment steps or
        # whatever you want to evaluate. It is possible that during a single
        # epoch, the same worker receives multiple evaluation requests.
        #
        # Note that the length of the slice (stop-start) must be equal to the
        # number of rewards in a single `evochi.Eval` object.
        # TODO: implement a proper evaluation step
        return [
            evochi.Eval(
                slice=slice,
                rewards=[random.randint(-42, 42) for _ in range(slice.start, slice.stop)],
            )
            for slice in slices
        ]

    def optimize(self, epoch: int, rewards: list[float]) -> State:
        # This method is called at the end of each epoch. The accumulated rewards
        # of the total population are sent to all workers to perform an optimization
        # step, which is performed in this method.
        #
        # It makes sense that the workers' states must be equal, which is ensured
        # using a `seed` in the state. After the optimization step, a worker
        # is requested to send its state to the server, which is then used for
        # new workers to join the training.
        # TODO: update state parameters of the model
        return State(seed=self.state.seed)


async def main() -> None:
    # Create a gRPC channel to the server. Here, the evochi server is assumed to
    # be running on localhost:8080.
    channel = grpc.insecure_channel("localhost:8080")

    # The number of cores determines the max length of slices (stop-start) that
    # the server will send to the worker to evaluate. Of course, this must not
    # necessarily be equal to the number of cores in CPU/GPU.
    worker = AwesomeWorker(channel=channel, cores=5)

    await worker.start()
