package main

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
	"reflect"
	"testing"
	"time"

	"github.com/sociam/xray-archiver/pipeline/db"
)

type xrayDb struct {
	*sql.DB
}

var useDB bool
var database xrayDb

var ExampleApp = db.App{ID: "example", Vers: []int64{1.0}}

var ExamplePlaystore = db.PlayStoreInfo{
	Title:         "EXAMPLE TITLE",
	Summary:       "I AM A EXAMPLE",
	Description:   "I WISH TO EXAMPLE",
	StoreURL:      "/path/to/store_url",
	Price:         "Free",
	Free:          true,
	Rating:        "3.7",
	NumReviews:    1,
	Genre:         "GAME_CASUAL",
	FamilyGenre:   "",
	Installs:      db.Range{Min: int64(1), Max: int64(10)},
	Developer:     1,
	Updated:       time.Date(1996, 1, 1, 1, 1, 1, 1, time.Local),
	AndroidVer:    "4.1",
	Screenshots:   []string{},
	Video:         "",
	RecentChanges: []string{},
	CrawlDate:     time.Date(1996, 1, 1, 1, 1, 1, 2, time.Local),
	Permissions:   []string{},
}

var ExampleDev = db.Developer{
	Emails:    []string{"example@email.com"},
	Name:      "example_name",
	StoreSite: "/path/to/store_site",
	Site:      "path/to/site"}

var ExampleAppVer = db.AppVersion{ID: 1,
	App:         ExampleApp.ID,
	Store:       "play",
	Region:      "uk",
	Ver:         ExamplePlaystore.AndroidVer, //Urgh maybe not android ver might of just screwed up versions
	ScreenFlags: 0,
	StoreInfo:   ExamplePlaystore,
	Icon:        "",
	Dev:         ExampleDev,
	Hosts:       []string{"example.com"},
	Perms:       []string{},
	Packages:    []string{},
	IsAnalyzed:  true}

//var cfgFileTest = flag.String("cfg", "/etc/xray/config.json", "config file location")

//var portTest = flag.Uint("port", 8118, "Port to serve on.")

// func Test_Soon_Init(t *testing.T) {
// 	var err error
// 	err = util.LoadCfg(*cfgFileTest, util.APIServ) //TODO:might break all setup
// 	if err != nil {
// 		t.Errorf("Could not load and setup database config", err)
// 	}

// 	err = db.Open(util.Cfg, true)
// 	if err != nil {
// 		t.Errorf("Could not connect to database", err)
// 	}

// 	//ID should be 1..
// 	rows, err := database.Query("INSERT INTO apps VALUES ($1, $2)",
// 		ExampleApp.ID, ExampleApp.Vers[0])

// 	if rows != nil {
// 		rows.Close()
// 	}

// 	if err != nil {
// 		t.Errorf("Could not add exampel host", err)
// 	}

// 	// xrayDb.Query("INSERT INTO app_versions(app, store, region, version, downloaded, last_dl_attempt, analyzed) VALUES ($1, $2, $3, $4, $5, $6, $7) RETURNING id",
// 	// 	ExampleApp.ID,
// 	// 	ExampleAppVer.Store,
// 	// 	ExampleAppVer.Region,
// 	// 	ExampleAppVer.Ver,
// 	// 	0,
// 	// 	"epoch",
// 	// 	0)

// 	if rows != nil {
// 		rows.Close()
// 	}

// 	if err != nil {
// 		t.Errorf("Could not add appversions example data", err)
// 	}

// 	// rows, err := xrayDb.Query("INSERT INTO playstore_apps VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, $13, $14, $15, $16, $17, $18, $19, $20, $21)",
// 	// 	ExamplePlaystore.AndroidVer,
// 	// 	ExamplePlaystore.Title,
// 	// 	ExamplePlaystore.Summary,
// 	// 	ExamplePlaystore.Description,
// 	// 	ExamplePlaystore.StoreURL,
// 	// 	ExamplePlaystore.Price,
// 	// 	ExamplePlaystore.Free,
// 	// 	ExamplePlaystore.Rating,
// 	// 	ExamplePlaystore.NumReviews,
// 	// 	ExamplePlaystore.Genre,
// 	// 	ExamplePlaystore.FamilyGenre,
// 	// 	ExamplePlaystore.Installs.Min,
// 	// 	ExamplePlaystore.Installs.Max,
// 	// 	ExamplePlaystore.Developer,
// 	// 	ExamplePlaystore.Updated,
// 	// 	ExamplePlaystore.AndroidVer,
// 	// 	ExamplePlaystore.ContentRating,
// 	// 	ExamplePlaystore.Screenshots,
// 	// 	ExamplePlaystore.Video,
// 	// 	ExamplePlaystore.RecentChanges,
// 	// 	ExamplePlaystore.Updated)

