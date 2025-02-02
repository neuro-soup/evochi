{
  description = "evochi";
  inputs = {
    nixpkgs.url = "github:NixOS/nixpkgs/nixos-unstable";
    utils.url = "github:numtide/flake-utils";
  };
  outputs =
    {
      self,
      nixpkgs,
      utils,
      ...
    }:
    (utils.lib.eachDefaultSystem (
      system:
      let
        pkgs = import nixpkgs { inherit system; };
      in
      {
        devShells = {
          default = pkgs.callPackage ./nix/shell.nix { };
        };

        packages = {
          server = pkgs.callPackage ./nix/packages/server.nix { };
        };
      }
    ))
    // {
      nixosModules = {
        server = import ./nix/nixosModules/server.nix self;
      };
    };
}
