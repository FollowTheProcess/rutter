_: {
  perSystem =
    { pkgs, ... }:
    {
      devShells.default = pkgs.mkShell {
        packages = with pkgs; [
          go_1_27
          golangci-lint
          goose
          goperf
          gopls
          goreleaser
          gotools
          just
          nilaway
          nix-update
          pkgsite
          sqlfluff
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