// 	// if rows != nil {
// 	// 	rows.Close()
// 	// }

// 	// if err != nil {
// 	// 	t.Errorf("Could not add playstore_app  example data", err)
// 	// }

// 	// rows, err := xrayDb.Query("SELECT * FROM apps WHERE id = $1", ExampleApp.ID)
// 	// if rows != nil {
// 	// 	rows.Close()
// 	// }

// 	// if err != nil {
// 	// 	t.Errorf("Could not select  example data", err)
// 	// }

// 	// if rows != 1 {
// 	// 	t.Errorf("Correct number of example row data was not added", rows)
// 	// }

// }

/*API ENDPOINT TEST*/
//Test the app endppoint or the query? the problem is main barely have shiz...
//just querying the real running endpoint?

var apiUrl = "http://localhost:8118"

var myClient = &http.Client{Timeout: 10 * time.Second}

func getJSON(response []byte) (data []map[string]interface{}, err error) {

	err = json.Unmarshal(response, &data)

	return data, err
}

func Test_AccessEndpoint(t *testing.T) {
	// client := &http.Client{
	// 	CheckRedirect: redirectPolicyFunc,
	// }
	resp, err := http.Get(apiUrl)

	if err != nil {
		t.Errorf("Could not access endpoint", err)
	}

	defer resp.Body.Close()
	//body, err := ioutil.ReadAll(resp.Body)
}

func constructAPIURL(resource string, formData url.Values) string {
	u, _ := url.ParseRequestURI(apiUrl)
	u.Path = resource
	if formData != nil {
		u.RawQuery = formData.Encode()
	}
	return fmt.Sprintf("%v", u)
}

func Test_DefaultLimit(t *testing.T) {
	urlStr := constructAPIURL("/api/apps", nil)

	hc := http.Client{}

	req, err := http.NewRequest("POST", urlStr, nil)
	if err != nil {
		t.Errorf("Bad post", err)
	}

	req.Header.Add("Accept", "application/json") //TODO: generate func to pass to req

	resp, err := hc.Do(req)

	if err != nil {
		t.Errorf("Could not access endpoint", err)
	}

	if resp.StatusCode != http.StatusOK {
		t.Errorf("Bad request", resp.Status)
	}

	respData, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		t.Errorf("Umable to read response", err)
	}

	data, err := getJSON(respData)

	if err != nil {
		t.Errorf("Unable to retrieve response data", err)
	}

	if len(data) > 10 {
		t.Errorf("Default size returned greater than 10:", data)
	}

	if _, ok := data[0]["title"]; !ok {
		t.Errorf("No title present")
	}

	if _, ok := data[0]["app"]; !ok {
		t.Error("No app present")
	}
}

func Test_IsFulParam(t *testing.T) {
	// ("/api/apps?isFull=True")
	form := url.Values{}
	form.Add("isFull", "true")

	urlStr := constructAPIURL("/api/apps", form)

	hc := http.Client{}

	req, err := http.NewRequest("POST", urlStr, nil)
	if err != nil {
		t.Errorf("Bad post", err)
	}
	req.Header.Add("Accept", "application/json")

	resp, err := hc.Do(req)

	if err != nil {
		t.Errorf("Could not access endpoint", err)
	}

	if resp.StatusCode != http.StatusOK {
		t.Errorf("Bad request", resp.Status)
	}

	respData, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		t.Errorf("Umable to read response", err)
	}

	var data []db.AppVersion

	err = json.Unmarshal(respData, &data)

	if err != nil {
		t.Errorf("Unable to retrieve response data", err)
	}

	if data == nil {
		t.Errorf("Empty data returned")
	}

	//TODO: validate data
}

