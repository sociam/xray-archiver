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

var useDB bool
var db xrayDb

// Open opens the database with the given config. If enable is false, the
// functions that modify the database are noops.
func Open(cfg util.Config, enable bool) error {
	if enable {
		useDB = true
		sqlDb, err := sql.Open("postgres",
			fmt.Sprintf("dbname='%s' user='%s' password='%s' host='%s' port='%d' sslmode='disable'",
				cfg.DB.Database, cfg.DB.User, cfg.DB.Password, cfg.DB.Host, cfg.DB.Port))
		if err != nil {
			return err
		}
		db = xrayDb{sqlDb}
	}
	return nil
}

// func (db *XrayDb) insertApp(app *App) (int, error) {
// 	if !useDB {
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

//TODO: make Add* functions take a db id and what to add instead of a util.App

// AddPackages is a function that allows you to add packages to the Xray DB. The
// argument app must contain a DB ID and an array of package names.
func AddPackages(app *util.App) error {
	if !useDB || app.DBID == 0 {
		return nil
	}

	var dbPkgs []string
	// Check if app and packages already exist.
	err := db.QueryRow("SELECT packages FROM app_packages WHERE id = $1", app.DBID).
		Scan(pq.Array(&dbPkgs))
	if err != nil {
		if err != sql.ErrNoRows {
			return err
		}
		// If the app doesn't exist, insert with whole list of packages
		rows, err := db.Query("INSERT INTO app_packages VALUES ($1, $2)",
			app.DBID, pq.Array(&app.Packages))
		rows.Close()
		if err != nil {
			return err
		}
	} else {
		// The app already exists, only add new packages for that app.
		bothPkgs := util.UniqAppend(app.Packages, dbPkgs)
		rows, err := db.Query("UPDATE app_packages SET perms = $1 WHERE id = $2",
			pq.Array(&bothPkgs), app.DBID)
		rows.Close()
		if err != nil {
			return err
		}
	}

	return nil
}

// AddPerms is a function that allows you to add permissions to the Xray DB. The
// argument app must contain a DB ID and an array of permissions.
func AddPerms(app *util.App) error {
	if !useDB || app.DBID == 0 {
		return nil
	}

	sPerms := make([]string, 0, len(app.Perms))
	for _, perm := range app.Perms {
		sPerms = append(sPerms, perm.ID)
	}

	var dbPerms []string
	err := db.QueryRow("SELECT perms FROM app_perms WHERE id = $1", app.DBID).
		Scan(pq.Array(&dbPerms))
	if err != nil {
		if err != sql.ErrNoRows {
			return err
		}
		rows, err := db.Query("INSERT INTO app_perms VALUES ($1, $2)",
			app.DBID, pq.Array(&sPerms))
		rows.Close()
		if err != nil {
			return err
		}
	} else {
		bothPerms := util.UniqAppend(sPerms, dbPerms)
		rows, err := db.Query("UPDATE app_perms SET perms = $1 WHERE id = $2",
			pq.Array(&bothPerms), app.DBID)
		rows.Close()
		if err != nil {
			return err
		}
	}

	return nil
}

// AddHosts is a function that allows you to add hosts to the Xray DB. The
// argument app must contain a DB ID.
func AddHosts(app *util.App, hosts []string) error {
	if !useDB || app.DBID == 0 {
		return nil
	}

	var dbHosts []string
	err := db.QueryRow("SELECT hosts FROM app_hosts WHERE id = $1", app.DBID).
		Scan(pq.Array(&dbHosts))
	if err != nil {
		if err != sql.ErrNoRows {
			return err
		}
		rows, err := db.Query("INSERT INTO app_hosts(id, hosts) VALUES ($1, $2)",
			app.DBID, pq.Array(&hosts))
		rows.Close()
		if err != nil {
			return err
		}
	} else {
		bothHosts := util.UniqAppend(hosts, dbHosts)
		rows, err := db.Query("UPDATE app_host SET hosts = $1 WHERE id = $2",
			pq.Array(&bothHosts), app.DBID)
		rows.Close()
		if err != nil {
			return err
		}
	}

	return nil
}

