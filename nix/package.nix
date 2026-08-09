{ inputs, ... }:
let
  version = "0.1.0";
  rev = inputs.self.rev or inputs.self.dirtyRev or "unknown";
in
{
  perSystem =
    {
      pkgs,
      lib,
      self',
      ...
    }:
    {
      packages.default = pkgs.buildGoModule {
        meta = {
          description = "Sail through your shell history ⚓";
          homepage = "https://github.com/FollowTheProcess/rutter";
          license = lib.licenses.asl20;
          platforms = lib.platforms.unix;
          mainProgram = "rutter";
        };

        pname = "rutter";
        inherit version;
        src = lib.sources.cleanSource inputs.self;
        vendorHash = "sha256-ap+rlbLtrhOppB06dJsv9Iuz7RUAr0JG3QkXQD2w+Yg=";
        ldflags = [
          "-s"
          "-w"
          "-X go.followtheprocess.codes/rutter/internal/cmd.version=${version}"
          "-X go.followtheprocess.codes/rutter/internal/cmd.commit=${rev}"
          "-X go.followtheprocess.codes/rutter/internal/cmd.date=${inputs.self.lastModifiedDate}"
        ];

        env = {
          CGO_ENABLED = 0;
        };

        checkPhase = ''
          runHook preCheck
          CGO_ENABLED=1 go test -race ./...
          runHook postCheck
        '';
      };

      checks.rutter = self'.packages.default;
    };
}
