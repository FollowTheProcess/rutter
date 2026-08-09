_: {
  perSystem =
    { pkgs, ... }:
    {
      devShells.default = pkgs.mkShell {
        packages = with pkgs; [
          go
          golangci-lint
          goose
          goperf
          gopls
          goreleaser
          gotools
          mise
          nilaway
          nix-update
          pkgsite
          sqlc
          sqlite
          typos
        ];

        shellHook = ''
          echo "👋🏻 Welcome to the project devShell!"
        '';
      };
    };
}
