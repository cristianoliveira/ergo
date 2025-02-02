{ pkgs ? import <nixpkgs> {} }:
{
  default = pkgs.callPackage ./nix/package.nix { inherit pkgs; };
  nightly = pkgs.callPackage ./nix/package-nightly.nix { 
  inherit pkgs;
  };
}
