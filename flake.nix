{
  description = "New flake";
  inputs = {
    nixpkgs.url = "github:NixOS/nixpkgs";
    utils.url = "github:numtide/flake-utils";
  };
  outputs = { self, nixpkgs, utils }: 
    utils.lib.eachDefaultSystem (system:
      let
        pkgs = import nixpkgs { 
          inherit system;
          # Inject the namespace copkgs into the package set
          overlays = [ 
              (_: prev: {
                copkgs = {
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
          default = pkgs.copkgs.ergoProxy;
          ergoProxy = pkgs.copkgs.ergoProxy;
          nightly = pkgs.copkgs.nightly;
        };
    });
}
