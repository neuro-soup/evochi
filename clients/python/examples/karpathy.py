"""
A bare bones examples of optimizing a black-box function (f) using
Natural Evolution Strategies (NES), where the parameter distribution is a
gaussian of fixed standard deviation.

Originally from https://gist.github.com/karpathy/77fbb6a8dac5395f1b73e7a89300318d
"""

from typing import cast
import asyncio
from dataclasses import dataclass
import os
import time

from evochi.v1 import Worker, Eval

import grpc.aio as grpc
import numpy as np

np.random.seed(0)


# the function we want to optimize
def f(w) -> np.ndarray:
    # here we would normally:
    # ... 1) create a neural network with weights w
    # ... 2) run the neural network on the environment for some time
    # ... 3) sum up and return the total reward

    # but for the purposes of an example, lets try to minimize
    # the L2 distance to a specific solution vector. So the highest reward
    # we can achieve is 0, when the vector w is exactly equal to solution
    reward = -np.sum(np.square(solution - w))
    # time.sleep(0.1)
    return reward


sigma = 0.1
alpha = 0.001

# start the optimization
solution = np.array([0.5, 0.1, -0.3])


@dataclass
class State:
    weights: np.ndarray
    N: np.ndarray


def initialize(w: Worker[State]) -> State:
    return State(weights=np.random.randn(3), N=np.random.randn(w.pop_size, 3))


def evaluate(worker: Worker[State], _: int, slices: list[slice]) -> list[Eval]:
    N = worker.state.N
    evals: list[Eval] = []
    for sl in slices:
        rewards: list[float] = []
        for j in range(sl.start, sl.stop):
            w_try = worker.state.weights + sigma * N[j]
            rewards.append(f(w_try).item())
        evals.append(Eval(slice=sl, rewards=rewards))
    return evals


def optimize(worker: Worker[State], epoch: int, rewards: list[float]) -> State:
    R = np.array(rewards)
    A = (R - np.mean(R)) / np.std(R)
    N = worker.state.N
    npop = worker.pop_size
    w = worker.state.weights
    w = w + alpha / (npop * sigma) * np.dot(N.T, A)
    if epoch % 20 == 0:
        print("epoch %d. w: %s, reward: %f" % (epoch, str(w), f(w)))
    return State(weights=w, N=np.random.randn(npop, 3))


async def main() -> None:
    channel = grpc.insecure_channel("localhost:8080")
    worker = Worker(
        channel=channel,
        cores=cast(int, os.cpu_count()),
        evaluate=evaluate,
        initialize=initialize,
        optimize=optimize,
    )
    await worker.start()


if __name__ == "__main__":
    asyncio.run(main())