func Test_Title(t *testing.T) {
	//("/api/apps?title=DinnerAtNoon")
	form := url.Values{}
	form.Add("title", "EXAMPLE TITLE")

	urlStr := constructAPIURL("/api/apps", form)

	hc := http.Client{}

	req, err := http.NewRequest("POST", urlStr, nil)
	if err != nil {
		t.Errorf("Bad post", err)
	}
	req.Header.Add("Accept", "application/json")

	resp, err := hc.Do(req)

	if err != nil {
		t.Errorf("Could not access endpoint", err)
	}

	if resp.StatusCode != http.StatusOK {
		t.Errorf("Bad request", resp.Status)
	}

	respData, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		t.Errorf("Umable to read response", err)
	}

	data, err := getJSON(respData)

	if err != nil {
		t.Errorf("Unable to retrieve response data", err)
	}

	if len(data) != 1 {
		t.Errorf("Shoul be only one match", data)
	}

	if _, ok := data[0]["title"]; !ok {
		t.Errorf("No title present")
	}

	if _, ok := data[0]["app"]; !ok {
		t.Error("No app present")
	}
}

func Test_NonTitle(t *testing.T) {
	//("/api/apps?title=DinnerAtNoon")
	form := url.Values{}
	form.Add("title", "EXAMPLE NOT A TITLE PLEASE")

	urlStr := constructAPIURL("/api/apps", form)

	hc := http.Client{}

	req, err := http.NewRequest("POST", urlStr, nil)
	if err != nil {
		t.Errorf("Bad post", err)
	}
	req.Header.Add("Accept", "application/json")

	resp, err := hc.Do(req)

	if err != nil {
		t.Errorf("Could not access endpoint", err)
	}

	if resp.StatusCode != http.StatusOK {
		t.Errorf("Bad request", resp.Status)
	}

	respData, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		t.Errorf("Umable to read response", err)
	}

	data, err := getJSON(respData)

	if err != nil {
		t.Errorf("Unable to retrieve response data", err)
	}

	if len(data) < 1 {
		t.Errorf("Should be no matchs", data)
	}

}

func Test_MultiTitle(t *testing.T) {
	//("/api/apps?title=DinnerAtNoon,cats")
	//("/api/apps?title=dinner&title=at&title=noon")
	form := url.Values{}
	form.Add("title", "EXAMPLE TITLE")
	form.Add("title", "EXAMPLE TITLE TWO")

	urlStr := constructAPIURL("/api/apps", form)

	hc := http.Client{}

	req, err := http.NewRequest("POST", urlStr, nil)
	if err != nil {
		t.Errorf("Bad post", err)
	}
	req.Header.Add("Accept", "application/json")

	resp, err := hc.Do(req)

	if err != nil {
		t.Errorf("Could not access endpoint", err)
	}

	if resp.StatusCode != http.StatusOK {
		t.Errorf("Bad request", resp.Status)
	}

	respData, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		t.Errorf("Umable to read response", err)
	}

	data, err := getJSON(respData)

	if err != nil {
		t.Errorf("Unable to retrieve response data", err)
	}

	if len(data) != 2 {
		t.Errorf("Should be only one match", data)
	}

	if _, ok := data[0]["title"]; !ok {
		t.Errorf("No title present")
	}

	if _, ok := data[0]["app"]; !ok {
		t.Error("No app present")
	}

}

func Test_Developer(t *testing.T) {
	//("/api/apps?developer=zynga")
	form := url.Values{}
	form.Add("developer", "EXAMPLE")

	urlStr := constructAPIURL("/api/apps", form)

	hc := http.Client{}

	req, err := http.NewRequest("POST", urlStr, nil)
	if err != nil {
		t.Errorf("Bad post", err)
	}
	req.Header.Add("Accept", "application/json")

	resp, err := hc.Do(req)

	if err != nil {
		t.Errorf("Could not access endpoint", err)
	}

	if resp.StatusCode != http.StatusOK {
		t.Errorf("Bad request", resp.Status)
	}

	respData, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		t.Errorf("Umable to read response", err)
	}

	data, err := getJSON(respData)

	if err != nil {
		t.Errorf("Unable to retrieve response data", err)
	}

	if len(data) != 1 {
		t.Errorf("Should be only one match", data)
	}

	if _, ok := data[0]["title"]; !ok {
		t.Errorf("No title present")
	}

	if _, ok := data[0]["app"]; !ok {
		t.Error("No app present")
	}
}

