{
  pkgs ? import <nixpkgs> { },
}:

pkgs.mkShell {
  buildInputs = with pkgs; [
    # go
    go_1_23
    gopls
    gotools
    gofumpt
    golangci-lint

    # general
    just

    # protobuf
    grpcurl
    buf

    # python
    uv
    ruff
    python312
    python312Packages.grpcio
    python312Packages.grpcio-tools
  ];
}
