package db

import (
	"time"

	"github.com/sociam/xray-archiver/pipeline/util"
)

type Range struct {
	Min int64 `json:"min"`
	Max int64 `json:"max"`
}

type App struct {
	ID   string  `json:"id"`
	Vers []int64 `json:"vers"`
	Icon string  `json:"icon"`
}

type AppVersion struct {
	ID          int64               `json:"id"`
	App         string              `json:"app"`
	Store       string              `json:"string"`
	Region      string              `json:"region"`
	Ver         string              `json:"ver"`
	ScreenFlags int64               `json:"screenFlags"`
	Hosts       map[string][]string `json:"hosts"`
	Perms       []util.Permission   `json:"permissions"`
	StoreInfo   storeInfo           `json:"storeinfo"`
}

type storeInfo interface{}

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

type Developer struct {
	ID        int64    `json:"id"`
	Emails    []string `json:"emails"`
	Name      string   `json:"name"`
	StoreSite string   `json:"storeSite"`
	Site      string   `json:"site"`
}

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
