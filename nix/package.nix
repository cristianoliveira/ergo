{ pkgs, ... }:
  let 
    _version = "v0.4.1";
  in
    pkgs.buildGoModule {
      # name of our derivation
      name = "ergo-proxy";
      version = "${_version}";

      # sources that will be used for our derivation.
      src = pkgs.fetchFromGitHub {
        owner = "cristianoliveira";
        repo = "ergo";
        rev = _version;
        sha256 = "sha256-/VhHuJsQrxewmfP/V25r5rZady0snfYNOGyIrv3vpGA=";
      };

      modSha256 = "0fagi529m1gf5jrqdlg9vxxq4yz9k9q8h92ch0gahp43kxfbgr4q";

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