func Test_NonDeveloper(t *testing.T) {
	//("/api/apps?title=DinnerAtNoon")
	form := url.Values{}
	form.Add("title", "EXAMPLE NOT A DEVELOPER PLEASE")

	urlStr := constructAPIURL("/api/apps", form)

	hc := http.Client{}

	req, err := http.NewRequest("POST", urlStr, nil)
	if err != nil {
		t.Errorf("Bad post", err)
	}
	req.Header.Add("Accept", "application/json")

	resp, err := hc.Do(req)

	if err != nil {
		t.Errorf("Could not access endpoint", err)
	}

	if resp.StatusCode != http.StatusOK {
		t.Errorf("Bad request", resp.Status)
	}

	respData, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		t.Errorf("Umable to read response", err)
	}

	data, err := getJSON(respData)

	if err != nil {
		t.Errorf("Unable to retrieve response data", err)
	}

	if len(data) < 1 {
		t.Errorf("Should be no matchs", data)
	}

}
func Test_MultiDevelopers(t *testing.T) {
	//TODO:("/api/apps?developer=zynga,cat")
	//("/api/apps?developer=zynga&developer=cat")
	form := url.Values{}
	form.Add("title", "EXAMPLE DEVELOPER")
	form.Add("title", "EXAMPLE DEVELOPER TWO")

	urlStr := constructAPIURL("/api/apps", form)

	hc := http.Client{}

	req, err := http.NewRequest("POST", urlStr, nil)
	if err != nil {
		t.Errorf("Bad post", err)
	}
	req.Header.Add("Accept", "application/json")

	resp, err := hc.Do(req)

	if err != nil {
		t.Errorf("Could not access endpoint", err)
	}

	if resp.StatusCode != http.StatusOK {
		t.Errorf("Bad request", resp.Status)
	}

	respData, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		t.Errorf("Umable to read response", err)
	}

	data, err := getJSON(respData)

	if err != nil {
		t.Errorf("Unable to retrieve response data", err)
	}

	if len(data) != 2 {
		t.Errorf("Should be only one match", data)
	}

	if _, ok := data[0]["title"]; !ok {
		t.Errorf("No title present")
	}

	if _, ok := data[0]["app"]; !ok {
		t.Error("No app present")
	}
}

func Test_Genre(t *testing.T) {
	//("/api/apps?genre=entertainment")
	//TODO: test every genre
	form := url.Values{}
	form.Add("genre", "EXAMPLE")

	urlStr := constructAPIURL("/api/apps", form)

	hc := http.Client{}

	req, err := http.NewRequest("POST", urlStr, nil)
	if err != nil {
		t.Errorf("Bad post", err)
	}
	req.Header.Add("Accept", "application/json")

	resp, err := hc.Do(req)

	if err != nil {
		t.Errorf("Could not access endpoint", err)
	}

	if resp.StatusCode != http.StatusOK {
		t.Errorf("Bad request", resp.Status)
	}

	respData, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		t.Errorf("Umable to read response", err)
	}

	data, err := getJSON(respData)

	if err != nil {
		t.Errorf("Unable to retrieve response data", err)
	}

	if len(data) > 10 {
		t.Errorf("Should be only one match", data)
	}

	if _, ok := data[0]["title"]; !ok {
		t.Errorf("No title present")
	}

	if _, ok := data[0]["app"]; !ok {
		t.Error("No app present")
	}

}

func Test_MultiGenre(t *testing.T) {
	//("/api/apps?genre=lifestyle,food_and_drink")
	//("/api/apps?genre=lifestyle&genre=food_and_drink")
	form := url.Values{}
	form.Add("genre", "EXAMPLE GENRE")
	form.Add("genre", "EXAMPLE GENRE TWO")

	urlStr := constructAPIURL("/api/apps", form)

	hc := http.Client{}

	req, err := http.NewRequest("POST", urlStr, nil)
	if err != nil {
		t.Errorf("Bad post", err)
	}
	req.Header.Add("Accept", "application/json")

	resp, err := hc.Do(req)

	if err != nil {
		t.Errorf("Could not access endpoint", err)
	}

	if resp.StatusCode != http.StatusOK {
		t.Errorf("Bad request", resp.Status)
	}

	respData, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		t.Errorf("Umable to read response", err)
	}

	data, err := getJSON(respData)

	if err != nil {
		t.Errorf("Unable to retrieve response data", err)
	}

	if len(data) != 2 {
		t.Errorf("Should be only one match", data)
	}

	if _, ok := data[0]["title"]; !ok {
		t.Errorf("No title present")
	}

	if _, ok := data[0]["app"]; !ok {
		t.Error("No app present")
	}
}

