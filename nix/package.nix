{ buildGoPackage, nix-gitignore }:
buildGoPackage rec {
  name = "vit-kedqr-${version}";
  version = "0.0.0";
  goPackagePath = "github.com/input-output-hk/vit-kedqr";
  src = nix-gitignore.gitignoreSource [] ../.;
  goDeps = ./deps.nix;
}
