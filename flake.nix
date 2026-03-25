{
  description = "A Nix-flake development environment";

  inputs.nixpkgs.url = "https://flakehub.com/f/NixOS/nixpkgs/0.1"; # unstable Nixpkgs

  outputs =
    { self, ... }@inputs:

    let
      goVersion = 26;
      nodeVersion = 24;

      supportedSystems = [
        "x86_64-linux"
        "aarch64-linux"
        "aarch64-darwin"
      ];
      forEachSupportedSystem =
        f:
        inputs.nixpkgs.lib.genAttrs supportedSystems (
          system:
          f {
            inherit system;
            pkgs = import inputs.nixpkgs {
              inherit system;
              overlays = [ inputs.self.overlays.default ];
            };
          }
        );
    in
    {
      overlays.default = final: prev: {
        go = final."go_1_${toString goVersion}";
        nodejs = final."nodejs_${toString nodeVersion}";
      };

      devShells = forEachSupportedSystem (
        { pkgs, system }:
        {
          default = pkgs.mkShellNoCC {
            packages = with pkgs; [
              # Golang (version is specified by overlay)
              go
              gotools
              golangci-lint
              gopls
              delve
              air

              # NodeJS (version is specified by overlay)
              nodejs
              nodePackages.pnpm
              turbo
              nodePackages.typescript
              nodePackages.typescript-language-server
              eslint
              prettier

              # Python
              uv

              # Git
              git

              # System deps
              gcc
              #gnumake
              #pkg-config

              # Infra
              #postgres
              #redis

              # utils
              wget
              curl
              jq
              nixfmt

              # Formatter
              self.formatter.${system}
            ];

            shellHook = ''
              echo -e "\n---------NIX DEVELOP---------\n"

              echo "Go: $(go version)"
              echo "NodeJS: $(node -v)"
              echo "pnpm: $(pnpm -v)"

              echo -e "\n---------NIX DEVELOP---------\n"

            '';
          };
        }
      );

      formatter = forEachSupportedSystem ({ pkgs, ... }: pkgs.nixfmt);
    };
}