func Test_GameGenre(t *testing.T) {
	//("/api/apps?genre=GAMES")
	form := url.Values{}
	form.Add("genre", "GAMES")
	form.Add("isFull", "true") //Need fulll data
	form.Add("limit", "100")   //Need high limit inorder to hope two different game genres
	urlStr := constructAPIURL("/api/apps", form)

	hc := http.Client{}

	req, err := http.NewRequest("POST", urlStr, nil)
	if err != nil {
		t.Errorf("Bad post", err)
	}
	req.Header.Add("Accept", "application/json")

	resp, err := hc.Do(req)

	if err != nil {
		t.Errorf("Could not access endpoint", err)
	}

	if resp.StatusCode != http.StatusOK {
		t.Errorf("Bad request", resp.Status)
	}

	respData, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		t.Errorf("Umable to read response", err)
	}

	var data []db.AppVersion

	err = json.Unmarshal(respData, &data)

	if err != nil {
		t.Errorf("Unable to retrieve response data", err)
	}

	if data == nil {
		t.Errorf("Empty data returned")
	}

	//Test to see if a unqiue genre exists

	multiGenresFound := false
	cmpGenre := data[0].StoreInfo.(db.PlayStoreInfo).Genre
	for _, app := range data {
		appGenre := app.StoreInfo.(db.PlayStoreInfo).Genre
		if cmpGenre != appGenre {
			multiGenresFound = true
		}
	}

	if !multiGenresFound {
		t.Errorf("Games did not get mutlple unique genres, maybe increase limit size. Data:\n", data)
	}
}

func Test_AppId(t *testing.T) {
	//("/api/apps?appId")
	form := url.Values{}
	form.Add("appid", "EXAMPLE APP")

	urlStr := constructAPIURL("/api/apps", form)

	hc := http.Client{}

	req, err := http.NewRequest("POST", urlStr, nil)
	if err != nil {
		t.Errorf("Bad post", err)
	}
	req.Header.Add("Accept", "application/json")

	resp, err := hc.Do(req)

	if err != nil {
		t.Errorf("Could not access endpoint", err)
	}

	if resp.StatusCode != http.StatusOK {
		t.Errorf("Bad request", resp.Status)
	}

	respData, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		t.Errorf("Umable to read response", err)
	}

	data, err := getJSON(respData)

	if err != nil {
		t.Errorf("Unable to retrieve response data", err)
	}

	if len(data) != 1 {
		t.Errorf("Should be only one match", data)
	}

	if _, ok := data[0]["title"]; !ok {
		t.Errorf("No title present")
	}

	if _, ok := data[0]["app"]; !ok {
		t.Error("No app present")
	}

	if data[0]["app"] != "EXAMPLE APP" {
		t.Error("Incorrect appid found: ", data[0]["app"])
	}

}

func Test_MultiAppId(t *testing.T) {
	//("/api/apps?appId=example1&appId=example2")
	//("/api/apps?appId=example1,example2")
	form := url.Values{}
	form.Add("title", "EXAMPLE APP")
	form.Add("title", "EXAMPLE APP TWO")

	urlStr := constructAPIURL("/api/apps", form)

	hc := http.Client{}

	req, err := http.NewRequest("POST", urlStr, nil)
	if err != nil {
		t.Errorf("Bad post", err)
	}
	req.Header.Add("Accept", "application/json")

	resp, err := hc.Do(req)

	if err != nil {
		t.Errorf("Could not access endpoint", err)
	}

	if resp.StatusCode != http.StatusOK {
		t.Errorf("Bad request", resp.Status)
	}

	respData, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		t.Errorf("Umable to read response", err)
	}

	data, err := getJSON(respData)

	if err != nil {
		t.Errorf("Unable to retrieve response data", err)
	}

	if len(data) != 2 {
		t.Errorf("Should be only one match", data)
	}

	if _, ok := data[0]["title"]; !ok {
		t.Errorf("No title present")
	}

	if _, ok := data[0]["app"]; !ok {
		t.Error("No app present")
	}
}

func Test_invalidAppid(t *testing.T) {
	//("/api/apps?appId=INVALID")
	form := url.Values{}
	form.Add("appid", "EXAMPLE NOT A APP ID PLEASE")

	urlStr := constructAPIURL("/api/apps", form)

	hc := http.Client{}

	req, err := http.NewRequest("POST", urlStr, nil)
	if err != nil {
		t.Errorf("Bad post", err)
	}
	req.Header.Add("Accept", "application/json")

	resp, err := hc.Do(req)

	if err != nil {
		t.Errorf("Could not access endpoint", err)
	}

	if resp.StatusCode != http.StatusOK {
		t.Errorf("Bad request", resp.Status)
	}

	respData, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		t.Errorf("Umable to read response", err)
	}

	data, err := getJSON(respData)

	if err != nil {
		t.Errorf("Unable to retrieve response data", err)
	}

	if len(data) < 1 {
		t.Errorf("Should be no matchs", data)
	}

}

