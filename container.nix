{ config, pkgs, ... }:
let xraypkg = (pkgs.callPackage ./default.nix {}); in
{
  imports = [
    ./module.nix
  ];

  services.postgresql.enable = true;
  services.postgresql.package = pkgs.postgresql;

  networking.nameservers = [ "8.8.8.8" "8.8.4.4" ];

  environment.systemPackages = with pkgs; [
    apktool
    python27Packages.gplaycli
  ];

  services.xray.enable = true;
  services.xray.package = xraypkg;

  environment.etc."xray/config.json".text = ''
    {
        "datadir": "/var/lib/xray",
        "unpackdir": "/tmp/unpacked_apks",
        "credDownload": "/etc/xray/credentials.conf",
        "db": {
            "database": "xraydb",
            "host": "localhost",
            "port": 5432
        },
        "retriever": {
            "db": {
                "user": "postgres",
                "password": "aaaaa"
            }
        },
        "explorer": {
            "db": {
                "user": "postgres",
                "password": "aaaaa"
            }
        },
        "downloader": {
            "db": {
                "user": "postgres",
                "password": "aaaaa"
            }
        },
        "analyzer": {
            "db": {
                "user": "postgres",
                "password": "aaaaa"
            }
        },
        "apiserv": {
            "db": {
                "user": "postgres",
                "password": "aaaaa"
            }
        },
        "suggester": {
          "db": {
                "user": "postgres",
                "password": "aaaaa"
            }
        }
    }
  '';

  services.postgresql.initialScript = builtins.toFile "init_db.sql" (''
    create database xraydb;
    \connect xraydb
  '' + (builtins.readFile "${xraypkg}/lib/xray/init_db.sql") + ''
    alter user postgres with password 'aaaaa';
  '');
}
