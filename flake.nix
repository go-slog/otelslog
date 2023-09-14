{
  inputs = {
    # nixpkgs.url = "github:NixOS/nixpkgs/nixpkgs-unstable";
    nixpkgs.url = "github:NixOS/nixpkgs/master";
    flake-parts.url = "github:hercules-ci/flake-parts";
    devenv.url = "github:cachix/devenv";
  };

  outputs = inputs@{ flake-parts, ... }:
    flake-parts.lib.mkFlake { inherit inputs; } {
      imports = [
        inputs.devenv.flakeModule
      ];

      systems = [ "x86_64-linux" "x86_64-darwin" "aarch64-darwin" ];

      perSystem = { config, self', inputs', pkgs, system, ... }: rec {
        devenv.shells = {
          default = {
            languages = {
              go.enable = true;
              go.package = pkgs.lib.mkDefault pkgs.go_1_21;
            };

            packages = with pkgs; [
              # golangci-lint
            ] ++ [
              self'.packages.golangci-lint
            ];

            # https://github.com/cachix/devenv/issues/528#issuecomment-1556108767
            containers = pkgs.lib.mkForce { };
          };

          ci = devenv.shells.default;

          ci_1_21 = {
            imports = [ devenv.shells.ci ];

            languages = {
              go.package = pkgs.go_1_21;
            };
          };
        };

        packages = {
          golangci-lint = pkgs.buildGo121Module rec {
            pname = "golangci-lint";
            version = "1.54.2";

            src = pkgs.fetchFromGitHub {
              owner = "golangci";
              repo = "golangci-lint";
              rev = "v${version}";
              hash = "sha256-7nbgiUrp7S7sXt7uFXX8NHYbIRLZZQcg+18IdwAZBfE=";
            };

            vendorHash = "sha256-IyH5lG2a4zjsg/MUonCUiAgMl4xx8zSflRyzNgk8MR0=";

            subPackages = [ "cmd/golangci-lint" ];

            nativeBuildInputs = [ pkgs.installShellFiles ];

            ldflags = [
              "-s"
              "-w"
              "-X main.version=${version}"
              "-X main.commit=v${version}"
              "-X main.date=19700101-00:00:00"
            ];

            postInstall = ''
              for shell in bash zsh fish; do
                HOME=$TMPDIR $out/bin/golangci-lint completion $shell > golangci-lint.$shell
                installShellCompletion golangci-lint.$shell
              done
            '';

            meta = with pkgs.lib; {
              description = "Fast linters Runner for Go";
              homepage = "https://golangci-lint.run/";
              changelog = "https://github.com/golangci/golangci-lint/blob/v${version}/CHANGELOG.md";
              license = licenses.gpl3Plus;
              maintainers = with maintainers; [ anpryl manveru mic92 ];
            };
          };
        };
      };
    };
}
