{
  pkgs ? import <nixpkgs> { },
}:

pkgs.mkShell rec {
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
    python312Packages.wandb

    glib
    glibc
    zstd
    wayland
    glfw-wayland
    qt5.full
    libsForQt5.qt5.qtwayland
    libGL
    libGLU
    gcc
    zlib
    qt5.qtbase
    libsForQt5.qt5.qtwayland
    SDL2
    SDL2.dev
    xorg.libX11
    xorg.libX11.dev
    xorg.libXext
    xorg.libXext.dev
    xorg.libXrandr
    xorg.libXrandr.dev
  ];

  shellHook = ''
    export LD_LIBRARY_PATH=${pkgs.lib.makeLibraryPath buildInputs}
  '';
}
