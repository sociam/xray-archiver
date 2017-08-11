{ stdenv, go, nodejs-8_x, bash, apktool, python27Packages }:

stdenv.mkDerivation rec {
  name = "xray-${version}";
  version = "0.1.0";

  buildInputs = [ go nodejs-8_x bash ];
  propagatedBuildInputs = [ apktool python27Packages.gplaycli ];

  src = ./pipeline;

  configurePhase = ''
    patchShebangs scripts/install.sh
  '';

  buildPhase = ''
    makeFlagsArray=(PREFIX="$out" INSTALLFLAGS="-n")
  '';
}
