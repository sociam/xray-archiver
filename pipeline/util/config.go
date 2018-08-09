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
	GeoIPEndpoint string        `json:"geoipurl"`
	StorageConfig StorageConfig `json:"storage_config"`
	SystemConfig  SystemConfig  `json:"system_config"`
	Analyzer      AnalyzerCfg   `json:"analyzer"`
	APIServ       APIServCfg    `json:"apiserv"`
	DB            DBCfg         `json:"db"`
}

// SystemConfig represents the config info related to the system the program
// is running on.
type SystemConfig struct {
	VMName                string `json:"vm_name"`
	DownloaderCredentials string `json:"downloader_credentials"`
}

// StorageConfig holds the config data related to where APK data
// may be stored.
type StorageConfig struct {
	APKDownloadDirectories []APKDownloadDirectory `json:"apk_download_directories"`
	APKUnpackDirectory     string                 `json:"apk_unpack_directory"`
	MinimumGBRequired      string                 `json:"minimum_gb_required"`
}

// APKDownloadDirectory represents a possible location an APK could be stored on.
type APKDownloadDirectory struct {
	Name string `json:"name"`
	Path string `json:"path"`
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

	Cfg.StorageConfig.APKUnpackDirectory = path.Clean(Cfg.StorageConfig.APKUnpackDirectory)

	switch requester {
	case Analyzer:
		Cfg.DB.User = Cfg.Analyzer.DB.User
		Cfg.DB.Password = Cfg.Analyzer.DB.Password
	case APIServ:
		Cfg.DB.User = Cfg.APIServ.DB.User
		Cfg.DB.Password = Cfg.APIServ.DB.Password
	}

	fmt.Println("Config:")
	fmt.Println("\tApp directories:", Cfg.StorageConfig.APKDownloadDirectories)
	fmt.Println("\tUnpacked app directory:", Cfg.StorageConfig.APKUnpackDirectory)

	return nil
}
