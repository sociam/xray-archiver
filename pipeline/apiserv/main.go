package main

import (
	"flag"
	"fmt"
	"log"
	"net/http"
	"regexp"
	"strconv"

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


func parseNumCheck(num string) (val int, oops string, err error) {
	//oops error
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

// func validGenre(gen string ) (val string, error) {
// 	//TODO: Check against genre constnats
// 	return gen
// }

func gatherAppsEndpoint(w http.ResponseWriter, r *http.Request) {
	mime := r.Header.Get("Accept")
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

		limit := db.FormParam{"limit", "10"}

		offset := db.FormParam{"offset", "0"}
		isFull := true

		formParams := make([]db.FormParam, 0, 10) //:= make([]db.formParams, 3, 7)

		titles := []string{""}
		developers := []string{""}
		genres := []string{""}
		permisions := []string{""}
		appIDs := []string{""}

		fmt.Printf("Parsing app form parameters, params size %s", fmt.Sprint(len(r.Form)))
		//Should not complain if form is 0...

		for name, val := range r.Form {
			oops := ""

			switch name {
			case limit.Name:
				fmt.Println("Got range of limits", val)
				limit.Val, oops, _ = parseLimit(val[0])
				if oops != "" {
					writeErr(w, mime, http.StatusBadRequest, "bad_form", oops)
					return
				}
				fmt.Println("Parsing app form parameters ")

			case offset.Name:
				offset.Val, oops, _ = parseOffset(val[0])
				if oops != "" {
					writeErr(w, mime, http.StatusBadRequest, "bad_form", oops)
					return
				}

			case "isFull":
				_, err := strconv.ParseBool(val[0])
				if err != nil {
					writeErr(w, mime, http.StatusBadRequest, "bad_form", "isFull needs to be a boolean value, true or false")
					return
				}

			case "title":
				fmt.Println("titles:", len(val))
				titles = val
				// form := db.FormParam{"title", val[0]}
				// formParams = append(formParams, form)

			case "developer":
				developers = val
				//formParams = append(formParams, db.FormParam{"developer", val[0]})

			case "genre":
				genres = val
				//formParams = append(formParams, db.FormParam{"genre", val[0]})
				//Valid genre check

			case "appId":
				appIDs = val
				//formParams = append(formParams, db.FormParam{"appId", val[0]})

			default:
				writeErr(w, mime, http.StatusBadRequest, "bad_form", "passed form values did not match params", name)
				return
			}
		}
		println("Finished passing parameters, formParams", len(formParams))

		if len(titles) == 0 {
			developers = []string{}
		}

		if isFull {
			//TODO: pass store paramters
			fmt.Println("Gather full details")

			// string  int    int     string[] string[]..
			results, err := db.QuickQuery("playstore_apps", limit.Val, offset.Val, developers, genres, permisions, appIDs, titles)

			//results, err := db.RetrieveFullFrom("playstore_apps", limit.Val, offset.Val, formParams) // titles /*=[a, b, c]*/, developers /*=[1, 2, 3]*/)

			if err != nil {
				fmt.Println("Error querying database: ", err.Error())
				writeErr(w, mime, http.StatusInternalServerError, "internal_error", "An internal error occurred")
				return
			}

			util.WriteJSON(w, results)

		} else {
			//TODO: non full
		}
	} else {
		writeErr(w, mime, http.StatusBadRequest, "bad_method", "You must POST or GET this endpoint!")
	}

}

var cfgFile = flag.String("cfg", "/etc/xray/config.json", "config file location")
var port = flag.Uint("port", 8118, "Port to serve on.")

func init() {
	util.LoadCfg(*cfgFile, util.APIServ)
	db.Open(util.Cfg, true)
}

func main() {
	http.HandleFunc("/", hello)

	http.HandleFunc("/api/apps/", gatherAppsEndpoint)
	//@deprecated
	// http.HandleFunc("/api/apps", appsEndpoint) //Returned chunked apps
	// http.HandleFunc("/api/apps/", appEndpoint) //?full=True&title=Title&developer=dev&genre=GENRELIST&permisions=PERMISSIONLIST&appId=id
	// http.HandleFunc("/api/developers", devsEndpoint)
	// http.HandleFunc("/api/developers/", devEndpoint)
	// http.HandleFunc("/api/companies", compsEndpoint)
	// http.HandleFunc("/api/companies/", compEndpoint)
	// http.HandleFunc("/api/latest", latestsEndpoint)
	// http.HandleFunc("/api/search/apps/", searchAppsEndpoint)
	// http.HandleFunc("/api/search/apps", searchAppsEndpoint)

	//TODO:
	//http.HandleFunc("/api/latest/", latestEndpoint)

	panic(http.ListenAndServe(fmt.Sprintf(":%d", *port), nil))
}
