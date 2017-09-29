package util

import (
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"os"
	"path"
)

// DBCfg Struct for the Database Config File information
type DBCfg struct {
	Database string `json:"database"`
	User     string `json:"-"`
	Password string `json:"-"`
	Host     string `json:"host"`
	Port     int    `json:"port"`
}

// DBCreds Struct for the Database Credentials
type DBCreds struct {
	User     string `json:"user"`
	Password string `json:"password"`
}

// AnalyzerCfg Represents the Credentials used to connect to the DB
// as the Analyser
type AnalyzerCfg struct {
	DB DBCreds `json:"db"`
}

// APIServCfg Represents the Credentials used to connect to the DB
// as the API Server
type APIServCfg struct {
	DB DBCreds `json:"db"`
}

// Config Struct for the Xray Config information relating to Dir and file
// locations. As well as holding DB, Analyser and APIServ Config
// information.
type Config struct {
	GeoIPEndpoint string      `json:"geoipurl"`
	DataDir       string      `json:"datadir"`
	AppDir        string      `json:"-"`
	UnpackDir     string      `json:"unpackdir"`
	Analyzer      AnalyzerCfg `json:"analyzer"`
	APIServ       APIServCfg  `json:"apiserv"`
	DB            DBCfg       `json:"db"`
}

// Cfg is an instance of Config that will contain DB settings
// as well as Dir/file locations and information.
var Cfg Config

// Enum used when loading the config in LoadCfg
const (
	Analyzer = iota
	APIServ
)

// LoadCfg Opens a config file and creates a series of objects
// using the information located in the file. It constructs a
// Config, populating information for the Analyser Config,
// API Server Config and the DB config.
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

	if Cfg.GeoIPEndpoint == "" {
		Cfg.GeoIPEndpoint = "http://localhost/geoip"
	}
	if Cfg.DataDir == "" {
		Cfg.DataDir = "/var/xray"
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
