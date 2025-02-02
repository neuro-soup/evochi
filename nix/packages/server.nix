{
  buildGoModule,
  lib,
  ...
}:
buildGoModule {
  pname = "evochi";
  version = lib.fileContents ../../version;

  src = lib.fileset.toSource {
    root = ../../server;
    fileset = lib.fileset.unions [
      ../../server/go.mod
      ../../server/go.sum

      ../../server/cmd
      ../../server/internal
      ../../server/pkg
    ];
  };
  subPackages = [ "cmd/server" ];
  vendorHash = "sha256-JUnar1qY5phW9OZ3Egzczlq/c73oLQZ/5RuwCjh1BH0=";

  postInstall = ''
    mv $out/bin/server $out/bin/evochi
  '';

  meta = {
    description = "A distributed training orchestrator inspired by OpenAI's Evolution Strategies paper.";
    homepage = "https://github.com/neuro-soup/evochi";
    license = lib.licenses.mit;
    maintainers = with lib.maintainers; [
      lukasl-dev
      MaxWolf-01
    ];
    platforms = lib.platforms.linux;
    mainProgram = "evochi";
  };
}
