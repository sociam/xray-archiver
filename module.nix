{ config, lib, pkgs, ... }:

with lib;

let
  cfg = config.services.xray;
  defaultUser = "xray";
in {
  options = {
    services.xray = {
      enable = mkEnableOption ''
        XRAY - some thing that does some other stuff.
      '';

      user = mkOption {
        type = types.string;
        default = defaultUser;
        description = ''
          XRAY will run with this user.
        '';
      };

      dataDir = mkOption {
        type = types.path;
        default = "/var/lib/xray";
      };

      package = mkOption {
        type = types.package;
        default = "";
        defaultText = "pkgs.xray";
        example = literalExample "pkgs.xray";
        description = "XRAY package to use.";
      };
    };
  };

  config = mkIf cfg.enable {
    systemd.packages = [ pkgs.nodejs-8_x cfg.package ];

    users = mkIf (cfg.user == defaultUser) {
      extraUsers."${defaultUser}" = {
        group = "nogroup";
        home = cfg.dataDir;
        createHome = true;
        uid = 500;
        description = "XRAY daemon user";
      };
    };

    systemd.services = {
      xray-analyzer = {
        after = [ "postgresql.service" ];
        requires = [ "postgresql.service" ];
        description = "XRAY analyzer";
        wantedBy = [ "multi-user.target" ];
        serviceConfig = {
          Restart = "always";
          User = cfg.user;
          Group = "nogroup";
          ExecStart = "${cfg.package}/bin/analyzer -daemon -db";
        };
      };

      xray-apiserv = {
        after = [ "network.target" "postgresql.service" ];
        requires = [ "postgresql.service" ];
        description = "XRAY API server";
        wantedBy = [ "multi-user.target" ];
        serviceConfig = {
          Restart = "always";
          User = cfg.user;
          Group = "nogroup";
          ExecStart = "${cfg.package}/bin/apiserv";
        };
      };

      xray-explorer = {
        after = [ "network.target" "postgresql.service" ];
        requires = [ "postgresql.service" ];
        description = "XRAY explorer";
        wantedBy = [ "multi-user.target" ];
        environment = {
          NODE_ENV = "production";
        };
        serviceConfig = {
          Restart = "always";
          User = cfg.user;
          Group = "nogroup";
          ExecStart = "${pkgs.nodejs-8_x}/bin/node ${cfg.package}/lib/xray/archiver/explorer/explorer.js";
        };
      };

      xray-downloader = {
        after = [ "network.target" "postgresql.service" ];
        requires = [ "postgresql.service" ];
        description = "XRAY explorer";
        wantedBy = [ "multi-user.target" ];
        environment = {
          NODE_ENV = "production";
        };
        serviceConfig = {
          Restart = "always";
          User = cfg.user;
          Group = "nogroup";
          ExecStart = "${pkgs.nodejs-8_x}/bin/node ${cfg.package}/lib/xray/archiver/downloader/downloader.js";
        };
      };

      xray-retriever = {
        after = [ "network.target" "postgresql.service" ];
        requires = [ "postgresql.service" ];
        description = "XRAY explorer";
        wantedBy = [ "multi-user.target" ];
        environment = {
          NODE_ENV = "production";
        };
        serviceConfig = {
          Restart = "always";
          User = cfg.user;
          Group = "nogroup";
          ExecStart = "${pkgs.nodejs-8_x}/bin/node ${cfg.package}/lib/xray/archiver/retriever/retriever.js";
        };
      };
    };
  };
}