func Test_startsWithSingleLetter(t *testing.T) {
	//("/api/apps?startsWith=i")
	form := url.Values{}
	form.Add("startsWith", "a")

	urlStr := constructAPIURL("/api/apps", form)

	hc := http.Client{}

	req, err := http.NewRequest("POST", urlStr, nil)
	if err != nil {
		t.Errorf("Bad post", err)
	}
	req.Header.Add("Accept", "application/json")

	resp, err := hc.Do(req)

	if err != nil {
		t.Errorf("Could not access endpoint", err)
	}

	if resp.StatusCode != http.StatusOK {
		t.Errorf("Bad request", resp.Status)
	}

	respData, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		t.Errorf("Umable to read response", err)
	}

	data, err := getJSON(respData)

	if err != nil {
		t.Errorf("Unable to retrieve response data", err)
	}

	if len(data) > 10 {
		t.Errorf("Should be more thabn 10 matchs", data)
	}

	if _, ok := data[0]["title"]; !ok {
		t.Errorf("No title present")
	}

	if _, ok := data[0]["app"]; !ok {
		t.Error("No app present")
	}
}

func Test_startsWithDoubleLetter(t *testing.T) {
	//("/api/apps?startsWith=i")
	form := url.Values{}
	form.Add("startsWith", "aa")

	urlStr := constructAPIURL("/api/apps", form)

	hc := http.Client{}

	req, err := http.NewRequest("POST", urlStr, nil)
	if err != nil {
		t.Errorf("Bad post", err)
	}
	req.Header.Add("Accept", "application/json")

	resp, err := hc.Do(req)

	if err != nil {
		t.Errorf("Could not access endpoint", err)
	}

	if resp.StatusCode != http.StatusOK {
		t.Errorf("Bad request", resp.Status)
	}

	respData, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		t.Errorf("Umable to read response", err)
	}

	data, err := getJSON(respData)

	if err != nil {
		t.Errorf("Unable to retrieve response data", err)
	}

	if len(data) > 10 {
		t.Errorf("Should be more than 10 matchs", data)
	}

	if _, ok := data[0]["title"]; !ok {
		t.Errorf("No title present")
	}

	if _, ok := data[0]["app"]; !ok {
		t.Error("No app present")
	}
}

func Test_startsWithTripleLetter(t *testing.T) {
	//("/api/apps?startsWith=i")
	form := url.Values{}
	form.Add("startsWith", "cat")

	urlStr := constructAPIURL("/api/apps", form)

	hc := http.Client{}

	req, err := http.NewRequest("POST", urlStr, nil)
	if err != nil {
		t.Errorf("Bad post", err)
	}
	req.Header.Add("Accept", "application/json")

	resp, err := hc.Do(req)

	if err != nil {
		t.Errorf("Could not access endpoint", err)
	}

	if resp.StatusCode != http.StatusOK {
		t.Errorf("Bad request", resp.Status)
	}

	respData, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		t.Errorf("Umable to read response", err)
	}

	data, err := getJSON(respData)

	if err != nil {
		t.Errorf("Unable to retrieve response data", err)
	}

	if len(data) > 10 {
		t.Errorf("Should be more than 10 matchs", data)
	}

	if _, ok := data[0]["title"]; !ok {
		t.Errorf("No title present")
	}

	if _, ok := data[0]["app"]; !ok {
		t.Error("No app present")
	}
}

func Test_onlyAnalyzed(t *testing.T) {
	//("/api/apps?onlyAnalyzed")

}

func Test_limitOne(t *testing.T) {
	//("/api/apps?limit=10")
	form := url.Values{}
	form.Add("limit", "1")

	urlStr := constructAPIURL("/api/apps", form)

	hc := http.Client{}

	req, err := http.NewRequest("POST", urlStr, nil)
	if err != nil {
		t.Errorf("Bad post", err)
	}
	req.Header.Add("Accept", "application/json")

	resp, err := hc.Do(req)

	if err != nil {
		t.Errorf("Could not access endpoint", err)
	}

	if resp.StatusCode != http.StatusOK {
		t.Errorf("Bad request", resp.Status)
	}

	respData, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		t.Errorf("Umable to read response", err)
	}

	data, err := getJSON(respData)

	if err != nil {
		t.Errorf("Unable to retrieve response data", err)
	}

	if len(data) != 1 {
		t.Errorf("Should be more than 10 matchs", data)
	}

	if _, ok := data[0]["title"]; !ok {
		t.Errorf("No title present")
	}

	if _, ok := data[0]["app"]; !ok {
		t.Error("No app present")
	}
}

