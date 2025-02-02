{
  description = "New flake";
  inputs = {
    nixpkgs.url = "github:NixOS/nixpkgs";
    utils.url = "github:numtide/flake-utils";
    conixpkgs = {
      url = "github:cristianoliveira/nixpkgs";
      flake = true;
    };
  };
  outputs = { self, nixpkgs, utils, conixpkgs }: 
    utils.lib.eachDefaultSystem (system:
      let
        pkgs = import nixpkgs { inherit system; };
        copkgs = import conixpkgs { inherit pkgs; };
      in {
        devShells.default = import ./nix/dev-env.nix {
          inherit pkgs;
          inherit copkgs;
        };

        packages = pkgs.callPackage ./default.nix {
          inherit pkgs;
        };
    });
}
