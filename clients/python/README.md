# Evochi Python Client

<div align="center">
    <img src="https://img.shields.io/badge/Written_In-Python-f7d44f?style=for-the-badge&logo=python" alt="Python" />
</div>

<br />

This project provides an easy-to-use Python client for the Evochi API such that
users don't have to worry about the low-level networking details of the API protocol.

## Installation

Install the package from PyPI:

```bash
pip install evochi
```

Or, if you want to install from (Git) source:

```bash
pip install "evochi @ git+https://github.com/neuro-soup/evochi.git/#subdirectory=clients/python"
```

## Basic Usage

### Minimal Example

The following example is a minimal example of how to use the Evochi Python
client. It is assumed that the Evochi server is running on `localhost:8080`.

```py
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
        # The total number of individuals across all slices will be <= self.cores.
        #
        # Note that the length of the slice (stop-start) must be equal to the
        # number of rewards in a single `evochi.Eval` object.
        # TODO: implement a proper evaluation step
        return [
            evochi.Eval(
                slice=slice,
                rewards=[
                    random.randint(-42, 42)
                    for _ in range(slice.start, slice.stop)
                ],
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
```

## Understanding Slice Distribution

When workers connect to the evochi server, they specify their `cores` parameter, which represents their parallel evaluation capacity. The server uses this information to distribute work efficiently:

### How Slices Work

- **Slices** represent contiguous segments of the population (e.g., `slice(0, 5)` means individuals 0-4)
- Each worker receives one or more slices per evaluation request
- The total number of individuals across all slices will be `<= cores`

### Example

If a worker has 8 cores and the unassigned population segments are fragmented:
```python
# Worker might receive:
slices = [slice(0, 3), slice(7, 12)]  # Total: 3 + 5 = 8 individuals

# Worker evaluates all 8 individuals in parallel and returns:
[
    Eval(slice(0, 3), [reward0, reward1, reward2]),
    Eval(slice(7, 12), [reward7, reward8, reward9, reward10, reward11])
]
```

### Dynamic Work Distribution

- New workers joining mid-epoch immediately receive available work from the unassigned pool
- If a worker disconnects, its unfinished slices are redistributed to remaining workers
- No worker waits idle if there's unassigned work available

## API Reference

### Worker Class

The `Worker` class is the base class for all evochi workers. It handles the communication protocol with the server and provides lifecycle hooks for your implementation.

#### Properties

- `cores` → `int`: Maximum number of individuals to evaluate in a single evaluate call. Usually the number of cores/gpus/... this worker uses for parallel evaluation.
- `population_size` → `int`: Total population size (available after initialization)
- `max_epochs` → `int`: Maximum number of epochs (available after initialization)
- `state` → `S`: Current shared state (available after initialization)

#### Lifecycle Methods (to override)

- `initialize() → S`: Called on the first worker to initialize shared state
- `evaluate(epoch: int, slices: list[slice]) → list[Eval]`: Called to evaluate population slices
- `optimize(epoch: int, rewards: list[float]) → S`: Called to perform optimization with full population rewards
- `on_stop(cancel: bool) → None`: Called when worker is requested to stop (optional)
- `on_state_change(state: S) → None`: Called when shared state changes (optional)

#### Public Methods

- `async start() → None`: Start the worker and begin processing tasks
- `async close() → None`: Close the worker's connection

### Eval Class

Represents evaluated slices with their rewards.

- `slice`: The population slice that was evaluated
- `rewards`: List of rewards for each individual in the slice
- `from_flat(slices: list[slice], rewards: list[float]) → list[Eval]`: Helper to construct Eval objects from flat reward list
