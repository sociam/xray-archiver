package main

import (
	"github.com/sociam/xray-archiver/pipeline/db"
)

func main() {
	// Select app Host app IDs.
	// for all app_host records
	// for all hosts in app host_records
	// Map host name to company.
	// insert company if new
	// insert company app association if new.

	appIDs, _ := db.GetAppHostIDs()
}
