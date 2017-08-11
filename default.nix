{ pkgs, system, stdenv, fetchFromGitHub, go, nodejs-8_x, bash, apktool, python27Packages }:

let
  pqsrc = fetchFromGitHub {
    owner = "lib";
    repo = "pq";
    rev = "e42267488fe361b9dc034be7a6bffef5b195bceb";
    sha256 = "082s0kb0qi1m32nq6jsnrn9wsiwm7pljhn9mh4sgkqajsxp830yk";
  };

  crap = import ./pipeline/default.nix { inherit pkgs system; };

  nodeDependencies = crap.shell.nodeDependencies;
in
stdenv.mkDerivation rec {
  name = "xray-${version}";
  version = "0.1.0";

  buildInputs = [ go nodejs-8_x bash nodeDependencies ];
  propagatedBuildInputs = [ apktool python27Packages.gplaycli ];

  srcs = ./pipeline;

  patchPhase = ''
    patchShebangs scripts/install.sh
  '';

  configurePhase = ''
    cd "$NIX_BUILD_TOP"

    mkdir goPath
    (cd goPath; unpackFile "${pqsrc}")
    mkdir -p "go/src/github.com/lib"
    chmod -R u+w goPath
    mv goPath/* "go/src/github.com/lib/pq"
    rmdir goPath

    mkdir -p "go/src/github.com/sociam/xray-archiver"
    mv pipeline "go/src/github.com/sociam/xray-archiver"

    cd go/src/github.com/sociam/xray-archiver/pipeline

    rm -rf node_modules
    [[ -d ${nodeDependencies}/lib/node_modules ]]
    ln -s ${nodeDependencies}/lib/node_modules .

    export GOPATH="$NIX_BUILD_TOP/go:$GOPATH"
  '';

  installPhase = ''
    makeFlagsArray=(PREFIX="$out" INSTALLFLAGS="-n")
  '';
}
