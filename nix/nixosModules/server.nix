self:
{
  config,
  lib,
  pkgs,
  ...
}:

let
  evochi = config.services.evochi;
in
{
  options = {
    services.evochi = {
      enable = lib.mkEnableOption "evochi";

      package = lib.mkPackageOption self.packages.${pkgs.system} "server" {
        default = "server";
        pkgsText = "evochi.packages.\${pkgs.system}.server";
      };

      service = lib.mkOption {
        description = ''Systemd service configuration for evochi'';
        type = lib.types.submodule {
          options = {

            enable = lib.mkOption {
              description = ''Whether to run evochi asa systemd service'';
              type = lib.types.bool;
              default = true;
            };

          };
        };
        default = { };
      };

      config = lib.mkOption {
        description = ''Application-specific configuration for evochi'';
        type = lib.types.submodule {
          options = {

            # required
            secret = lib.mkOption {
              description = ''The JWT secret to use.'';
              type = lib.types.submodule {
                options = {

                  file = lib.mkOption {
                    description = ''The path to the JWT secret file.'';
                    type = lib.types.nullOr lib.types.path;
                    default = null;
                  };

                  value = lib.mkOption {
                    description = ''The JWT secret value.'';
                    type = lib.types.nullOr lib.types.str;
                    default = null;
                  };

                };
              };
              default = { };
            };

            logging = lib.mkOption {
              description = ''Logging configuration for evochi'';
              type = lib.types.nullOr (
                lib.types.submodule {
                  options = {

                    level = lib.mkOption {
                      description = ''The log level to use.'';
                      type = lib.types.nullOr (
                        lib.types.enum [
                          "debug"
                          "info"
                          "warn"
                          "error"
                        ]
                      );
                      default = null;
                    };

                  };
                }
              );
              default = null;
            };

            server = lib.mkOption {
              type = lib.types.nullOr (
                lib.types.submodule {
                  options = {

                    port = lib.mkOption {
                      description = ''The port to listen on.'';
                      type = lib.types.nullOr lib.types.int;
                      default = null;
                    };

                  };
                }
              );
              default = null;
            };

            worker = lib.mkOption {
              type = lib.types.nullOr (
                lib.types.submodule {
                  options = {

                    timeout = lib.mkOption {
                      description = ''The task and heartbeat timeout for workers (in ms)'';
                      type = lib.types.nullOr lib.types.str;
                      default = null;
                    };

                    maximum = lib.mkOption {
                      description = ''The maximum number of workers to run.'';
                      type = lib.types.nullOr lib.types.int;
                      default = null;
                    };

                  };
                }
              );
              default = null;
            };

            training = lib.mkOption {
              type = lib.types.nullOr (
                lib.types.submodule {
                  options = {

                    population = lib.mkOption {
                      description = ''The size of the population to be evaluated'';
                      type = lib.types.nullOr lib.types.int;
                      default = null;
                    };

                    epochs = lib.mkOption {
                      description = ''The maximum number of epochs to run'';
                      type = lib.types.nullOr lib.types.int;
                      default = null;
                    };

                  };
                }
              );
              default = null;
            };

          };
        };
      };
    };
  };

  config = lib.mkIf evochi.enable (
    let
      envs = [
        {
          key = "EVOCHI_JWT_SECRET";
          value =
            if evochi.config.secret != null then
              if evochi.config.secret.file != null then
                lib.fileContents evochi.config.secret.file
              else if evochi.config.secret.value != null then
                evochi.config.secret.value
              else
                builtins.throw "services.evochi.config.secret.{file|value} (one of both) is required"
            else
              builtins.throw "services.evochi.config.secret.{file|value} (one of both) is required";
        }
        {
          key = "EVOCHI_LOG_LEVEL";
          value = if evochi.config.logging != null then evochi.config.logging.level else null;
        }
        {
          key = "EVOCHI_SERVER_PORT";
          value = if evochi.config.server != null then evochi.config.server.port else null;
        }
        {
          key = "EVOCHI_WORKER_TIMEOUT";
          value = if evochi.config.worker != null then evochi.config.worker.timeout else null;
        }
        {
          key = "EVOCHI_POPULATION_SIZE";
          value =
            if evochi.config.training != null then
              if evochi.config.training.population != null || evochi.config.training.population == 0 then
                evochi.config.training.population
              else
                builtins.throw "services.evochi.config.training.population is required"
            else
              builtins.throw "services.evochi.config.training.population is required";
        }
        {
          key = "EVOCHI_MAX_EPOCHS";
          value = if evochi.config.training != null then evochi.config.training.epochs else null;
        }
      ];
    in
    {
      environment.systemPackages = [ evochi.package ];

      systemd.services.evochi = lib.mkIf evochi.service.enable {
        enable = true;
        description = "evochi";
        after = [ "network.target" ];
        wantedBy = [ "multi-user.target" ];
        serviceConfig = {
          Environment = lib.map (env: "${env.key}=${toString env.value}") (
            lib.filter (env: env.value != null) envs
          );
          ExecStart = lib.getExe evochi.package;
          Restart = "on-failure";
        };
      };
    }
  );
}
