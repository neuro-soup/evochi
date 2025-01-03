# server

## Environment Variables

The server is configurable via environment variables. The following configuration
options are available:

| Variable | Options | Default | Description |
| -------- | ------- | ------- | ----------- |
| `EVOCHI_JWT_SECRET` | `string` | **required** | The secret to use for JWT tokens.
| `EVOCHI_POPULATION_SIZE` | `uint` | **required** | The size of the population to use.
| `EVOCHI_LOG_LEVEL` | `debug`, `info`, `warn`, or `error` | `info` | The log level to use.
| `EVOCHI_SERVER_PORT` | `uint` | `8080` | The port to listen on.
| `EVOCHI_WORKER_TIMEOUT` | `duration` | `1m` | The task and heartbeat timeout for workers.
| `EVOCHI_MAX_WORKERS` | `uint` | `0` | The maximum number of workers to run. If set to 0, no limit is set.
| `EVOCHI_MAX_EPOCHS` | `uint` | `0` | The maximum number of epochs to run. If set to 0, no limit is set.
