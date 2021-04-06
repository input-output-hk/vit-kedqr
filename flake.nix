{
  inputs = {
    nixpkgs.url = "github:NixOS/nixpkgs/nixos-unstable";
    utils.url = "github:kreisys/flake-utils";
  };
  outputs = { self, nixpkgs, utils }:
  utils.lib.simpleFlake {
    inherit nixpkgs;
    systems = [ "x86_64-linux" "aarch64-linux" ];
    overlay = final: prev:
    let lib = prev.lib;
    in {
      vit-kedqr = final.buildGoPackage rec {
        name = "vit-kedqr-${version}";
        version = "0.0.0";
        goPackagePath = "github.com/input-output-hk/vit-kedqr";
        src = final.nix-gitignore.gitignoreSource [] ./.;
        goDeps = ./deps.nix;
      };
    };
    packages = { vit-kedqr }@pkgs: pkgs;
    devShell = { mkShell, vgo2nix }: mkShell {
      buildInputs = [ vgo2nix ];
      shellHook = ''
        echo 'Run `vgo2nix` to regenerate deps.nix'
      '';
    };
  };
}