// GetAppVersion gets an app version from the database. The argument app is the
// app id, in the form com.example.app.
func GetAppVersion(app, store, region, version string) (AppVersion, error) {
	var appVer AppVersion

	err := db.QueryRow(
		"SELECT * FROM app_versions WHERE app = $1 AND store = $2 AND region = $3 AND version = $4",
		app,
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

// GetAppVersionByID gets an app version given its ID in the database.
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

// GetDeveloper gets a developer given its ID in the database.
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

// GetDevelopers returns a list of developers.
func GetDevelopers(num, start int) ([]Developer, error) {
	rows, err := db.Query("SELECT * FROM developers LIMIT $1 OFFSET $2", num, start)
	defer rows.Close()
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

// GetCompany returns a company given its ID in the database.
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

// GetCompanies returns a list of companies.
func GetCompanies(num, start int) ([]Company, error) {
	rows, err := db.Query("SELECT * FROM companies LIMIT $1 OFFSET $2", num, start)
	defer rows.Close()
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

// GetApp gets an app given its id, in the form com.example.app.
func GetApp(id string) (App, error) {
	var app App
	err := db.QueryRow("SELECT * FROM apps WHERE id = $1", id).Scan(&app.ID, pq.Array(&app.Vers))
	if err != nil {
		return App{}, err
	}

	return app, nil
}

// GetApps returns a list of apps.
func GetApps(num, start int) ([]App, error) {
	rows, err := db.Query("SELECT * FROM apps LIMIT $1 OFFSET $2", num, start)
	defer rows.Close()
	if err != nil {
		return []App{}, err
	}
	ret := make([]App, 0, num)
	for i := 0; rows.Next(); i++ {
		var app App
		err := rows.Scan(
			&app.ID,
			pq.Array(&app.Vers))
		if err != nil {
			fmt.Println("Database err:", err)
		} else {
			ret = append(ret, app)
		}
	}

	if rows.Err() != sql.ErrNoRows && rows.Err() != nil {
		fmt.Println("Databse err", rows.Err())
		return []App{}, rows.Err()
	}

	return ret, nil
}

// GetLatestApps returns a list of the newest app versions.
//
// TODO: THIS IS WRONG.
func GetLatestApps(num, start int) ([]App, error) {
	//TOOD: db join for only only latest to get through
	rows, err := db.Query("SELECT * FROM apps LIMIT $1 OFFSET $2", num, start)
	defer rows.Close()
	if err != nil {
		return []App{}, err
	}

	ret := make([]App, 0, num)

	for i := 0; rows.Next(); i++ {
		ret = append(ret, App{})
		rows.Scan(&ret[i].ID, pq.Array(&ret[i].Vers))
	}

	if rows.Err() != sql.ErrNoRows && rows.Err() != nil {
		fmt.Println("Databse err", rows.Err())
		return []App{}, rows.Err()
	}

	return ret, nil
}

// SearchApps returns a list of apps matching a search term.
func SearchApps(searchTerm string) ([]PlaystoreInfo, error) {
	searchTerm = "%" + searchTerm + "%"

	rows, err := db.Query("SELECT * from playstore_apps WHERE title like $1 ORDER BY rating USING> LIMIT 120", searchTerm)
	defer rows.Close()
	if err != nil {
		return []PlaystoreInfo{}, err
	}

	var ret []PlaystoreInfo

	var id string

	for i := 0; rows.Next(); i++ {
		var inf PlaystoreInfo
		err := rows.Scan(
			&id,
			&inf.Title,
			&inf.Summary,
			&inf.Description,
			&inf.StoreURL,
			&inf.Price,
			&inf.Free,
			&inf.Rating,
			&inf.NumReviews,
			&inf.Genre,
			&inf.FamilyGenre,
			&inf.Installs.Min,
			&inf.Installs.Max,
			&inf.Developer,
			&inf.Updated,
			&inf.AndroidVer,
			&inf.ContentRating,
			pq.Array(&inf.Screenshots),
			&inf.Video,
			pq.Array(&inf.RecentChanges),
			&inf.CrawlDate,
			pq.Array(&inf.Permissions))
		if err != nil {
			fmt.Println(err)
		} else {
			ret = append(ret, inf)
		}
	}

	if rows.Err() != sql.ErrNoRows && rows.Err() != nil {
		fmt.Println("Database err", rows.Err())
		return []PlaystoreInfo{}, rows.Err()
	}

	return ret, nil
}

// GetAppsToAnalyze returns a list of up to 10 apps that have analyzed=False and
// downloaded=True for the analyzer.
func GetAppsToAnalyze() ([]AppVersion, error) {
	rows, err := db.Query("SELECT id, app, store, region, version, screen_flags, icon FROM app_versions WHERE analyzed = False AND downloaded = True LIMIT 10")
	defer rows.Close()
	if err != nil {
		return []AppVersion{}, err
	}
	ret := make([]AppVersion, 0, 10)
	for i := 0; rows.Next(); i++ {
		var cur AppVersion
		var screenFlags sql.NullInt64
		var icon sql.NullString
		err := rows.Scan(&cur.ID, &cur.App, &cur.Store, &cur.Region, &cur.Ver, &screenFlags, &icon)
		cur.ScreenFlags = screenFlags.Int64
		cur.Icon = icon.String

		if err != nil {
			fmt.Println(err)
		} else {
			ret = append(ret, cur)
		}
	}

	if rows.Err() != sql.ErrNoRows && rows.Err() != nil {
		return []AppVersion{}, err
	}

	return ret, nil
}

// UnsetDownloaded sets an downloaded=False for given app.
func UnsetDownloaded(id int64) error {
	rows, err := db.Query("UPDATE app_versions SET downloaded = False WHERE id = $1", id)
	rows.Close()
	return err
}

// SetAnalyzed sets analyzed=True for a given app.
func SetAnalyzed(id int64) error {
	rows, err := db.Query("UPDATE app_versions SET analyzed = True WHERE id = $1", id)
	rows.Close()
	return err
}
