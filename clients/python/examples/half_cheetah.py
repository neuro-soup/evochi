import os
from typing import Callable, cast
import asyncio
from dataclasses import dataclass

import grpc.aio as grpc

from evochi.v1 import Worker, Eval

import numpy as np
import gymnasium as gym
import torch
from torch import nn
from jaxtyping import Float


class SimpleMLP(nn.Module):
    def __init__(self, obs_dim: int, act_dim: int, hidden_dim: int) -> None:
        super().__init__()
        self._net = nn.Sequential(
            nn.Linear(obs_dim, hidden_dim),
            nn.Tanh(),
            nn.Linear(hidden_dim, hidden_dim),
            nn.Tanh(),
            nn.Linear(hidden_dim, act_dim),
        )
        self.apply(self._init_weights)

    @staticmethod
    def _init_weights(m: nn.Module) -> None:
        if isinstance(m, nn.Linear):
            nn.init.xavier_uniform_(m.weight)
            nn.init.zeros_(m.bias)

    def forward(self, obs: torch.Tensor) -> torch.Tensor:
        return self._net(obs)


type SamplingStrategy = Callable[[int, int, torch.Generator], torch.Tensor]
type RewardTransform = Callable[
    [Float[torch.Tensor, "n_pop"]],
    Float[torch.Tensor, "n_pop"],
]


@dataclass
class State:
    seed: int
    learning_rate: float
    noise_std: float
    weight_decay: float
    params: torch.Tensor


class HalfCheetah(Worker[State]):
    def __init__(
        self,
        channel: grpc.Channel,
        cores: int,
        vectorization: gym.VectorizeMode = gym.VectorizeMode.SYNC,
        device: torch.device = torch.device("cpu"),
    ) -> None:
        super().__init__(
            channel=channel,
            cores=cores,
            initialize=lambda _: self._initialize(),
            evaluate=lambda _, __, slices: self._evaluate(slices),
            optimize=lambda _, epoch, rewards: self._optimize(epoch, rewards),
        )
        self._device = device
        self._env = gym.make(
            "HalfCheetah-v5",
            num_envs=cores,  # TODO: separate variable?
            vectorization_mode=vectorization,
        )
        self._rng: torch.Generator | None = None
        self._noise: torch.Tensor | None = None
        self._mlp: SimpleMLP | None = None
        self._obs: torch.Tensor = self._env.reset()[0]

    def _initialize(self) -> State:
        """Initialize shared information across all workers."""
        self._mlp = self._empty_mlp()
        self._generate_noise()
        # TODO: make configurable
        return State(
            seed=np.random.randint(0, 2**32),
            learning_rate=0.1,
            noise_std=0.1,
            weight_decay=0.1,
            params=nn.utils.parameters_to_vector(self._mlp.parameters()),
        )

    def _prepare(self, slices: list[slice]) -> None:
        """Prepares this worker for evaluation."""
        if self._rng is None:
            self._rng = torch.Generator(device=self._device).manual_seed(
                self.state.seed
            )
        if self._mlp is None:
            self._mlp = self._mlp_from_params(slices)

    def _evaluate(self, slices: list[slice]) -> list[Eval]:
        """Evaluates the model by computing the rewards for each slice."""
        self._prepare(slices)

        raise NotImplementedError("not implemented")  # TODO: implement

        actions = self._choose_actions()
        self._obs, rewards, terminations, truncations, _ = self._env.step(actions)
        return []

    def _optimize(self, _: int, raw_rewards: list[float]) -> State:
        rewards = self._transform_reward(np.array(raw_rewards))

        raise NotImplementedError("not implemented")  # TODO: implement

        self._generate_noise()
        return self.state

    def _choose_actions(self) -> np.ndarray:
        assert self._mlp is not None, "Cannot determine action without model"
        # TODO: use old implementation here
        return self._mlp(self._obs).argmax(dim=1).detach().numpy()

    def _mlp_from_params(self, slices: list[slice]) -> SimpleMLP:
        raise NotImplementedError("not implemented")  # TODO: implement
        # mlp = self._empty_mlp()
        # # nn.utils.vector_to_parameters(
        # #     vec=self._perturb(self.state.params),
        # #     parameters=mlp.parameters(),
        # # )
        # return mlp.to(self._device)

    def _generate_noise(self) -> None:
        assert self._mlp is not None, "Cannot determine noise without model"
        n_params = nn.utils.parameters_to_vector(self._mlp.parameters()).numel()
        self._noise = torch.randn((self.population_size, n_params), generator=self._rng)

    def _transform_reward(self, rewards: np.ndarray) -> np.ndarray:
        return (rewards - rewards.mean()) / rewards.std()

    def _empty_mlp(self) -> SimpleMLP:
        return SimpleMLP(obs_dim=4, act_dim=1, hidden_dim=32).to(self._device)


async def main() -> None:
    channel = grpc.insecure_channel("localhost:8080")
    worker = HalfCheetah(channel=channel, cores=cast(int, os.cpu_count()))
    await worker.start()


if __name__ == "__main__":
    asyncio.run(main())
