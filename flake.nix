{
  description = "New flake";
  inputs = {
    nixpkgs.url = "github:NixOS/nixpkgs";
    utils.url = "github:numtide/flake-utils";
    conixpkgs.url = "github:cristianoliveira/nixpkgs";
  };
  outputs = { self, nixpkgs, utils, conixpkgs }: 
    utils.lib.eachDefaultSystem (system:
      let
        pkgs = import nixpkgs { 
          inherit system;
          # Inject the namespace copkgs into the package set
          overlays = [ 
              (_: prev: {
                copkgs = {
                  funzzy = conixpkgs.packages."${system}".funzzyNightly;
                  ergoProxy = pkgs.callPackage ./nix/package.nix { inherit pkgs; };
                  nightly = pkgs.callPackage ./nix/package-nightly.nix { 
                    inherit pkgs;
                  };
                };
              }
            )
          ];
        };
      in {
        devShells.default = import ./nix/dev-env.nix {
          inherit pkgs;
        };

        packages = {
          ergoProxy = pkgs.copkgs.ergoProxy;
          nightly = pkgs.copkgs.nightly;
        };
    });
}
