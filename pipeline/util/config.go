package util

import (
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"os"
	"path"
)

type DbCfg struct {
	Database string `json:"database"`
	User     string `json:"-"`
	Password string `json:"-"`
	Host     string `json:"host"`
	Port     int    `json:"port"`
}

type DbCreds struct {
	User     string `json:"user"`
	Password string `json:"password"`
}

type AnalyzerCfg struct {
	Db DbCreds `json:"db"`
}

type ApiServCfg struct {
	Db DbCreds `json:"db"`
}

type Config struct {
	EdiHostname string      `json:"edihost"`
	DataDir     string      `json:"datadir"`
	AppDir      string      `json:"-"`
	UnpackDir   string      `json:"unpackdir"`
	SockPath    string      `json:"sockpath"`
	Analyzer    AnalyzerCfg `json:"analyzer"`
	ApiServ     ApiServCfg  `json:"apiserv"`
	Db          DbCfg       `json:"db"`
}

var Cfg Config

const (
	Analyzer = iota
	ApiServ
)

func LoadCfg(cfgFile string, requester int) error {
	file, err := os.Open(cfgFile)
	bytes, err := ioutil.ReadAll(file)
	if err != nil {
		return errors.New("Couldn't read config file " + cfgFile)
	}
	err = json.Unmarshal(bytes, &Cfg)
	if err != nil {
		return errors.New("Error reading JSON: " + err.Error())
	}

	if Cfg.EdiHostname == "" {
		Cfg.EdiHostname = "edi.sociam.org"
	}
	if Cfg.DataDir == "" {
		Cfg.DataDir = "/usr/local/var/xray"
	}
	Cfg.AppDir = path.Join(Cfg.DataDir, "apps")
	if Cfg.UnpackDir == "" {
		Cfg.UnpackDir, err = ioutil.TempDir("", "xray-analyzer")
	}
	if Cfg.SockPath == "" {
		Cfg.SockPath = "/var/run/apkScraper"
	}

	Cfg.AppDir = path.Clean(Cfg.AppDir)
	Cfg.UnpackDir = path.Clean(Cfg.UnpackDir)
	Cfg.SockPath = path.Clean(Cfg.SockPath)

	switch requester {
	case Analyzer:
		Cfg.Db.User = Cfg.Analyzer.Db.User
		Cfg.Db.Password = Cfg.Analyzer.Db.Password
	case ApiServ:
		Cfg.Db.User = Cfg.ApiServ.Db.User
		Cfg.Db.Password = Cfg.ApiServ.Db.Password
	}

	fmt.Println("Config:")
	fmt.Println("\tApp directory:", Cfg.AppDir)
	fmt.Println("\tUnpacked app directory:", Cfg.UnpackDir)
	fmt.Println("\tMessage socket path:", Cfg.SockPath)

	return nil
}
