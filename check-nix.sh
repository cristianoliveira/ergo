nix-build --no-out-link -E 'with import <nixpkgs> {}; callPackage ./default.nix {}'
