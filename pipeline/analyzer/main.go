package main

import (
	"flag"
	"fmt"
	"github.com/sociam/xray-archiver/pipeline/db"
	"github.com/sociam/xray-archiver/pipeline/util"
	"log"
	"os"
	"time"
)

func analyze(app *util.App) error {
	var err error

	// if app.store == "cli" {
	// 	app.dbId, err = db.insertApp(app)
	// 	if err != nil {
	// 		fmt.Printf("Error getting database id of app %s: %s", app.id, err.Error())
	// 	}
	// }

	fmt.Print("Unpacking... ")
	err = app.Unpack()
	if err != nil {
		fmt.Println()
		fmt.Println(err.Error())
		if os.IsNotExist(err) {
			err := db.UnsetDownloaded(app.DBID)
			if err != nil {
				fmt.Printf("Failed to set %d not downloaded: %s", app.DBID, err.Error())
			}
		}
		return fmt.Errorf("Error unpacking apk: %s", err.Error())
	}
	fmt.Println("done.")

	fmt.Println("Getting permissions...")
	manifest, err := parseManifest(app)
	if err != nil {
		fmt.Println("Error parsing manifest: ", err.Error())
	} else {
		app.Perms = manifest.getPerms()
		fmt.Printf("Permissions found: %v\n\n", app.Perms)
		err = db.AddPerms(app)
		if err != nil {
			fmt.Printf("Error writing permissions to DB: %s\n", err.Error())
		}
	}

	fmt.Println("Running simple analysis... ")
	app.Hosts, err = simpleAnalyze(app)
	if err != nil {
		fmt.Printf("Error getting hosts: %s\n", err.Error())
	}
	fmt.Printf("Hosts found: %v\n\n", app.Hosts)
	// TODO: put this in a better place
	// err = db.AddHosts(app, app.Hosts)
	// if err != nil {
	// 	fmt.Printf("Error writing hosts to DB: %s\n", err.Error())
	// }

	app.Packages, err = findPackages(app)
	if err != nil {
		fmt.Println("Error finding packages: ", err.Error())
	} else {
		fmt.Println("Packages found: ", app.Packages)
		err = db.AddPackages(app)
		if err != nil {
			fmt.Printf("Error writing packages to DB: %s\n", err.Error())
		}
	}

	err = db.SetAnalyzed(app.DBID)
	if err != nil {
		fmt.Println("Error setting analyzed for app %d! This will result in looping!", app.DBID)
	}

	app.Cleanup()

	return nil
}

func runServer() {
	util.CheckDir(util.Cfg.UnpackDir, "Unpacked APK directory")

	for {
		apps, err := db.GetAppsToAnalyze()
		if err != nil || len(apps) == 0 {
			if err != nil {
				fmt.Println("Error getting apps to analyze from DB:", err.Error())
			}

			// Got no apps or a DB error occurred, sleep for 30 seconds.
			fmt.Println("Got no apps or errored, sleeping")
			time.Sleep(30 * time.Second)
		}

		for _, dbApp := range apps {
			app := dbApp.UtilApp()
			fmt.Printf("Got app %v\n", app)
			analyze(app)
		}
	}
}

var cfgFile = flag.String("cfg", "/etc/xray/config.json", "config file location")
var daemon = flag.Bool("daemon", false, "run analyzer as a daemon")
var useDb = flag.Bool("db", false, "add app information to the db specified in the config file")

func init() {
	var err error
	flag.Parse()
	util.LoadCfg(*cfgFile, util.Analyzer)
	err = db.Open(util.Cfg, *useDb)
	if err != nil {
		log.Fatal("Failed to open a connection to the database %s", err.Error())
	}
}

func main() {
	if err := os.MkdirAll(util.Cfg.DataDir, 0755); err != nil {
		panic(err)
	}
	if err := os.MkdirAll(util.Cfg.AppDir, 0755); err != nil {
		panic(err)
	}
	if err := os.MkdirAll(util.Cfg.UnpackDir, 0755); err != nil {
		panic(err)
	}

	if *daemon {
		fmt.Println("Starting xray analyzer daemon")
		runServer()
	} else {
		if flag.NArg() == 0 {
			flag.Usage()
			os.Exit(64)
		}

		for _, appPath := range flag.Args() {
			app := util.AppByPath(appPath)
			app.Store = "cli"
			fmt.Println("Analyzing apk ", appPath)
			analyze(app)
		}
	}
}
