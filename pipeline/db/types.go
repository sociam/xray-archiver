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
	ID          int64     `json:"id"`
	App         string    `json:"app"`
	Store       string    `json:"string"`
	Region      string    `json:"region"`
	Ver         string    `json:"ver"`
	ScreenFlags int64     `json:"screenFlags"`
	StoreInfo   StoreInfo `json:"storeinfo"`
	Icon        string    `json:"icon"`
	Dev         Developer `json:"developer"`
	Hosts       []string  `json:"hosts"`
	Perms       []string  `json:"perms"`
	Packages    []string  `json:"packages"`
	IsAnalyzed  bool      `json:"isAnalyzed"`
}

// GenreAvg represents a row from the Genre average statistics table
type GenreStats struct {
	Category  string  `json:"category"`
	HostCount int64   `json:"hostCount"`
	AppCount  int64   `json:"appCount"`
	GenreAvg  float64 `json:"genreAvg"`
}

// AltApp represents Alternative App title and playstore url. Used in API
type AltApp struct {
	AltAppTitle     string `json:"altAppTitle"`
	AltToURL        string `json:"altToURL"`
	GPlayURL        string `json:"gPlayURL"`
	GPlayID         string `json:"gPlayID"`
	IconURL         string `json:"iconURL"`
	OfficialSiteURL string `json:"officialSiteURL"`
	IsScraped       bool   `json:"isScraped"`
}

// StoreInfo represents the information contained about an app in its respective
// store
type StoreInfo interface{}

// PlayStoreInfo represents the data contained in the google play store
type PlayStoreInfo struct {
	Title         string    `json:"title"`
	Summary       string    `json:"summary"` //TODO:Omitempty for null vlaues?
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
	Emails    []string `json:"emails"`
	Name      string   `json:"name"`
	StoreSite string   `json:"storeSite"`
	Site      string   `json:"site"`
}

// AppStub represents the 'stub' of app data located in the Database.
// it is formed of only the title given to the app by the developer
// and the App ID that it can be uniquely identified as.
type AppStub struct {
	Title string `json:"title"`
	App   string `json:"app"`
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
