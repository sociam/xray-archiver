package util

import (
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"os"
	"path"
)

type DBCfg struct {
	Database string `json:"database"`
	User     string `json:"-"`
	Password string `json:"-"`
	Host     string `json:"host"`
	Port     int    `json:"port"`
}

type DBCreds struct {
	User     string `json:"user"`
	Password string `json:"password"`
}

type AnalyzerCfg struct {
	DB DBCreds `json:"db"`
}

type APIServCfg struct {
	DB DBCreds `json:"db"`
}

type Config struct {
	EDIHostname string      `json:"edihost"`
	DataDir     string      `json:"datadir"`
	AppDir      string      `json:"-"`
	UnpackDir   string      `json:"unpackdir"`
	Analyzer    AnalyzerCfg `json:"analyzer"`
	APIServ     APIServCfg  `json:"apiserv"`
	DB          DBCfg       `json:"db"`
}

var Cfg Config

const (
	Analyzer = iota
	APIServ
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

	if Cfg.EDIHostname == "" {
		Cfg.EDIHostname = "edi.sociam.org"
	}
	if Cfg.DataDir == "" {
		Cfg.DataDir = "/usr/local/var/xray"
	}
	Cfg.AppDir = path.Join(Cfg.DataDir, "apps")
	if Cfg.UnpackDir == "" {
		Cfg.UnpackDir, err = ioutil.TempDir("", "xray-analyzer")
	}

	Cfg.AppDir = path.Clean(Cfg.AppDir)
	Cfg.UnpackDir = path.Clean(Cfg.UnpackDir)

	switch requester {
	case Analyzer:
		Cfg.DB.User = Cfg.Analyzer.DB.User
		Cfg.DB.Password = Cfg.Analyzer.DB.Password
	case APIServ:
		Cfg.DB.User = Cfg.APIServ.DB.User
		Cfg.DB.Password = Cfg.APIServ.DB.Password
	}

	fmt.Println("Config:")
	fmt.Println("\tApp directory:", Cfg.AppDir)
	fmt.Println("\tUnpacked app directory:", Cfg.UnpackDir)

	return nil
}
