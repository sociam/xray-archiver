package main

import (
	"bytes"
	"encoding/json"
	"flag"
	"log"
	"net/http"

	"github.com/sociam/xray-archiver/pipeline/db"
	"github.com/sociam/xray-archiver/pipeline/util"
)

func requestTrackerMapping(appHostRecord db.AppHostRecord) []db.TrackerMapperCompany {
	tmReqData := db.TrackerMapperRequest{appHostRecord.HostNames}
	// BODY: {"host_names":["facebook.com", "360.jp.co"]}
	// URL: localhost:8080/hosts
	// REQUEST TYPE: Post

	url := "http://127.0.0.1:8080/hosts" // Get from some config file or something...

	// Encode Object
	ioBuffer := new(bytes.Buffer)
	json.NewEncoder(ioBuffer).Encode(tmReqData)

	// Form Request and set headers.
	req, err := http.NewRequest("POST", url, ioBuffer)
	req.Header.Set("Content-Type", "application/json")

	// Check for errors forming request.
	if err != nil {
		util.Log.Err("Error forming TrackerMapper API Request.", err)
	}

	// carry out the request.
	client := &http.Client{}
	resp, err := client.Do(req)

	// check for errors carrying out the request
	if err != nil {
		util.Log.Err("Client Error issueing Tracker Mapper API request..", err)
	}

	// Check there is a response body.
	if resp.Body != nil {
		defer resp.Body.Close()
	}

	// Decode the response and check for error.
	var tmCompanies []db.TrackerMapperCompany
	if err := json.NewDecoder(resp.Body).Decode(&tmCompanies); err != nil {
		util.Log.Err("Error Decoding Response Body from TrackerMapper API.", err)
	}
	return tmCompanies
}

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

		tmCompanies := requestTrackerMapping(appHostRecord)

		for j := 0; j < len(tmCompanies); j++ {
			// Insert Company App Association into the Database.
			db.InsertCompanyName(tmCompanies[j].CompanyName)
			db.InsertCompanyAppAssociation(appIDs[i], tmCompanies[j].CompanyName)

			util.Log.Debug("Company Name: %s, Host Name: %s", tmCompanies[j].CompanyName, tmCompanies[j].HostName)
		}
	}
}
