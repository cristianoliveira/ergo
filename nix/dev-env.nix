{ pkgs ? import <nixpkgs> {}, copkgs }:
  pkgs.mkShell {
    # buildInputs is for dependencies you'd need "at run time",
    # were you to to use nix-build not nix-shell and build whatever you were working on
    buildInputs = [
      pkgs.go
      copkgs.funzzy
    ];

    shell = pkgs.zsh;
  }
