package main

import (
	"flag"
	"fmt"
	"log"
	"net/http"
	"regexp"
	"strconv"
	"strings"

	"github.com/sociam/xray-archiver/pipeline/db"
	"github.com/sociam/xray-archiver/pipeline/util"
)

// Err - convenience struct for marshalling errors
type Err struct {
	Code    string `json:"err"`
	Message string `json:"err_msg"`
}

var unit = util.Unit{}

var supportedMimes = map[string]util.Unit{
	"application/json":    unit,
	"application/nahmate": unit,
}

func toBytes(data interface{}) []byte {
	switch v := data.(type) {
	case Err:
		return []byte(v.Message)
	default:
		return []byte(fmt.Sprintf("%v", v))
	}
}

func writeErr(w http.ResponseWriter, mime string, status int, err, msg string, vals ...interface{}) {
	writeData(w, mime, status, Err{err, fmt.Sprintf(msg, vals...)})
}

func writeData(w http.ResponseWriter, mime string, status int, data interface{}) {
	w.WriteHeader(status)
	var err1 error
	w.Header().Set("Content-Type", mime)
	switch mime {
	case "application/nahmate":
		err1 = util.WriteDEAN(w, data)
	case "text/plain":
		_, err1 = w.Write(toBytes(data))
	default:
		w.Header().Set("Content-Type", "application/json")
		fallthrough
	case "application/json":
		err1 = util.WriteJSON(w, data)
	}
	if err1 != nil {
		log.Println(err1)
	}
}

func hello(w http.ResponseWriter, r *http.Request) {
	fmt.Println("Got spurious request on " + r.URL.Path)
	writeErr(w, r.Header.Get("Accept"), http.StatusNotFound, "not_found", "Nah mate!")
}

var appPrefixRe = regexp.MustCompile("^/api/apps/")
var dbIDRe = regexp.MustCompile("^\\d+$")
var appIDRe = regexp.MustCompile("^[[:alpha:]][\\w$]*(\\.[[:alpha:]][\\w$]*)*$")

func appEndpoint(w http.ResponseWriter, r *http.Request) {
	mime := r.Header.Get("Accept")
	if r.Method == "POST" || r.Method == "GET" {
		if _, ok := supportedMimes[mime]; !ok {
			writeErr(w, mime, http.StatusNotAcceptable, "not_acceptable", "This API only supports JSON at the moment.")
			return
		}

		split := strings.SplitN(r.URL.Path, "/", 7)

		if len(split) < 4 {
			writeErr(w, mime, http.StatusBadRequest, "bad_app", "Bad app slashes specified")
			return
		}

		appID := split[3]
		fmt.Println("AppId searching:", appID)
		if appID == "" {
			appsEndpoint(w, r)
		} else if len(split) == 4 && dbIDRe.MatchString(appID) {
			// Is a DB ID
			///api/apps/<dbid>

			dbID, err := strconv.Atoi(appID)
			if err != nil {
				writeErr(w, mime, http.StatusBadRequest, "big_int", "dbID is too big")
			}

			appVer, err := db.GetAppVersionById(dbID)
			if err != nil {
				writeErr(w, mime, http.StatusBadRequest, "bad_app", "App could not be found")
				return
			}

			// util.WriteJSON(w, app)

		} else if appIDRe.MatchString(appID) {
			// Is an app ID
			switch len(split) {
			case 4:
				app, err := db.GetApp(appID)
				if err != nil {
					fmt.Println("Error querying database: ", err.Error())
					writeErr(w, mime, http.StatusInternalServerError, "internal_error", "An internal error occurred")
					return
				}

				writeData(w, mime, http.StatusOK, app)
			case 5:
				///api/apps/<appid>/<version string>

				appVer, err := db.GetAppVersion(store, region, ver)
				if err != nil {
					fmt.Println("Error querying database: ", err.Error())
					writeErr(w, mime, http.StatusInternalServerError, "internal_error", "An internal error occurred")
				}
			case 7:
				///api/apps/<appid>/<store>/<region>/<version string>

				// There are parts after the app ID
				store, region, ver := split[4], split[5], split[6]

				appVer, err := db.GetAppVersion(store, region, ver)
				if err != nil {
					fmt.Println("Error querying database: ", err.Error())
					writeErr(w, mime, http.StatusInternalServerError, "internal_error", "An internal error occurred")
				}
			default:
				writeErr(w, mime, "bad_req", "Number of parts is not 1, 2, or 4")
			}

		} else {
			writeErr(w, mime, http.StatusBadRequest, "bad_app", "Invalid app ID specified")
		}
	}
}

func appsEndpoint(w http.ResponseWriter, r *http.Request) {
	//TODO: handle OPTIONS and HEAD
	mime := r.Header.Get("Accept")
	if r.Method == "POST" || r.Method == "GET" {
		if _, ok := supportedMimes[mime]; !ok {
			writeErr(w, mime, http.StatusNotAcceptable, "not_acceptable", "This API only supports JSON at the moment.")
			return
		}

		err := r.ParseForm()
		if err != nil {
			writeErr(w, mime, http.StatusBadRequest, "bad_form", "Error parsing form input: %s", err.Error())
			return
		}

		num, start := 10, 0
		if v, ok := r.Form["num"]; ok && len(v) > 0 {
			if len(v) > 1 {
				writeErr(w, mime, http.StatusBadRequest, "bad_form", "num must have a single value")
				return
			}
			num, err = strconv.Atoi(v[0])
			if err != nil {
				writeErr(w, mime, http.StatusBadRequest, "bad_form", "num value must be a number")
				return
			}
			if num < 1 || num > 100 {
				writeErr(w, mime, http.StatusBadRequest, "bad_form", "num value must be between 1 and 100")
				return
			}
		}

		if v, ok := r.Form["start"]; ok && len(v) > 0 {
			if len(v) > 1 {
				writeErr(w, mime, http.StatusBadRequest, "bad_form", "start must have a single value")
				return
			}
			start, err = strconv.Atoi(v[0])
			if err != nil {
				writeErr(w, mime, http.StatusBadRequest, "bad_form", "start value must be a number")
				return
			}
			if start < 0 {
				writeErr(w, mime, http.StatusBadRequest, "bad_form", "start value must positive")
				return
			}
		}

		apps, err := db.GetApps(num, start)
		if err != nil {
			fmt.Println("Error querying database: ", err.Error())
			writeErr(w, mime, http.StatusInternalServerError, "internal_error", "An internal error occurred")
			return
		}

		util.WriteJSON(w, apps)
	} else {
		writeErr(w, mime, http.StatusBadRequest, "bad_method", "You must POST or GET this endpoint!")
	}
}

var cfgFile = flag.String("cfg", "/etc/xray/config.json", "config file location")

func init() {
	util.LoadCfg(*cfgFile, util.ApiServ)
	db.Open(util.Cfg)
}

func main() {
	http.HandleFunc("/", hello)
	http.HandleFunc("/api/apps", appsEndpoint)
	http.HandleFunc("/api/apps/", appEndpoint)
	http.HandleFunc("/api/developers", devEndpoint)
	http.HandleFunc("/api/developers/", devsEndpoint)
	http.HandleFunc("/api/companies", compEndpoint)
	http.HandleFunc("/api/companies/", compsEndpoint)
	panic(http.ListenAndServe(":8080", nil))
}
