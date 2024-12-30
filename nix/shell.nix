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
    python312Packages.zstandard
    python312Packages.torch
    python312Packages.gymnasium
    python312Packages.jaxtyping
    python312Packages.numpy
  ];
}
