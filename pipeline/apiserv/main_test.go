package main

import (
	"net/http"
	"reflect"
	"testing"
	"github.com/sociam/xray-archiver/pipeline/util"
	"github.com/sociam/xray-archiver/pipeline/db"
)





/* Db config */
type xrayDb struct {
	*sql.DB
}
var useDB bool
var db xrayDb

var EXAMPLE_APP = App{ID: "example", Ver: {1.0}}	
var EXAMPLE_STOREINFO = PlayStoreInfo{Title: "EXAMPLE TITLE",
									Summary: "I AM A EXAMPLE",
									Description: "I WISH TO EXAMPLE",
									StoreURL: "/path/to/store_url",
									Price: "Free",
									Free: true,
									Rating: "3.7",
									NumReviews: 1,
									Genre: "GAME_CASUAL",
									FamilyGenre: "",
									Installs: 20,
									Developer: 1,
									Updated: "2017-02-21",
									AndroidVer: 4.1,
									Screenshots: {},
									Video: "",
									RecentChanges: {},
									CrawlDate: "2017-08-16",
									Permissions: "" //where are the permission
									}

var EXAMPLE_DEV = Developer{Emails: {"example@email.com"},
							Name: "example_name",
							StoreSite: "/path/to/store_site",
							Site: "path/to/site"}

var EXAMPLE_APP_VER = AppVersion{ID: 1, 
							App:EXAMPLE_APP.ID, 
							Store: "play", 
							Region: "uk", 
							Ver:EXAMPLE_APP.Vers[0]
							ScreenFlags: 0,
							StoreInfo: EXAMPLE_STOREINFO 
							Icon: "",
							Dev: 1,	
							Hosts: {"example.com"},	
							Perms:       {},	
							Packages:    {},	
							IsAnalyzed: true}

func TestInit(t *testing.T) {
	err := util.LoadCfg(*cfgFile, util.APIServ)
	if err != nil {
		t.Errorf("Could not load and setup database config", err)
	}
	
	err := db.Open(util.Cfg, true)
	if err != nil {
		t.Errorf("Could not connect to database", err)
	}
	

	//ID should be 1..
	rows, err := db.Query("INSERT INTO apps VALUES ($1, $2)",
			EXAMPLE_APP.ID, EXAMPLE_APP.Ver)
	
	if rows != nil {
		rows.Close()
	} 

	if err != nil {
		t.Errorf("Could not add exampel host", err)
	}


	rows, err := db.Query("INSERT INTO app_versions(app, store, region, version, downloaded, last_dl_attempt, analyzed )" +
					"VALUES ($1, $2, $3, $4, $5, $6, $7) RETURNING id", [app.appId, 
					EXAMPLE_APP_VER.Store, EXAMPLE_APP_VER.Region, 
					EXAMPLE_APP_VER.Ver, 0, "epoch", 0])
					
	if rows != nil {
		rows.Close()
	} 

	if err != nil {
		t.Errorf("Could not add appversions example data", err)
	}

	rows, err := db.Query("INSERT INTO playstore_apps VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, $13, $14, $15, $16, $17, $18, $19, $20, current_date)", [
                    verId,
                    EXAMPLE_STOREINFO.title,
                    EXAMPLE_STOREINFO.summary,
                    EXAMPLE_STOREINFO.description,
                    EXAMPLE_STOREINFO.url,
                    EXAMPLE_STOREINFO.price,
                    EXAMPLE_STOREINFO.free,
                    EXAMPLE_STOREINFO.score,
                    EXAMPLE_STOREINFO.reviews,
                    EXAMPLE_STOREINFO.genreId,
                    EXAMPLE_STOREINFO.familyGenreId,
                    EXAMPLE_STOREINFO.minInstalls,
                    EXAMPLE_STOREINFO.maxInstalls,
                    EXAMPLE_DEV.ID,
                    EXAMPLE_STOREINFO.updated,
                    EXAMPLE_STOREINFO.androidVersion,
                    EXAMPLE_STOREINFO.contentRating,
                    EXAMPLE_STOREINFO.screenshots,
                    EXAMPLE_STOREINFO.video,
                    EXAMPLE_STOREINFO.recentChanges
                ]);
    )
	
	if rows != nil {
		rows.Close()
	} 

	if err != nil {
		t.Errorf("Could not add playstore_app  example data", err)
	}

	



}


func insertTestData() {
	
	'INSERT INTO app_versions(app, store, region, version, downloaded, last_dl_attempt, analyzed )' +
                    'VALUES ($1, $2, $3, $4, $5, $6, $7) RETURNING id', [app.appId, 'play', region, app.version, 0, 'epoch', 0]

	               'INSERT INTO playstore_apps VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, $13, $14, $15, $16, $17, $18, $19, $20, current_date)', [
                    verId,
                    app.title,
                    app.summary,
                    app.description,
                    app.url,
                    app.price,
                    app.free,
                    app.score,
                    app.reviews,
                    app.genreId,
                    app.familyGenreId,
                    app.minInstalls,
                    app.maxInstalls,
                    devId,
                    app.updated,
                    app.androidVersion,
                    app.contentRating,
                    app.screenshots,
                    app.video,
                    app.recentChanges
                ]);
            await client.lquery('COMMIT');

	'SELECT * FROM apps WHERE id = $1', [app.appId]

}

func Test_IsFulParam(t *testing.T) {

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
