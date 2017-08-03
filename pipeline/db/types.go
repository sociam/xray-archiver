package db

import (
	"database/sql"
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
	ID   sql.NullString `json:"id"`
	Vers []int64        `json:"vers"`
}

// AppVersion represents all the information about an app version contained in
// the database
type AppVersion struct {
	ID          int64               `json:"id"`
	App         string              `json:"app"`
	Store       string              `json:"sql.NullString"`
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
	Title         sql.NullString   `json:"title"`
	Summary       sql.NullString   `json:"summary"`
	Description   sql.NullString   `json:"description"`
	StoreURL      sql.NullString   `json:"storeURL"`
	Price         sql.NullString   `json:"price"`
	Free          bool             `json:"free"`
	Rating        sql.NullString   `json:"rating"`
	NumReviews    int64            `json:"numReviews"`
	Genre         sql.NullString   `json:"genre"`
	FamilyGenre   sql.NullString   `json:"familyGenre"`
	Installs      Range            `json:"installs"`
	Developer     int64            `json:"developer"`
	Updated       time.Time        `json:"updated"`
	AndroidVer    sql.NullString   `json:"androidVer"`
	ContentRating sql.NullString   `json:"contentRating"`
	Screenshots   []sql.NullString `json:"screenshots"`
	Video         sql.NullString   `json:"video"`
	RecentChanges []sql.NullString `json:"recentChanges"`
	CrawlDate     time.Time        `json:"crawlDate"`
	Permissions   []sql.NullString `json:"permissions"`
}

// Developer represents a developer from the database
type Developer struct {
	ID        int64            `json:"id"`
	Emails    []sql.NullString `json:"emails"`
	Name      sql.NullString   `json:"name"`
	StoreSite sql.NullString   `json:"storeSite"`
	Site      sql.NullString   `json:"site"`
}

// Company represents a company from the database
type Company struct {
	ID           sql.NullString   `json:"id"`
	Name         sql.NullString   `json:"name"`
	Hosts        []sql.NullString `json:"nosts"`
	Founded      sql.NullString   `json:"founded"`
	Acquired     sql.NullString   `json:"acquired"`
	Type         []sql.NullString `json:"type"`
	TypeTag      sql.NullString   `json:"typeTag"`
	Jurisdiction sql.NullString   `json:"jurisdiction"`
	Parent       sql.NullString   `json:"parent"`
	Capital      sql.NullString   `json:"capital"`
	Equity       sql.NullString   `json:"equity"`
	Size         Range            `json:"size"`
	DataSources  []sql.NullString `json:"dataSources"`
	Description  []sql.NullString `json:"description"`
}

// UtilApp creates a *util.App from an AppVersion
func (a AppVersion) UtilApp() *util.App {
	return util.NewApp(a.ID, a.App, a.Store, a.Region, a.Ver)
}
