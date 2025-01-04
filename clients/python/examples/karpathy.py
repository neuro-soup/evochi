"""
A bare bones examples of optimizing a black-box function (f) using
Natural Evolution Strategies (NES), where the parameter distribution is a
gaussian of fixed standard deviation.

Originally from https://gist.github.com/karpathy/77fbb6a8dac5395f1b73e7a89300318d
"""

import asyncio
import os
from dataclasses import dataclass
from typing import cast

import evochi.v1 as evochi
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
    return reward


sigma = 0.1
alpha = 0.001

# start the optimization
solution = np.array([0.5, 0.1, -0.3])


@dataclass
class State:
    weights: np.ndarray
    N: np.ndarray


class Worker(evochi.Worker[State]):
    def initialize(self) -> State:
        return State(weights=np.random.randn(3), N=np.random.randn(self.population_size, 3))

    def evaluate(self, epoch: int, slices: list[slice]) -> list[evochi.Eval]:
        N = self.state.N
        evals: list[evochi.Eval] = []
        for sl in slices:
            rewards: list[float] = []
            for j in range(sl.start, sl.stop):
                w_try = self.state.weights + sigma * N[j]
                rewards.append(f(w_try).item())
            evals.append(evochi.Eval(slice=sl, rewards=rewards))
        return evals

    def optimize(self, epoch: int, rewards: list[float]) -> State:
        R = np.array(rewards)
        A = (R - np.mean(R)) / np.std(R)
        N = self.state.N
        npop = self.population_size
        w = self.state.weights
        w = w + alpha / (npop * sigma) * np.dot(N.T, A)
        if epoch % 20 == 0:
            print("epoch %d. w: %s, reward: %f" % (epoch, str(w), f(w)))
        return State(weights=w, N=np.random.randn(npop, 3))


async def main() -> None:
    channel = grpc.insecure_channel("localhost:8080")
    worker = Worker(channel=channel, cores=cast(int, os.cpu_count()))
    await worker.start()


if __name__ == "__main__":
    asyncio.run(main())
