{ pkgs, ... }:
  let 
    _version = "v0.4.1-nighlty";
  in
    pkgs.buildGoModule {
      # name of our derivation
      name = "ergo-proxy";
      version = "${_version}";
      doCheck = true;

      # sources that will be used for our derivation.
      src = ../.;

      modSha256 = "sha256-yXWM59zoZkQPpOouJoUd5KWfpdCBw47wknb+hYy6rh0=";

      vendorHash = "sha256-yXWM59zoZkQPpOouJoUd5KWfpdCBw47wknb+hYy6rh0=";

      ldflags = [
        "-s" "-w"
        "-X main.VERSION=${_version}"
      ];

      meta = with pkgs.lib; {
        description = "Ergo: The reverse proxy agent for local domain management";
        homepage = "https://github.com/cristianoliveira/ergo";
        license = licenses.mit;
        maintainers = with maintainers; [ cristianoliveira ];
      };
    }
