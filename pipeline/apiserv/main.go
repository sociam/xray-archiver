package main

import (
	"flag"
	"fmt"
	"log"
	"net/http"
	"os/exec"
	"regexp"
	"strconv"
	"strings"
	"sync"

	// "github.com/sociam/xray-archiver/pipeline/db"
	"../db"
	
	// "github.com/sociam/xray-archiver/pipeline/util"
	"../util"
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
		util.Log.Err("error writing data", err1)
	}
}

// //yeah, a map of string to some data is fine but the array setup? check for empty? or null?
// func writeMultiSuccessData(w http.ResponseWriter, mime string, data map[string][]interface{}) {
// 	w.WriteHeader(http.StatusMultiStatus)
// 	var err1 error
// 	w.Header().Set("Content-Type", mime)
// 	switch mime {
// 	case "application/nahmate":
// 		err1 = util.WriteDEAN(w, data)
// 	case "text/plain":
// 		_, err1 = w.Write(toBytes(data))
// 	default:
// 		w.Header().Set("Content-Type", "application/json")
// 		fallthrough
// 	case "application/json":
// 		//if array empty failed the lookup... can not have a domain without a host
// 		// for _, i := range data {
// 		// 	if len(i) > 1
// 		// 		response := fmt.print(data, http.StatusOK)
// 		// 		err1 = util.WriteJSON(w, response) //Need to append into data key... still want key on outside
// 		// 	else
// 		// 		err1 = util.WriteJSON(w,data)
// 		// }
// 	}
// 	if err1 != nil {
// 		util.Log.Err("error writing data", err1)
// 	}

// }

func mimeCheck(mime string) string {
	mimes := strings.Split(mime, ",")
	for _, mime := range mimes {
		mime = strings.TrimSpace(mime)
		if _, ok := supportedMimes[mime]; ok {
			return mime
		}
	}
	return ""
}

