{
  pkgs,
  ...
}:
pkgs.buildGoModule {
  pname = "evochi";
  version = pkgs.lib.fileContents ../version;

  src = pkgs.lib.fileset.toSource {
    root = ../server;
    fileset = pkgs.lib.fileset.unions [
      ../server/go.mod
      ../server/go.sum

      ../server/cmd
      ../server/internal
      ../server/pkg
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
    license = pkgs.lib.licenses.mit;
    maintainers = with pkgs.lib.maintainers; [
      lukasl-dev
      MaxWolf-01
    ];
    platforms = pkgs.lib.platforms.linux;
    mainProgram = "evochi";
  };
}
