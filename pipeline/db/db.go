package db

import (
	"database/sql"

	"github.com/lib/pq"
	"github.com/sociam/xray-archiver/pipeline/util"

	"fmt"
)

type xrayDb struct {
	*sql.DB
}

var useDb bool
var db xrayDb

func Open(cfg util.Config) error {
	sqlDb, err := sql.Open("postgres",
		fmt.Sprintf("dbname='%s' user='%s' password='%s' host='%s' port='%d' sslmode='disable'",
			cfg.Db.Database, cfg.Db.User, cfg.Db.Password, cfg.Db.Host, cfg.Db.Port))
	if err != nil {
		return err
	}
	db = xrayDb{sqlDb}
	return nil
}

// func (db *XrayDb) insertApp(app *App) (int, error) {
// 	if !useDb {
// 		return 0, nil
// 	}

// 	var ret int
// 	err := db.QueryRow(
// 		"SELECT id FROM app_versions WHERE app = $1 AND store = $2 AND region = $3 AND version = $4",
// 		app.id, app.store, app.region, app.ver).Scan(&ret)
// 	if err != nil {
// 		return 0, err
// 	}

// 	return ret, nil
// }

func AddPerms(app *util.App, perms []util.Permission) error {
	if !useDb || app.DbId == 0 {
		return nil
	}

	sPerms := make([]string, 0, len(perms))
	for _, perm := range perms {
		sPerms = append(sPerms, perm.Id)
	}

	var dbPerms []string
	err := db.QueryRow("SELECT perms FROM app_perms WHERE id = $1", app.DbId).
		Scan(pq.Array(&dbPerms))
	if err != nil {
		if err != sql.ErrNoRows {
			return err
		}
		_, err := db.Query("INSERT INTO app_perms VALUES ($1, $2)",
			app.DbId, pq.Array(&sPerms))
		if err != nil {
			return err
		}
	} else {
		bothPerms := util.UniqAppend(sPerms, dbPerms)
		_, err := db.Query("UPDATE app_perms SET perms = $1 WHERE id = $2",
			pq.Array(&bothPerms), app.DbId)
		if err != nil {
			return err
		}
	}

	return nil
}

func AddHosts(app *util.App, hosts []string) error {
	if !useDb || app.DbId == 0 {
		return nil
	}

	var dbHosts []string
	err := db.QueryRow("SELECT hosts FROM app_hosts WHERE id = $1", app.DbId).
		Scan(pq.Array(&dbHosts))
	if err != nil {
		if err != sql.ErrNoRows {
			return err
		}
		_, err := db.Query("INSERT INTO app_hosts VALUES ($1, $2)",
			app.DbId, pq.Array(&hosts))
		if err != nil {
			return err
		}
	} else {
		bothHosts := util.UniqAppend(hosts, dbHosts)
		_, err := db.Query("UPDATE app_host SET hosts = $1 WHERE id = $2",
			pq.Array(&bothHosts), app.DbId)
		if err != nil {
			return err
		}
	}

	return nil
}

func GetAppVersion(store, region, version string) (AppVersion, error) {
	var appVer AppVersion

	err := db.QueryRow(
		"SELECT * FROM app_versions WHERE store = $1 AND region = $2 AND version = $3",
		store,
		region,
		version,
	).Scan(
		&appVer.ID,
		&appVer.App,
		&appVer.Store,
		&appVer.Region,
		&appVer.Ver,
		&appVer.ScreenFlags,
	)

	if err != nil {
		return AppVersion{}, err
	}

	// TODO: Check implementation
	err = db.QueryRow("SELECT * FROM app_hosts WHERE id = $1", appVer.ID).Scan(
		pq.Array(&appVer.Hosts))

	if err != nil {
		return AppVersion{}, err
	}

	// TODO: Check implementation
	err = db.QueryRow("SELECT * FROM app_perms WHERE id = $1", appVer.ID).Scan(
		pq.Array(&appVer.Perms))

	if err != nil {
		return AppVersion{}, err
	}

	return appVer, nil
}

func GetAppVersionByID(id int64) (AppVersion, error) {
	var appVer AppVersion

	err := db.QueryRow(
		"SELECT * FROM app_versions WHERE id = $1", id).Scan(
		&appVer.ID,
		&appVer.App,
		&appVer.Store,
		&appVer.Region,
		&appVer.Ver,
		&appVer.ScreenFlags)

	if err != nil {
		return appVer, err
	}

	// TODO: Check implementation

	err = db.QueryRow("SELECT * FROM app_hosts WHERE id = $1", appVer.ID).Scan(
		pq.Array(&appVer.Hosts))
	if err != nil {
		return AppVersion{}, err
	}

	// TODO: Check implementation
	err = db.QueryRow("SELECT * FROM app_perms WHERE id = $1", appVer.ID).Scan(
		pq.Array(&appVer.Perms))

	if err != nil {
		return AppVersion{}, err
	}

	return appVer, nil
}