func Test_limitHundo(t *testing.T) {
	//("/api/apps?limit=10")
	form := url.Values{}
	form.Add("limit", "100")

	urlStr := constructAPIURL("/api/apps", form)

	hc := http.Client{}

	req, err := http.NewRequest("POST", urlStr, nil)
	if err != nil {
		t.Errorf("Bad post", err)
	}
	req.Header.Add("Accept", "application/json")

	resp, err := hc.Do(req)

	if err != nil {
		t.Errorf("Could not access endpoint", err)
	}

	if resp.StatusCode != http.StatusOK {
		t.Errorf("Bad request", resp.Status)
	}

	respData, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		t.Errorf("Umable to read response", err)
	}

	data, err := getJSON(respData)

	if err != nil {
		t.Errorf("Unable to retrieve response data", err)
	}

	if len(data) != 100 {
		t.Errorf("Should be more than 100 matchs", data)
	}

	if _, ok := data[0]["title"]; !ok {
		t.Errorf("No title present")
	}

	if _, ok := data[0]["app"]; !ok {
		t.Error("No app present")
	}
}

func Test_limitThoudo(t *testing.T) {
	//("/api/apps?limit=10")
	form := url.Values{}
	form.Add("limit", "1000")

	urlStr := constructAPIURL("/api/apps", form)

	hc := http.Client{}

	req, err := http.NewRequest("POST", urlStr, nil)
	if err != nil {
		t.Errorf("Bad post", err)
	}
	req.Header.Add("Accept", "application/json")

	resp, err := hc.Do(req)

	if err != nil {
		t.Errorf("Could not access endpoint", err)
	}

	if resp.StatusCode != http.StatusOK {
		t.Errorf("Bad request", resp.Status)
	}

	respData, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		t.Errorf("Umable to read response", err)
	}

	data, err := getJSON(respData)

	if err != nil {
		t.Errorf("Unable to retrieve response data", err)
	}

	if len(data) != 1000 {
		t.Errorf("Should be more than 1000 matchs", data)
	}

	if _, ok := data[0]["title"]; !ok {
		t.Errorf("No title present")
	}

	if _, ok := data[0]["app"]; !ok {
		t.Error("No app present")
	}
}

func Test_limitHundoThoudo(t *testing.T) {
	//("/api/apps?limit=10")
	form := url.Values{}
	form.Add("limit", "100000")

	urlStr := constructAPIURL("/api/apps", form)

	hc := http.Client{}

	req, err := http.NewRequest("POST", urlStr, nil)
	if err != nil {
		t.Errorf("Bad post", err)
	}
	req.Header.Add("Accept", "application/json")

	resp, err := hc.Do(req)

	if err != nil {
		t.Errorf("Could not access endpoint", err)
	}

	if resp.StatusCode != http.StatusOK {
		t.Errorf("Bad request", resp.Status)
	}

	respData, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		t.Errorf("Umable to read response", err)
	}

	data, err := getJSON(respData)

	if err != nil {
		t.Errorf("Unable to retrieve response data", err)
	}

	if len(data) != 100000 {
		t.Errorf("Should be more than  100000 matchs", data)
	}

	if _, ok := data[0]["title"]; !ok {
		t.Errorf("No title present")
	}

	if _, ok := data[0]["app"]; !ok {
		t.Error("No app present")
	}
}
func Test_limitExcess(t *testing.T) {
	//("/api/apps?limit=10")
	form := url.Values{}
	form.Add("limit", "10000000")

	urlStr := constructAPIURL("/api/apps", form)

	hc := http.Client{}

	req, err := http.NewRequest("POST", urlStr, nil)
	if err != nil {
		t.Errorf("Bad post", err)
	}
	req.Header.Add("Accept", "application/json")

	resp, err := hc.Do(req)

	if err != nil {
		t.Errorf("Could not access endpoint", err)
	}

	if resp.StatusCode != http.StatusOK {
		t.Errorf("Bad request", resp.Status)
	}

	respData, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		t.Errorf("Umable to read response", err)
	}

	data, err := getJSON(respData)

	if err != nil {
		t.Errorf("Unable to retrieve response data", err)
	}

	if len(data) != 10000000 {
		t.Errorf("Should be more than 10 matchs", data)
	}

	if _, ok := data[0]["title"]; !ok {
		t.Errorf("No title present")
	}

	if _, ok := data[0]["app"]; !ok {
		t.Error("No app present")
	}
}
func Test_offset(t *testing.T) {
	//("/api/apps?offset=10")

}

