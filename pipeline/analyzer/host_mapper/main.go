package main

import (
	"flag"
	"log"
	"strings"

	"github.com/sociam/xray-archiver/pipeline/db"
	"github.com/sociam/xray-archiver/pipeline/util"
)

var cfgFile = flag.String("cfg", "/etc/xray/config.json", "config file location")

func init() {
	var err error
	flag.Parse()
	err = util.LoadCfg(*cfgFile, util.Analyzer)
	if err != nil {
		log.Fatalf("Failed to read config: %s", err.Error())
	}
	err = db.Open(util.Cfg, true)
	if err != nil {
		log.Fatalf("Failed to open a connection to the database: %s", err.Error())
	}
}

func main() {
	// Select app Host app IDs.
	// for all app_host records
	// for all hosts in app host_records
	// Map host name to company.
	// insert company if new
	// insert company app association if new.

	appIDs, _ := db.GetAppHostIDs()

	for i := 0; i < len(appIDs); i++ {
		appHostRecord, _ := db.GetAppHostsByID(appIDs[i])
		util.Log.Debug("Host Names found in app with Version ID: %d", appHostRecord.ID)
		util.Log.Debug(strings.Join(appHostRecord.HostNames, ", "))
	}
}
