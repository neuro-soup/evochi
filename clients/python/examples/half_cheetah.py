import asyncio
from dataclasses import dataclass

import grpc.aio as grpc

from evochi.v1 import Worker, Eval

import numpy as np
import gymnasium as gym
import torch
from torch import nn


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


@dataclass
class State:
    seed: int
    max_steps: int
    learning_rate: float
    noise_std: float
    params: torch.Tensor


class HalfCheetah(Worker[State]):
    def __init__(
        self,
        channel: grpc.Channel,
        cores: int,
        device: torch.device = torch.device("cpu"),
        vectorization_mode: gym.VectorizeMode = gym.VectorizeMode.SYNC,
        render: bool = False,
    ) -> None:
        super().__init__(
            channel=channel,
            cores=cores,
            initialize=lambda _: self._handle_init(),
            evaluate=lambda _, epoch, slices: self._handle_eval(epoch, slices),
            optimize=lambda _, epoch, rewards: self._handle_optim(epoch, rewards),
        )
        self._device = device
        self._vectorization_mode = vectorization_mode
        self._render = render
        self._noise: torch.Tensor | None = None
        self._mlp: SimpleMLP = SimpleMLP(obs_dim=17, act_dim=6, hidden_dim=64).to(
            self._device
        )

    def _handle_init(self) -> State:
        """Initialize shared information across all workers."""
        seed = np.random.randint(0, 2**32)
        return State(
            seed=seed,
            max_steps=1_000,
            learning_rate=0.1,
            noise_std=0.1,
            params=nn.utils.parameters_to_vector(self._mlp.parameters()),
        )

    @torch.inference_mode()
    def _handle_eval(self, epoch: int, slices: list[slice]) -> list[Eval]:
        """Evaluates the model by computing the rewards for each slice."""
        if self._mlp is None:
            self._mlp = self._construct_mlp(slices)
        if self._noise is None:
            self._generate_noise()

        total_width = sum(sl.stop - sl.start for sl in slices)
        env = gym.make_vec(
            "HalfCheetah-v5",
            num_envs=total_width,
            vectorization_mode=self._vectorization_mode,
            render_mode="human" if self._render else None,
        )
        obs, _ = env.reset(seed=self.state.seed)

        dones = np.zeros(total_width, dtype=bool)
        returns = np.zeros(total_width, dtype=np.float32)

        for _ in range(self.state.max_steps):
            actions = self._choose_actions(obs)
            obs, reward, terminations, truncations, _ = env.step(actions)
            dones = dones | terminations | truncations
            returns += reward * ~dones
            if dones.all():
                break

        acc_rewards: list[tuple[slice, list[float]]] = []
        offset = 0
        for sl in slices:
            rewards = returns[offset : offset + (sl.stop - sl.start)].tolist()
            acc_rewards.append((sl, rewards))
            offset += sl.stop - sl.start

        return [
            Eval(
                slice=slice,
                rewards=rewards,
            )
            for slice, rewards in acc_rewards
        ]

    def _handle_optim(self, epoch: int, raw_rewards: list[float]) -> State:
        """Perform an optimization step."""
        assert self._mlp is not None, "Cannot optimize without model"

        lr = self.state.learning_rate
        noise_std = self.state.noise_std
        epsilon = self._noise
        pop_size = self.population_size

        rewards = self._transform_reward(torch.tensor(raw_rewards, device=self._device))

        rewards = rewards.unsqueeze(dim=1)

        params = nn.utils.parameters_to_vector(self._mlp.parameters())
        epsilon = self._noise
        pert = torch.sum(epsilon * rewards, dim=0)

        # Update
        self.state.params = params + lr / (pop_size * noise_std) * pert
        self._mlp = self._construct_mlp([])

        # reset worker state
        self._noise = None
        self._rng = None

        self._generate_noise()
        return self.state

    def _choose_actions(self, obs: np.ndarray) -> np.ndarray:
        """Choose actions based on the given observation."""
        assert self._mlp is not None, "Cannot determine action without model"
        actions: torch.Tensor = self._mlp(
            torch.from_numpy(obs).float().to(self._device)
        )
        return actions.detach().numpy()

    def _construct_mlp(self, slices: list[slice]) -> SimpleMLP:
        """Constructs a MLP where the parameters at the given slice indices are
        perturbed."""
        mlp = self._mlp
        nn.utils.vector_to_parameters(
            vec=self._perturbed_params(slices)
            if len(slices) > 0
            else self.state.params,
            parameters=mlp.parameters(),
        )
        return mlp.to(self._device)

    def _perturbed_params(self, slices: list[slice]) -> torch.Tensor:
        """Perturb the state parameters with noise. The given slice indices are
        perturbed. The remaining are left untouched."""
        assert self._noise is not None, "Cannot perturb without noise"
        params = self.state.params
        sigma = self.state.noise_std
        eps = self._noise[slices]
        return params + sigma * eps

    def _generate_noise(self) -> None:
        """Generate noise for the perturbed parameters. This is done once per
        epoch."""
        assert self._mlp is not None, "Cannot determine noise without model"
        rng = torch.Generator(device=self._device).manual_seed(self.state.seed)
        n_params = nn.utils.parameters_to_vector(self._mlp.parameters()).numel()
        self._noise = torch.randn((self.population_size, n_params), generator=rng)

    @staticmethod
    def _transform_reward(rewards: torch.Tensor) -> torch.Tensor:
        """Transform the rewards to be normalized."""
        return (rewards - rewards.mean()) / rewards.std()


async def main() -> None:
    channel = grpc.insecure_channel("localhost:8080")
    worker = HalfCheetah(channel=channel, cores=1, render=True)
    await worker.start()


if __name__ == "__main__":
    asyncio.run(main())
