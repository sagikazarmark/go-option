{
  description = "Optional value using generics in Go";

  inputs = {
    nixpkgs.url = "nixpkgs/nixos-unstable";
    flake-utils.url = "github:numtide/flake-utils";
    gobin.url = "github:sagikazarmark/go-bin-flake";
    gobin.inputs.nixpkgs.follows = "nixpkgs";
  };

  outputs = { self, nixpkgs, flake-utils, gobin, ... }:
    flake-utils.lib.eachDefaultSystem (system:
      let
        pkgs = import nixpkgs {
          inherit system;

          overlays = [
            (final: prev: {
              golangci-lint = gobin.packages.${system}.golangci-lint-bin;
            })
          ];
        };

        buildDeps = with pkgs; [ git go_1_19 gnumake ];
        devDeps = with pkgs; buildDeps ++ [ golangci-lint ];
      in { devShell = pkgs.mkShell { buildInputs = devDeps; }; });
}
