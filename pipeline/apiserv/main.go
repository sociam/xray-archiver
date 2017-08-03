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

func isFormatRequestSupported(w http.ResponseWriter, mime string, r *http.Request) bool {

	if _, ok := supportedMimes[mime]; !ok {
		writeErr(w, mime, http.StatusNotAcceptable, "not_acceptable", "This API only supports JSON at the moment.")
		return ok
	}

	return true
}

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

		if appID == "" {
			appsEndpoint(w, r)
		} else if len(split) == 4 && dbIDRe.MatchString(appID) {
			// Is a DB ID
			///api/apps/<dbid>

			dbID, err := strconv.Atoi(appID)
			if err != nil {
				writeErr(w, mime, http.StatusBadRequest, "big_int", "dbID is too big")
			}

			appVer, err := db.GetAppVersionByID(int64(dbID))
			if err != nil {
				writeErr(w, mime, http.StatusBadRequest, "bad_app", "App could not be found")
				return
			}

			writeData(w, mime, http.StatusOK, appVer)

			//util.WriteJSON(w, app)

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
				store, region, ver := split[4], split[5], split[6]

				appVer, err := db.GetAppVersion(store, region, ver)
				if err != nil {
					fmt.Println("Error querying database: ", err.Error())
					writeErr(w, mime, http.StatusInternalServerError, "internal_error", "An internal error occurred")
				}

				writeData(w, mime, http.StatusOK, appVer)

			case 7:
				///api/apps/<appid>/<store>/<region>/<version string>

				// There are parts after the app ID
				//TODO: do things not need to be done different here?
				store, region, ver := split[4], split[5], split[6]

				appVer, err := db.GetAppVersion(store, region, ver)
				if err != nil {
					fmt.Println("Error querying database: ", err.Error())
					writeErr(w, mime, http.StatusInternalServerError, "internal_error", "An internal error occurred")
				}

				writeData(w, mime, http.StatusOK, appVer)

			default:
				writeErr(w, mime, http.StatusBadRequest, "bad_req", "Number of parts is not 1, 2, or 4")
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
		fmt.Println("Gathering app with heads", r)
		err := r.ParseForm()
		if err != nil {
			writeErr(w, mime, http.StatusBadRequest, "bad_form", "Error parsing form input: %s", err.Error())
			return
		}

		num, start := 10, 0 //Default selection range
		for name, val := range r.Form {
			switch name {

			case "num":
				if len(val) > 1 {
					writeErr(w, mime, http.StatusBadRequest, "bad_form", "num must have a single value")
					return
				}
				num, err = strconv.Atoi(val[0])
				if err != nil {
					writeErr(w, mime, http.StatusBadRequest, "bad_form", "num value must be a number")
					return
				}
				if num < 1 || num > 100 {
					writeErr(w, mime, http.StatusBadRequest, "bad_form", "num value must be between 1 and 100")
					return
				}

			case "start":
				if len(val) > 1 {
					writeErr(w, mime, http.StatusBadRequest, "bad_form", "start must have a single value")
					return
				}
				start, err = strconv.Atoi(val[0])
				if err != nil {
					writeErr(w, mime, http.StatusBadRequest, "bad_form", "start value must be a number")
					return
				}
				if start < 0 {
					writeErr(w, mime, http.StatusBadRequest, "bad_form", "start value must positive")
					return
				}
			default:
				writeErr(w, mime, http.StatusBadRequest, "bad_form", "passed values did not match params", name)
				return
			}
		}

		fmt.Println("Selecting app range", start, num)

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

// What about? ^[a-z0-9_-]*$
//var compIDRe = regexp.MustCompile("^\\l+$")
var compIDRe = regexp.MustCompile("^[a-z0-9_-]*$")

func compEndpoint(w http.ResponseWriter, r *http.Request) {
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

		compID := split[3]
		fmt.Println("CompID searching:", compID)

		if len(split) > 4 {
			compsEndpoint(w, r)

		} else if len(split) == 4 && compIDRe.MatchString(compID) {
			company, err := db.GetCompany(string(compID))

			if err != nil {
				writeErr(w, mime, http.StatusBadRequest, "bad_company", "Company could not be found")
				return
			}

			writeData(w, mime, http.StatusOK, company)

			util.WriteJSON(w, company)

		} else {
			writeErr(w, mime, http.StatusBadRequest, "bad_company", "Invalid Company ID specified")
		}
	}
}

