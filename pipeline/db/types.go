package db

import (
	"time"

	"github.com/sociam/xray-archiver/pipeline/util"
)

// Range represents a range
type Range struct {
	Min int64 `json:"min"`
	Max int64 `json:"max"`
}

// App represents an app from the database
type App struct {
	ID   string  `json:"id"`
	Vers []int64 `json:"vers"`
}

// AppVersion represents all the information about an app version contained in
// the database
type AppVersion struct {
	ID          int64               `json:"id"`
	App         string              `json:"app"`
	Store       string              `json:"string"`
	Region      string              `json:"region"`
	Ver         string              `json:"ver"`
	ScreenFlags int64               `json:"screenFlags"`
	Hosts       map[string][]string `json:"hosts"`
	Perms       []util.Permission   `json:"permissions"`
	StoreInfo   StoreInfo           `json:"storeinfo"`
	Icon        string              `json:"icon"`
}

// StoreInfo represents the information contained about an app in its respective
// store
type StoreInfo interface{}

// PlaystoreInfo represents the data contained in the google play store
type PlaystoreInfo struct {
	Title         string    `json:"title"`
	Summary       string    `json:"summary"`
	Description   string    `json:"description"`
	StoreURL      string    `json:"storeURL"`
	Price         string    `json:"price"`
	Free          bool      `json:"free"`
	Rating        string    `json:"rating"`
	NumReviews    int64     `json:"numReviews"`
	Genre         string    `json:"genre"`
	FamilyGenre   string    `json:"familyGenre"`
	Installs      Range     `json:"installs"`
	Developer     int64     `json:"developer"`
	Updated       time.Time `json:"updated"`
	AndroidVer    string    `json:"androidVer"`
	ContentRating string    `json:"contentRating"`
	Screenshots   []string  `json:"screenshots"`
	Video         string    `json:"video"`
	RecentChanges []string  `json:"recentChanges"`
	CrawlDate     time.Time `json:"crawlDate"`
	Permissions   []string  `json:"permissions"`
}

// Developer represents a developer from the database
type Developer struct {
	ID        int64    `json:"id"`
	Emails    []string `json:"emails"`
	Name      string   `json:"name"`
	StoreSite string   `json:"storeSite"`
	Site      string   `json:"site"`
}

// AppData struct combines playstore info, appVersion abd developer info.
type AppData struct {
	PlaystoreInfo
	AppVersion
	Developer
}

// AppStub represents the 'stub' of app data located in the Database.
// it is formed of only the title given to the app by the developer
// and the App ID that it can be uniquely identified as.
type AppStub struct {
	Title string
	App   string
}

// Company represents a company from the database
type Company struct {
	ID           string   `json:"id"`
	Name         string   `json:"name"`
	Hosts        []string `json:"nosts"`
	Founded      string   `json:"founded"`
	Acquired     string   `json:"acquired"`
	Type         []string `json:"type"`
	TypeTag      string   `json:"typeTag"`
	Jurisdiction string   `json:"jurisdiction"`
	Parent       string   `json:"parent"`
	Capital      string   `json:"capital"`
	Equity       string   `json:"equity"`
	Size         Range    `json:"size"`
	DataSources  []string `json:"dataSources"`
	Description  []string `json:"description"`
}

// UtilApp creates a *util.App from an AppVersion
func (a AppVersion) UtilApp() *util.App {
	return util.NewApp(a.ID, a.App, a.Store, a.Region, a.Ver)
}

/* FormParam is is a struct to represent parameters passed through a URL
// please note that form params are infact case senstive*/
type FormParam struct {
	Name string `json:"name"`
	Val  string `json:"val"`
}

// NewFormParam constructs a new struct consisting of a
// form elements name and the value associated with it.
func NewFormParam(name string, val string) *FormParam {
	return &FormParam{Name: name, Val: val}
}

// UtilFormApps
func (f *FormParam) UtilFormApps() *FormParam {
	// doing something with id here
	return NewFormParam(f.Name, f.Val)
}
