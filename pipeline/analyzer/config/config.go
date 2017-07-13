package config

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
	"path"
)

type Config struct {
	AppDir    string `json:"appdir"`
	UnpackDir string `json:"unpackdir"`
	SockPath  string `json:"sockpath"`
}

var UnpackDir string

func Load(cfgFile string) Config {
	file, err := os.Open(cfgFile)
	bytes, err := ioutil.ReadAll(file)
	if err != nil {
		panic("Couldn't read config file " + cfgFile)
	}
	var cfg Config
	err = json.Unmarshal(bytes, &cfg)
	if err != nil {
		panic("Error reading JSON: " + err.Error())
	}

	if cfg.AppDir == "" {
		cfg.AppDir = "./apks"
	}
	if cfg.UnpackDir == "" {
		cfg.UnpackDir = "."
	}
	if cfg.SockPath == "" {
		cfg.SockPath = "/run/apkScraper"
	}

	cfg.AppDir = path.Clean(cfg.AppDir)
	cfg.UnpackDir = path.Clean(cfg.UnpackDir)
	cfg.SockPath = path.Clean(cfg.SockPath)

	fmt.Println("Config:")
	fmt.Println("\tApp directory:", cfg.AppDir)
	fmt.Println("\tUnpacked app directory:", cfg.UnpackDir)
	fmt.Println("\tMessage socket path:", cfg.SockPath)

	return cfg
}
