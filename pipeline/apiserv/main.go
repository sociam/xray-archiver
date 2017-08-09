package main

import (
	"flag"
	"fmt"
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
		fmt.Println(err1)
	}
}

func hello(w http.ResponseWriter, r *http.Request) {
	fmt.Println("Got spurious request on " + r.URL.Path)
	writeErr(w, r.Header.Get("Accept"), http.StatusNotFound, "not_found", "Nah mate!")
}

var appPrefixRe = regexp.MustCompile("^/api/apps/")
var dbIDRe = regexp.MustCompile("^\\d+$")
var appIDRe = regexp.MustCompile("^[[:alpha:]][\\w$]*(\\.[[:alpha:]][\\w$]*)*$")

func parseNumCheck(num string) (val int, oops string, err error) {
	val, err = strconv.Atoi(num)

	if err != nil {
		return 0, "num value must be a number", nil
	}

	if val < 0 {
		return 0, "num can not be a value less than 0...", nil
	}

	return val, "", nil
}

func parseLimit(num string) (val string, oops string, err error) {

	if len(val) > 1 {
		return "", "num must have a single value", nil
	}
	realNum := 0

	realNum, oops, err = parseNumCheck(num)

	if oops != "" {
		return num, oops, nil
	}

	if realNum > 1000000 {
		return num, "Limit to high. Please slow down. Chunk the request using the offset", nil
	}

	return num, "", err
}

func parseOffset(num string) (val string, oops string, err error) {
	realNum := 0

	realNum, oops, err = parseNumCheck(num)

	if oops != "" {
		return num, oops, nil
	}

	if err != nil {
		return num, "offset value must be a number", nil
	}
	if realNum < 0 {
		return num, "offset value must positive", nil
	}

	return num, "", err
}

func appsEndpoint(w http.ResponseWriter, r *http.Request) {
	mime := r.Header.Get("Accept")

	// DONT DELETE DEAN.
	w.Header().Set("Access-Control-Allow-Origin", "*")
	// DEAN... DONT DELETE.

	//Check input
	if r.Method == "POST" || r.Method == "GET" {
		if _, ok := supportedMimes[mime]; !ok {
			writeErr(w, mime, http.StatusNotAcceptable, "not_acceptable", "This API only supports JSON at the moment.")
			return
		}

		//XXX:Assuming all endpoints have a form to process...
		err := r.ParseForm()
		if err != nil {
			writeErr(w, mime, http.StatusBadRequest, "bad_form", "Error parsing form input: %s", err.Error())
			return
		}

		mime := r.Header.Get("Accept")
		//Default apps

		limit := "10"

		offset := "0"
		isFull := false

		titles := []string{""}
		developers := []string{""}
		genres := []string{""}
		permissions := []string{""}
		appIDs := []string{""}

		fmt.Printf("Parsing app form parameters, params size %s", fmt.Sprint(len(r.Form)))
		//Should not complain if form is 0...

		for name, val := range r.Form {
			oops := ""

			switch name {
			case "limit":
				fmt.Println("Got range of limits", val)
				limit, oops, _ = parseLimit(val[0])
				fmt.Println("Limit Value = ", limit)
				if oops != "" {
					writeErr(w, mime, http.StatusBadRequest, "bad_form", oops)
					return
				}

			case "offset":
				offset, oops, _ = parseOffset(val[0])
				if oops != "" {
					writeErr(w, mime, http.StatusBadRequest, "bad_form", oops)
					return
				}

			case "isFull":
				var err error
				isFull, err = strconv.ParseBool(val[0])
				if err != nil {
					writeErr(w, mime, http.StatusBadRequest, "bad_form", "isFull needs to be a boolean value, true or false")
					return
				}

			case "title":
				fmt.Println("titles:", len(val))
				titles = val

			case "developer":
				developers = val

			case "genre":
				genres = val
				//Valid genre constant check

			case "appId":
				appIDs = val

			default:
				writeErr(w, mime, http.StatusBadRequest, "bad_form", "passed form values did not match params", name)
				return
			}
		}

		if len(titles) == 0 {
			developers = []string{}
		}

		fmt.Println("Gathering full details")

		results, err := db.QuickQuery(isFull, "playstore_apps", limit, offset, developers, genres, permissions, appIDs, titles)

		if err != nil {
			fmt.Println("Error querying database: ", err.Error())
			writeErr(w, mime, http.StatusInternalServerError, "internal_error", "An internal error occurred")
			return
		}

		if !isFull {
			stubs := make([]db.AppStub, len(results), len(results))
			for i, result := range results {
				fmt.Println(result.App)
				stubs[i].Title = result.StoreInfo.(db.PlayStoreInfo).Title
				stubs[i].App = result.App
			}

			writeData(w, mime, http.StatusOK, stubs)
		} else {
			writeData(w, mime, http.StatusOK, results)
		}

	} else {
		writeErr(w, mime, http.StatusBadRequest, "bad_method", "You must POST or GET this endpoint!")
	}

}

// altAppsEndpoint allows for external entities to query for alternative apps based on app ID.
func altAppsEndpoint(w http.ResponseWriter, r *http.Request) {
	mime := r.Header.Get("Accept")
	w.Header().Set("Access-Control-Allow-Origin", "*")
	if r.Method == "POST" || r.Method == "GET" {
		if _, ok := supportedMimes[mime]; !ok {
			writeErr(w, mime, http.StatusNotAcceptable, "not_acceptable", "Yo Dawg, we deal with json only son.")
			return
		}

		split := strings.Split(r.URL.Path, "/")

		if len(split) < 3 {
			writeErr(w, mime, http.StatusBadRequest, "bad_app", "Bad app slashes specified")
			return
		}

		appID := split[3]

		alts, err := db.GetAltApps(appID)

		if err != nil {
			writeErr(w, mime, http.StatusBadRequest, "bad_app", "Seems like we couldn't find your app... Probs means that we don't have any alts")
			return
		}

		writeData(w, mime, http.StatusOK, alts)
	}
}

var cfgFile = flag.String("cfg", "/etc/xray/config.json", "config file location")
var port = flag.Uint("port", 8118, "Port to serve on.")

func init() {
	util.LoadCfg(*cfgFile, util.APIServ)
	db.Open(util.Cfg, true)
}

func main() {
	http.Handle("/", http.FileServer(http.Dir(util.Cfg.AppDir)))

	http.HandleFunc("/api/apps/", appsEndpoint)
	http.HandleFunc("/api/alt/", altAppsEndpoint)

	panic(http.ListenAndServe(fmt.Sprintf(":%d", *port), nil))
}
