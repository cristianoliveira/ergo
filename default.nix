# allow our nixpkgs import to be overridden if desired
{ pkgs ? import <nixpkgs> {}, fetchFromGitHub, buildGoModule, ... }:

# let's write an actual basic derivation
# this uses the standard nixpkgs mkDerivation function
buildGoModule rec {

  # name of our derivation
  name = "ergo-proxy";
  version = "0.3.2";

  # sources that will be used for our derivation.
  src = fetchFromGitHub {
    owner = "cristianoliveira";
    repo = "ergo";
    rev = "v0.3.2";
    sha256 = "sha256-C3lJWuRyGuvo33kvj3ktWKYuUZ2yJ8pDBNX7YZn6wNM=";
  };

  modSha256 = "0fagi529m1gf5jrqdlg9vxxq4yz9k9q8h92ch0gahp43kxfbgr4q";

  vendorHash = "sha256-yXWM59zoZkQPpOouJoUd5KWfpdCBw47wknb+hYy6rh0=";

  meta = with pkgs.lib; {
    description = "Ergo: The reverse proxy agent for local domain management";
    homepage = "https://github.com/cristianoliveira/ergo";
    license = licenses.mit;
    maintainers = with maintainers; [ cristianoliveira ];
  };
}