func hello(w http.ResponseWriter, r *http.Request) {
	util.Log.Warning("Got spurious request on " + r.URL.Path)
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
		return num, "Limit too high. Please slow down. Chunk the request using the offset", nil
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

func genreHostAvgEndpoint(w http.ResponseWriter, r *http.Request) {
	mime := r.Header.Get("Accept")
	w.Header().Set("Access-Control-Allow-Origin", "*")
	util.Log.Debug("Genre Host Requst...")
	//Check input
	if r.Method == "POST" || r.Method == "GET" {
		mime = mimeCheck(mime)
		if mime == "" {
			writeErr(w, mime, http.StatusNotAcceptable, "not_acceptable", "This API only supports JSON at the moment.")
			return
		}
		//alts, err := db.GetAltApps(appID)
		alts, err := db.GetGenreHostAverages()
		if err != nil {
			writeErr(w, mime, http.StatusBadRequest, "Fail", "For some reason we couldn't fetch the stats table you requested.")
			return
		}
		writeData(w, mime, http.StatusOK, alts)
	}
}

func companyGenreCoverageEndpoint(w http.ResponseWriter, r *http.Request) {
	mime := r.Header.Get("Accept")
	w.Header().Set("Access-Control-Allow-Origin", "*")
	util.Log.Debug("Genre Host Requst...")
	//Check input
	if r.Method == "POST" || r.Method == "GET" {
		mime = mimeCheck(mime)
		if mime == "" {
			writeErr(w, mime, http.StatusNotAcceptable, "not_acceptable", "This API only supports JSON at the moment.")
			return
		}
		//alts, err := db.GetAltApps(appID)
		alts, err := db.GetCompanyGenreCoverage()
		if err != nil {
			writeErr(w, mime, http.StatusBadRequest, "Fail", "For some reason we couldn't fetch the stats table you requested.")
			return
		}
		writeData(w, mime, http.StatusOK, alts)
	}
}

func appCompanyFreqEndpoint(w http.ResponseWriter, r *http.Request) {
	mime := r.Header.Get("Accept")
	w.Header().Set("Access-Control-Allow-Origin", "*")
	util.Log.Debug("Genre Host Requst...")
	//Check input
	if r.Method == "POST" || r.Method == "GET" {
		mime = mimeCheck(mime)
		if mime == "" {
			writeErr(w, mime, http.StatusNotAcceptable, "not_acceptable", "This API only supports JSON at the moment.")
			return
		}
		//alts, err := db.GetAltApps(appID)
		alts, err := db.GetAppCompanyFreq()
		if err != nil {
			writeErr(w, mime, http.StatusBadRequest, "Fail", "For some reason we couldn't fetch the stats table you requested.")
			return
		}
		writeData(w, mime, http.StatusOK, alts)
	}
}

func appTypeFreqEndpoint(w http.ResponseWriter, r *http.Request) {
	mime := r.Header.Get("Accept")
	w.Header().Set("Access-Control-Allow-Origin", "*")
	util.Log.Debug("Genre Host Requst...")
	//Check input
	if r.Method == "POST" || r.Method == "GET" {
		mime = mimeCheck(mime)
		if mime == "" {
			writeErr(w, mime, http.StatusNotAcceptable, "not_acceptable", "This API only supports JSON at the moment.")
			return
		}
		//alts, err := db.GetAltApps(appID)
		alts, err := db.GetAppTypeFreq()
		if err != nil {
			writeErr(w, mime, http.StatusBadRequest, "Fail", "For some reason we couldn't fetch the stats table you requested.")
			return
		}

		writeData(w, mime, http.StatusOK, alts)
	}
}

func fetchIDEndpoint(w http.ResponseWriter, r *http.Request) {
	mime := r.Header.Get("Accept")
	w.Header().Set("Access-Control-Allow-Origin", "*")

	//Check input
	if r.Method == "POST" || r.Method == "GET" {
		mime = mimeCheck(mime)
		if mime == "" {
			writeErr(w, mime, http.StatusNotAcceptable, "not_acceptable", "This API only supports JSON at the moment.")
			return
		}

		err := r.ParseForm()
		if err != nil {
			writeErr(w, mime, http.StatusBadRequest, "bad_form", "Error parsing form input: %s", err.Error())
			return
		}

		util.Log.Debug("Parsing Form Params.")

		for name, val := range r.Form {
			switch name {

			case "appID":
				util.Log.Debug("appID form param found.")
				util.Log.Debug("Value of appID: %s", val)

				_, err := exec.Command("node", "../archiver/retriever/idFetch.js", val[0]).Output()
				//out, err := exec.Command("ls", "/var/xray/pipeline/archiver/retriever/").Output()

				if err != nil {
					fmt.Printf("%s\n\n", err)
					writeErr(w, mime, http.StatusConflict, "conflict", "{'appID':'"+val[0]+"', 'status':'Couldnt fetch or app already exists'}")
					return
				}
				//outStr := string(out[:])
				//fmt.Printf("%s\n\n", out)
				//fmt.Println(outStr)
				//fmt.Printf("%s\n\n", out)
				writeData(w, mime, http.StatusOK, "{'appID':'"+val[0]+"', 'status':'success'}")
				//wg := new(sync.WaitGroup)
				//exeCmd("ls /var/xray/pipeline/archiver/retriever/ ", wg)
			}
		}
	}
}

func appsEndpoint(w http.ResponseWriter, r *http.Request) {
	mime := r.Header.Get("Accept")
	w.Header().Set("Access-Control-Allow-Origin", "*")

	//Check input
	if r.Method == "POST" || r.Method == "GET" {
		mime = mimeCheck(mime)
		if mime == "" {
			writeErr(w, mime, http.StatusNotAcceptable, "not_acceptable", "This API only supports JSON at the moment.")
			return
		}

		//XXX:Assuming all endpoints have a form to process...
		err := r.ParseForm()
		if err != nil {
			writeErr(w, mime, http.StatusBadRequest, "bad_form", "Error parsing form input: %s", err.Error())
			return
		}

		//Default apps

		limit := "10"

		offset := "0"
		isFull := false
		onlyAnalyzed := false //Default is true as most desire is for analyzed apps
		store := "play"

		titles := []string{}
		developers := []string{}
		genres := []string{}
		permissions := []string{}
		appIDs := []string{}
		startsWith := []string{}

		util.Log.Info("Parsing app form parameters, params size %s", fmt.Sprint(len(r.Form)))

		for name, val := range r.Form {
			oops := ""

			switch name {
			case "limit":
				util.Log.Debug("Got range of limits", val)
				limit, oops, _ = parseLimit(val[0])
				util.Log.Debug("Limit Value = ", limit)
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

			case "onlyAnalyzed":
				var err error
				onlyAnalyzed, err = strconv.ParseBool(val[0])
				if err != nil {
					writeErr(w, mime, http.StatusBadRequest, "bad_form", "onlyAnalyzed needs to be a boolean value, true or false")
					return
				}

			case "title":
				util.Log.Debug("titles:", len(val))
				titles = val

			case "startsWith":
				fmt.Println(val)
				startsWith = val

			case "developer":
				developers = val

			case "genre":
				genres = val
				//Valid genre constant check

			case "appId":
				appIDs = val

			case "nocache":
				util.Log.Debug("No Cache flag. Good Stuff.")

			default:
				writeErr(w, mime, http.StatusBadRequest, "bad_form", "passed form values did not match params", name)
				return
			}
		}

		if len(titles) == 0 {
			developers = []string{}
		}

		util.Log.Debug("Gathering full details")

		results, err := db.QueryAll(onlyAnalyzed, store, limit, offset, developers, genres, permissions, appIDs, titles, startsWith)

		util.Log.Debug("Results found %v", results)
		if err != nil {
			util.Log.Err("Error querying database: ", err.Error())
			writeErr(w, mime, http.StatusInternalServerError, "internal_error", "An internal error occurred")
			return
		}

		if !isFull {
			stubs := make([]db.AppStub, len(results), len(results))
			for i, result := range results {
				util.Log.Info(result.App)
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
		mime = mimeCheck(mime)
		if mime == "" {
			writeErr(w, mime, http.StatusNotAcceptable, "not_acceptable", "This API only supports JSON at the moment.")
			return
		}

		split := strings.Split(r.URL.Path, "/")

		if len(split) < 3 {
			writeErr(w, mime, http.StatusBadRequest, "bad_app", "Bad app slashes specified")
			return
		}

		appID := split[3]

		//alts, err := db.GetAltApps(appID)
		alts, err := db.GetManualAltApps(appID)
		if err != nil {
			writeErr(w, mime, http.StatusBadRequest, "bad_app", "Seems like we couldn't find your app... Probs means that we don't have any alts: "+appID)
			return
		}

		writeData(w, mime, http.StatusOK, alts)
	}
}

// func getHostLoci(w http.ResponseWriter, r *http.Request) {
// 	util.Log.Debug("Host lookup requested")
// 	mime := r.Header.Get("Accept")
// 	if r.Method == "POST" || r.Method == "GET" {
// 		mime = mimeCheck(mime)
// 		if mime == "" {
// 			writeErr(w, mime, http.StatusNotAcceptable, "not_acceptable", "This API only supports JSON at the moment.")
// 			return
// 		}

// 		split := strings.Split(r.URL.Path, "/")

// 		if len(split) < 3 {
// 			writeErr(w, mime, http.StatusBadRequest, "bad_hosts", "Bad app slashes specified")
// 			return
// 		}

// 		hostname := split[3]
// 		util.Log.Debug("Attempting to lookup hosts: ", hostname)
// 		geoip, err := util.GetHostGeoIP(hostname)

// 		if err != nil {
// 			writeErr(w, mime, http.StatusBadRequest, "bad_host", "the host could not be retrieved", err)
// 			return
// 		}

// 		writeData(w, mime, http.StatusOK, geoip)
// 	}
// }

func fetchHosts(w http.ResponseWriter, r *http.Request) {
	mime := r.Header.Get("Accept")
	w.Header().Set("Access-Control-Allow-Origin", "*")

	if r.Method == "POST" || r.Method == "GET" {
		mime = mimeCheck(mime)
		if mime == "" {
			writeErr(w, mime, http.StatusNotAcceptable, "not_acceptable", "This API only supports JSON at the moment.")
			return
		}

		util.Log.Debug("Parsing Form Params.")
		//need to take a hosts list of []string
		err := r.ParseForm()
		if err != nil {
			writeErr(w, mime, http.StatusBadRequest, "bad_form", "Error parsing form input: %s", err.Error())
			return
		}

		hostsParams := r.Form["hosts"][0]
		util.Log.Debug("Raw host params: %s\n", hostsParams)

		hosts := strings.Split(hostsParams, ",")
		util.Log.Debug("Checking over hosts: %s\n", hosts)

		hostToGeoip := map[string][]util.GeoIPInfo{}

		wg := sync.WaitGroup{}
		for i := range hosts {
			j := i
			util.Log.Debug("Getting host geo ip: %s\n", hosts[i])
			wg.Add(1)
			go func() {
				var geoip []util.GeoIPInfo

				geoip, err = util.GetHostGeoIP(util.Cfg.GeoIPEndpoint, hosts[j])

				if err != nil {
					// TODO: immedoiately fail? change status to accepted 202 and 200 and
					// BADREQUEST when all is well with all hosts.

					// immediately failing is impossible with parallelization (or very
					// hard) and I don't think we should use http statuses in a non-standard way -sauyon

					// writeErr(w, mime, http.StatusBadRequest, "bad_host", "the host could not be retrieved", err)
					util.Log.Notice("Host %s could not be found: %s", hosts[j], err.Error())
					hostToGeoip[hosts[j]] = nil
				} else {
					hostToGeoip[hosts[j]] = geoip
				}
				wg.Done()
			}()
		}

		wg.Wait()

		writeData(w, mime, http.StatusOK, hostToGeoip)
	}

}

var cfgFile = flag.String("cfg", "/etc/xray/config.json", "config file location")
var port = flag.Uint("port", 8118, "Port to serve on.")

func init() {
	var err error
	err = util.LoadCfg(*cfgFile, util.APIServ)
	if err != nil {
		log.Fatalf("Failed to read config: %s", err.Error())
	}
	err = db.Open(util.Cfg, true)
	if err != nil {
		log.Fatalf("Failed to open a connection to the database: %s", err.Error())
	}
}

func main() {

	http.Handle("/", http.FileServer(http.Dir(util.Cfg.AppDir)))

	http.HandleFunc("/api/apps", appsEndpoint)
	http.HandleFunc("/api/alt/", altAppsEndpoint)
	http.HandleFunc("/api/fetch", fetchIDEndpoint)
	http.HandleFunc("/api/stats/genre_host_averages", genreHostAvgEndpoint)
	http.HandleFunc("/api/stats/app_company_freq", appCompanyFreqEndpoint)
	http.HandleFunc("/api/stats/app_type_freq", appTypeFreqEndpoint)
	http.HandleFunc("/api/stats/company_genre_coverage", companyGenreCoverageEndpoint)
	http.HandleFunc("/api/hosts", fetchHosts)
	panic(http.ListenAndServe(fmt.Sprintf(":%d", *port), nil))
}
