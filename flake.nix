{
  description = "automerge - Poll GitHub status checks and exit based on results";
  
  # Set the flake name for profile installation
  nixConfig = {
    flake-registry = "https://github.com/NixOS/flake-registry/raw/master/flake-registry.json";
  };
  
  inputs = {
    nixpkgs.url = "github:NixOS/nixpkgs/nixos-unstable";
    flake-utils.url = "github:numtide/flake-utils";
  };
  
  outputs = { self, nixpkgs, flake-utils }:
    flake-utils.lib.eachDefaultSystem (system:
      let
        pkgs = nixpkgs.legacyPackages.${system};
        automerge = pkgs.callPackage ./automerge.nix { };
      in {
        packages = {
          default = automerge;
          automerge = automerge;
        };
        
        # Development shell for working on automerge
        devShells.default = pkgs.mkShell {
          buildInputs = with pkgs; [
            go
            git
            gh
          ];
          
          shellHook = ''
            echo "automerge development environment loaded"
            echo "Available commands: go build, go test, go run ."
          '';
        };
      }
    );
}