func GetDeveloper(id int64) (Developer, error) {
	var dev Developer

	err := db.QueryRow("SELECT * from developers WHERE id = $1", id).Scan(
		&dev.ID,
		pq.Array(&dev.Emails),
		&dev.Name,
		&dev.StoreSite,
		&dev.Site)

	if err != nil {
		return Developer{}, err
	}

	return dev, nil

}

func GetDevelopers(num, start int) ([]Developer, error) {
	rows, err := db.Query("SELECT * FROM developers LIMIT $1 OFFSET $2", num, start)
	if err != nil {
		return []Developer{}, err
	}
	ret := make([]Developer, 0, num)
	for i := 0; rows.Next(); i++ {
		ret = append(ret, Developer{})
		rows.Scan(
			&ret[i].ID,
			pq.Array(&ret[i].Emails),
			&ret[i].Name,
			&ret[i].Site,
			&ret[i].StoreSite)
	}

	if rows.Err() != sql.ErrNoRows {
		return []Developer{}, err
	}

	return ret, nil
}

func GetCompany(id string) (Company, error) {
	var comp Company

	err := db.QueryRow("SELECT * from companies WHERE id = $1", id).Scan(
		&comp.ID,
		&comp.Name,
		pq.Array(&comp.Hosts),
		&comp.Founded,
		&comp.Acquired,
		pq.Array(&comp.Type),
		&comp.TypeTag,
		&comp.Jurisdiction,
		&comp.Parent,
		&comp.Capital,
		&comp.Equity,
		&comp.Size,
		pq.Array(&comp.DataSources),
		&comp.Description)

	if err != nil {
		return Company{}, err
	}

	return comp, nil
}

func GetCompanies(num, start int) ([]Company, error) {
	rows, err := db.Query("SELECT * FROM companies LIMIT $1 OFFSET $2", num, start)
	if err != nil {
		return []Company{}, err
	}
	ret := make([]Company, 0, num)

	for i := 0; rows.Next(); i++ {
		ret = append(ret, Company{})
		rows.Scan(
			&ret[i].ID,
			&ret[i].Name,
			pq.Array(&ret[i].Hosts),
			&ret[i].Founded,
			&ret[i].Acquired,
			pq.Array(&ret[i].Type),
			&ret[i].TypeTag,
			&ret[i].Jurisdiction,
			&ret[i].Parent,
			&ret[i].Capital,
			&ret[i].Equity,
			&ret[i].Size,
			pq.Array(&ret[i].DataSources),
			&ret[i].Description)
	}

	if rows.Err() != sql.ErrNoRows {
		return []Company{}, err
	}

	return ret, nil
}

func GetApp(id string) (App, error) {
	var app App
	err := db.QueryRow("SELECT * FROM apps WHERE id = $1", id).Scan(&app.ID, pq.Array(&app.Vers), &app.Icon)

	if err != nil {
		return App{}, err
	}

	return app, nil
}

func GetApps(num, start int) ([]App, error) {
	rows, err := db.Query("SELECT * FROM apps LIMIT $1 OFFSET $2", num, start)
	if err != nil {
		return []App{}, err
	}
	ret := make([]App, 0, num)
	for i := 0; rows.Next(); i++ {
		ret = append(ret, App{})
		rows.Scan(&ret[i].ID, pq.Array(&ret[i].Vers), &ret[i].Icon)
	}

	if rows.Err() != sql.ErrNoRows {
		return []App{}, err
	}

	return ret, nil
}

func GetLatestApps(num, start int) ([]App, error) {
	//TOOD: db join for only only latest to get through
	rows, err := db.Query("SELECT * FROM apps LIMIT $1 OFFSET $2", num, start)
	if err != nil {
		return []App{}, err
	}
	ret := make([]App, 0, num)
	for i := 0; rows.Next(); i++ {
		ret = append(ret, App{})
		rows.Scan(&ret[i].ID, pq.Array(&ret[i].Vers), &ret[i].Icon)
	}

	if rows.Err() != sql.ErrNoRows {
		return []App{}, err
	}

	return ret, nil
}