func Test_geoIp(t *testing.T) {
	//("/api/")

}

func Test_iconGet(t *testing.T) {
	//("/api/")

}

func Test_host(t *testing.T) {
	//("/api/hosts")

}

/* BELOW ARE AUTO GENERATED TEST STRUCTURES */
func Test_toBytes(t *testing.T) {
	type args struct {
		data interface{}
	}
	tests := []struct {
		name string
		args args
		want []byte
	}{
	// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := toBytes(tt.args.data); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("toBytes() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_writeErr(t *testing.T) {
	type args struct {
		w      http.ResponseWriter
		mime   string
		status int
		err    string
		msg    string
		vals   []interface{}
	}
	tests := []struct {
		name string
		args args
	}{
	// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			writeErr(tt.args.w, tt.args.mime, tt.args.status, tt.args.err, tt.args.msg, tt.args.vals...)
		})
	}
}

func Test_writeData(t *testing.T) {
	type args struct {
		w      http.ResponseWriter
		mime   string
		status int
		data   interface{}
	}
	tests := []struct {
		name string
		args args
	}{
	// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			writeData(tt.args.w, tt.args.mime, tt.args.status, tt.args.data)
		})
	}
}

func Test_mimeCheck(t *testing.T) {
	type args struct {
		mime string
	}
	tests := []struct {
		name string
		args args
		want string
	}{
	// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := mimeCheck(tt.args.mime); got != tt.want {
				t.Errorf("mimeCheck() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_hello(t *testing.T) {
	type args struct {
		w http.ResponseWriter
		r *http.Request
	}
	tests := []struct {
		name string
		args args
	}{
	// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			hello(tt.args.w, tt.args.r)
		})
	}
}

func Test_parseNumCheck(t *testing.T) {
	type args struct {
		num string
	}
	tests := []struct {
		name     string
		args     args
		wantVal  int
		wantOops string
		wantErr  bool
	}{
	// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotVal, gotOops, err := parseNumCheck(tt.args.num)
			if (err != nil) != tt.wantErr {
				t.Errorf("parseNumCheck() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if gotVal != tt.wantVal {
				t.Errorf("parseNumCheck() gotVal = %v, want %v", gotVal, tt.wantVal)
			}
			if gotOops != tt.wantOops {
				t.Errorf("parseNumCheck() gotOops = %v, want %v", gotOops, tt.wantOops)
			}
		})
	}
}

func Test_parseLimit(t *testing.T) {
	type args struct {
		num string
	}
	tests := []struct {
		name     string
		args     args
		wantVal  string
		wantOops string
		wantErr  bool
	}{
	// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotVal, gotOops, err := parseLimit(tt.args.num)
			if (err != nil) != tt.wantErr {
				t.Errorf("parseLimit() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if gotVal != tt.wantVal {
				t.Errorf("parseLimit() gotVal = %v, want %v", gotVal, tt.wantVal)
			}
			if gotOops != tt.wantOops {
				t.Errorf("parseLimit() gotOops = %v, want %v", gotOops, tt.wantOops)
			}
		})
	}
}

func Test_parseOffset(t *testing.T) {
	type args struct {
		num string
	}
	tests := []struct {
		name     string
		args     args
		wantVal  string
		wantOops string
		wantErr  bool
	}{
	// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotVal, gotOops, err := parseOffset(tt.args.num)
			if (err != nil) != tt.wantErr {
				t.Errorf("parseOffset() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if gotVal != tt.wantVal {
				t.Errorf("parseOffset() gotVal = %v, want %v", gotVal, tt.wantVal)
			}
			if gotOops != tt.wantOops {
				t.Errorf("parseOffset() gotOops = %v, want %v", gotOops, tt.wantOops)
			}
		})
	}
}

func Test_fetchIDEndpoint(t *testing.T) {
	type args struct {
		w http.ResponseWriter
		r *http.Request
	}
	tests := []struct {
		name string
		args args
	}{
	// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			fetchIDEndpoint(tt.args.w, tt.args.r)
		})
	}
}

func Test_altAppsEndpoint(t *testing.T) {
	type args struct {
		w http.ResponseWriter
		r *http.Request
	}
	tests := []struct {
		name string
		args args
	}{
	// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			altAppsEndpoint(tt.args.w, tt.args.r)
		})
	}
}
