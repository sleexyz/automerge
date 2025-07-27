{
  description = "wait-ci - Poll GitHub status checks and exit based on results";
  
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
        wait-ci = pkgs.callPackage ./package.nix { };
      in {
        packages = {
          default = wait-ci;
          wait-ci = wait-ci;
        };
        
        # Development shell for working on wait-ci
        devShells.default = pkgs.mkShell {
          buildInputs = with pkgs; [
            go
            git
            gh
          ];
          
          shellHook = ''
            echo "wait-ci development environment loaded"
            echo "Available commands: go build, go test, go run ."
          '';
        };
      }
    );
}