func compsEndpoint(w http.ResponseWriter, r *http.Request) {
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

		companies, err := db.GetCompanies(num, start)
		if err != nil {
			fmt.Println("Error querying database: ", err.Error())
			writeErr(w, mime, http.StatusInternalServerError, "internal_error", "An internal error occurred")
			return
		}

		util.WriteJSON(w, companies)
	} else {
		writeErr(w, mime, http.StatusBadRequest, "bad_method", "You must POST or GET this endpoint!")
	}
}

func devsEndpoint(w http.ResponseWriter, r *http.Request) {
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

		developers, err := db.GetDevelopers(num, start)
		if err != nil {
			fmt.Println("Error querying database: ", err.Error())
			writeErr(w, mime, http.StatusInternalServerError, "internal_error", "An internal error occurred")
			return
		}

		util.WriteJSON(w, developers)

	} else {
		writeErr(w, mime, http.StatusBadRequest, "bad_method", "You must POST or GET this endpoint!")
	}

}

func devEndpoint(w http.ResponseWriter, r *http.Request) {
	mime := r.Header.Get("Accept")
	if r.Method == "POST" || r.Method == "GET" {
		if _, ok := supportedMimes[mime]; !ok {
			writeErr(w, mime, http.StatusNotAcceptable, "not_acceptable", "This API only supports JSON at the moment.")
			return
		}

		split := strings.SplitN(r.URL.Path, "/", 6)

		if len(split) < 4 {
			writeErr(w, mime, http.StatusBadRequest, "bad_app", "Bad app slashes specified")
			return
		}

		devID := split[3]

		num, err := strconv.Atoi(devID)

		if err != nil {
			writeErr(w, mime, http.StatusBadRequest, "bad_form", "num value must be a number")
			return
		}

		dev, err := db.GetDeveloper(int64(num))

		if err != nil {
			fmt.Println("Error querying database: ", err.Error())
			writeErr(w, mime, http.StatusInternalServerError, "internal_error", "An internal error occurred")
			return
		}

		util.WriteJSON(w, dev)
	}

}

func latestsEndpoint(w http.ResponseWriter, r *http.Request) {
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

		latestApps, err := db.GetDevelopers(num, start)
		if err != nil {
			fmt.Println("Error querying database: ", err.Error())
			writeErr(w, mime, http.StatusInternalServerError, "internal_error", "An internal error occurred")
			return
		}

		util.WriteJSON(w, latestApps)

	} else {
		writeErr(w, mime, http.StatusBadRequest, "bad_method", "You must POST or GET this endpoint!")
	}
}

func searchAppsEndpoint(w http.ResponseWriter, r *http.Request) {
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

		searchTerm := split[4]
		fmt.Println("Fetching matches for:", searchTerm)

		results, err := db.SearchApps(string(searchTerm))

		fmt.Println("This many apps found: " + fmt.Sprint(len(results)))

		if err != nil {
			fmt.Printf("Error running search for %s: %s\n", searchTerm, err)
			writeErr(w, mime, http.StatusBadRequest, "bad_search", "No apps could not be found")
			return
		}

		w.Header().Set("Access-Control-Allow-Origin", "*")

		util.WriteJSON(w, results)

	}
}

func searchCompEndpoint(w http.ResponseWriter, r *http.Request) {

}

func searchDevEndpoint(w http.ResponseWriter, r *http.Request) {

}

var cfgFile = flag.String("cfg", "/etc/xray/config.json", "config file location")
var port = flag.Uint("port", 8118, "Port to serve on.")

func init() {
	util.LoadCfg(*cfgFile, util.APIServ)
	db.Open(util.Cfg, true)
}

func main() {
	http.HandleFunc("/", hello)
	http.HandleFunc("/api/apps", appsEndpoint)
	http.HandleFunc("/api/apps/", appEndpoint)
	http.HandleFunc("/api/developers", devsEndpoint)
	http.HandleFunc("/api/developers/", devEndpoint)
	http.HandleFunc("/api/companies", compsEndpoint)
	http.HandleFunc("/api/companies/", compEndpoint)
	http.HandleFunc("/api/latest", latestsEndpoint)
	http.HandleFunc("/api/search/apps/", searchAppsEndpoint)
	http.HandleFunc("/api/search/apps", searchAppsEndpoint)
	http.HandleFunc("/api/search/companies/", searchCompEndpoint)
	http.HandleFunc("/api/search/companies", searchCompEndpoint)
	http.HandleFunc("/api/search/developers/", searchDevEndpoint)
	http.HandleFunc("/api/search/developers", searchDevEndpoint)
	//http.HandleFunc("/api/latest/", latestEndpoint)
	panic(http.ListenAndServe(fmt.Sprintf(":%d", *port), nil))
}
