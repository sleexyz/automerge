# package.nix
{ lib
, buildGoModule
, gh
, git
}:

buildGoModule {
  pname = "wait-ci";
  version = "1.0.0";

  src = ./.;

  vendorHash = null;

  postInstall = ''
    # No shell completions for now, but could add later
  '';

  nativeBuildInputs = [ gh git ];

  ldflags = [
    "-s"
    "-w"
    "-X main.version=1.0.0"
  ];

  # Skip tests for now (no tests written yet)
  doCheck = false;

  meta = with lib; {
    description = "Poll GitHub status checks and exit based on results";
    homepage = "https://github.com/sleexyz/wait-ci";
    license = licenses.mit;
    maintainers = [ ];
    mainProgram = "wait-ci";
  };
}