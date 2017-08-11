{ stdenv, fetchFromGitHub, maven }:

stdenv.mkDerivation rec {
  name = "token-dispenser-${version}";
  version = "0.1.0";

  buildInputs = [ maven ];

  src = fetchFromGitHub {
      owner = "yeriomin";
      repo = "token-dispenser";
      rev = "ae40bc47b408363a70d65a10909c9154799fb79a";
      sha256 = "0k40skmj2309hwxsamjwagcllh4jcm48w7rmdwyzg666xikpx995";
  };

  configurePhase = ''
    sed -i -e "s:passwords\\.txt:/etc/xray/passwords.txt:g" src/main/resources/config.properties

  '';
  
  buildPhase = ''
    export M2_REPO=$(pwd)/.m2
    mvn install -Dmaven.repo.local=$(pwd)/.m2
  '';

  installPhase = '' 
    install -Dm644 target/token-dispenser.jar $out/lib/token-dispenser
  '';
}
