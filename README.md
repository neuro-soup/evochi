$$\LARGE {\color{white}\textrm{Evo}}{\color{gray}\textrm{lution~Or}}{\color{white}\textrm{ch}}{\color{gray}\textrm{estrat}}{\color{white}\textrm{i}}{\color{gray}\textrm{on}} $$

<div align="center">
    <img src="https://img.shields.io/badge/Written_In-Go-00acd7?style=for-the-badge&logo=go" alt="Go" />
    <img src="https://img.shields.io/badge/Library-Python-f7d44f?style=for-the-badge&logo=python" alt="Python" />
</div>

<div align="center">
    <img width="300" src="/assets/evochi.png" alt="evochi" />
</div>

Evochi is a framework-agnostic distributed training orchestrator for reinforcement
learning agents using [OpenAI's Evolution Strategies](https://arxiv.org/abs/1703.03864).

## Features

- 🔥 **Agnostic:** Evochi doesn't depend on any specific framework (or even programming language) for your workers. You define the format of your state for all workers.
- ⚡ **Fast:** Evochi's server is written in [Go](https://go.dev/) and uses [gRPC](https://grpc.io/) for fast communication.
- 📦 **Lightweight:** Evochi is designed to be as lightweight as possible on the server side. The computational workload is handled on the worker side.
- 📈 **Dynamically Scalable:** Evochi is built to scale horizontally and dynamically. Workers can leave or join at any time. As long as one worker remains in the workforce, the training can continue.
- 🚦 **Fault-Tolerance:** Evochi is fault-tolerant. If a worker crashes, mission-critical tasks can be recovered and delegated to other workers. As long as there is at least one functional worker, fault tolerance is ensured.

## Getting Started

### Start the Server

Binary releases are available on [GitHub](https://github.com/neuro-soup/evochi/releases).

Alternatively, you can run Evochi from (master) source using the `go run` command:

```bash
go run github.com/neuro-soup/evochi/cmd/evochi@latest
```

> [!IMPORTANT]
> Evochi requires some environment variables to be set. See the [server README](server/README.md#Environment-Variables) for all configuration options.

**Full (minimal) example:**
```bash
EVOCHI_JWT_SECRET="secret" EVOCHI_POPULATION_SIZE=50 go run github.com/neuro-soup/evochi/cmd/evochi@latest
```

<details>
    <summary>Run with <bold>Nix</bold></summary>

If you are using [Nix](https://nixos.org/), you can use the `nix run` command to run directly from source:

```bash
nix run "github:neuro-soup/evochi#server"
```

You can also import the package into your own Nix flake:

```nix
# flake.nix
inputs = {
    evochi.url = "github:neuro-soup/evochi";
};

# evochi.nix
{ inputs, pkgs, ... }:
{
    environment.systemPackages = [
        # installs `evochi` binary
        inputs.evochi.packages.${pkgs.system}.server
    ];
}
```

Alternatively, you can use evochi as Nix module:

```nix
# flake.nix
inputs = {
    evochi.url = "github:neuro-soup/evochi";
};

# evochi.nix
{ inputs, ... }:
{
    imports = [
        inputs.evochi.nixosModules.server
    ];

    services.evochi = {
        enable = true;
        config = {
            secret.file = ./evochi.secret;
            training.population = 100;
        };
    };
}
```

</details>


## (Real-World) Example Implementations

- [es-torch](https://github.com/neuro-soup/es-torch)
