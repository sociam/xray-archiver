package main

import (
	"database/sql"
	"flag"
	"net/http"
	"reflect"
	"testing"
	"time"

	"github.com/sociam/xray-archiver/pipeline/db"
	"github.com/sociam/xray-archiver/pipeline/util"
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

var cfgFileTest = flag.String("cfg", "/etc/xray/config.json", "config file location")

var portTest = flag.Uint("port", 8118, "Port to serve on.")

func TestInit(t *testing.T) {
	var err error
	err = util.LoadCfg(*cfgFileTest, util.APIServ) //TODO:might break all setup
	if err != nil {
		t.Errorf("Could not load and setup database config", err)
	}

	err = db.Open(util.Cfg, true)
	if err != nil {
		t.Errorf("Could not connect to database", err)
	}

	//ID should be 1..
	rows, err := database.Query("INSERT INTO apps VALUES ($1, $2)",
		ExampleApp.ID, ExampleApp.Vers[0])

	if rows != nil {
		rows.Close()
	}

	if err != nil {
		t.Errorf("Could not add exampel host", err)
	}

	xrayDb.Query("INSERT INTO app_versions(app, store, region, version, downloaded, last_dl_attempt, analyzed) VALUES ($1, $2, $3, $4, $5, $6, $7) RETURNING id",
		ExampleApp.ID,
		ExampleAppVer.Store,
		ExampleAppVer.Region,
		ExampleAppVer.Ver,
		0,
		"epoch",
		0)

	if rows != nil {
		rows.Close()
	}

	if err != nil {
		t.Errorf("Could not add appversions example data", err)
	}

	rows, err := xrayDb.Query("INSERT INTO playstore_apps VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, $13, $14, $15, $16, $17, $18, $19, $20, $21)",
		ExamplePlaystore.AndroidVer,
		ExamplePlaystore.Title,
		ExamplePlaystore.Summary,
		ExamplePlaystore.Description,
		ExamplePlaystore.StoreURL,
		ExamplePlaystore.Price,
		ExamplePlaystore.Free,
		ExamplePlaystore.Rating,
		ExamplePlaystore.NumReviews,
		ExamplePlaystore.Genre,
		ExamplePlaystore.FamilyGenre,
		ExamplePlaystore.Installs.Min,
		ExamplePlaystore.Installs.Max,
		ExamplePlaystore.Developer,
		ExamplePlaystore.Updated,
		ExamplePlaystore.AndroidVer,
		ExamplePlaystore.ContentRating,
		ExamplePlaystore.Screenshots,
		ExamplePlaystore.Video,
		ExamplePlaystore.RecentChanges,
		ExamplePlaystore.Updated)

	if rows != nil {
		rows.Close()
	}

	if err != nil {
		t.Errorf("Could not add playstore_app  example data", err)
	}

	rows, err := xrayDb.Query("SELECT * FROM apps WHERE id = $1", ExampleApp.ID)
	if rows != nil {
		rows.Close()
	}

	if err != nil {
		t.Errorf("Could not select  example data", err)
	}

	if rows != 1 {
		t.Errorf("Correct number of example row data was not added", rows)
	}

}

/*API ENDPOINT TEST*/
//Test the app endppoint or the query? the problem is main barely have shiz...
//just querying the real running endpoint?

func Test_DefaultLimit() {
	appsEndpoint()
}

func Test_IsFulParam(t *testing.T) {
	appsEndpoint("/api/apps?isFull=True")
	t.Fail()
}

func Test_Title(t *testing.T) {
	appsEndpoint("/api/apps?title=DinnerAtNoon")
	t.Fail()
}

func Test_MultiTitle(t *testing.T) {
	appsEndpoint("/api/apps?title=DinnerAtNoon,cats")
	appsEndpoint("/api/apps?title=dinner&title=at&title=noon")
	t.Fail()
}

func Test_Developer(t *testing.T) {
	appsEndpoint("/api/apps?developer=zynga")
	t.Fail()
}

func Test_MultiDevelopers(t *testing.T) {
	appsEndpoint("/api/apps?developer=zynga,cat")
	appsEndpoint("/api/apps?developer=zynga&developer=cat")
	t.Fail()
}

func Test_Genre(t *testing.T) {
	appsEndpoint("/api/apps?genre=entertainment")
	t.Fail()
}

func Test_MultiGenre(t *testing.T) {
	appsEndpoint("/api/apps?genre=lifestyle,food_and_drink")
	appsEndpoint("/api/apps?genre=lifestyle&genre=food_and_drink")
	t.Fail()
}

func Test_GameGenre(t *testing.T) {
	appsEndpoint("/api/apps?genre=GAMES")
	t.Fail()
}

func Test_AppId(t *testing.T) {
	appsEndpoint("/api/apps?appId")
	t.Fail()
}

func Test_MultiAppId(t *testing.T) {
	appsEndpoint("/api/apps?appId=example1&appId=example2")
	appsEndpoint("/api/apps?appId=example1,example2")
	t.Fail()
}

func Test_invalidAppid(t *testing.T) {
	appsEndpoint("/api/apps?appId=INVALID")
	t.Fail()
}

func Test_startsWith(t *testing.T) {
	appsEndpoint("/api/apps?startsWith=i")
	t.Fail()
}

func Test_onlyAnalyzed(t *testing.T) {
	appsEndpoint("/api/apps?onlyAnalyzed")
	t.Fail()
}

func Test_limit(t *testing.T) {
	appsEndpoint("/api/apps?limit=10")

	t.Fail()
}

func Test_offset(t *testing.T) {
	appsEndpoint("/api/apps?offset=10")

	t.Fail()
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

func Test_appsEndpoint(t *testing.T) {
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
			appsEndpoint(tt.args.w, tt.args.r)
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
