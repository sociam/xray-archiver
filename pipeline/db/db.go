package db

import (
	"database/sql"
	"fmt"
	"time"

	"github.com/lib/pq"
	"github.com/sociam/xray-archiver/pipeline/util"
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

//TODO: make Add* functions take a db id and what to add instead of a util.App

// SetLastAnalyzeAttempt sets the last_analyzed_attempt of an app to the
// current time.
func SetLastAnalyzeAttempt(id int64) error {
	rows, err := db.Query("UPDATE app_versions SET last_analyze_attempt = $1 WHERE id = $2", time.Now(), id)
	if rows != nil {
		rows.Close()
	}
	if err != nil {
		return err
	}
	return nil
}

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
		if rows != nil {
			rows.Close()
		}
		if err != nil {
			return err
		}
	} else {
		// The app already exists, only add new packages for that app.
		bothPkgs := util.UniqAppend(app.Packages, dbPkgs)
		rows, err := db.Query("UPDATE app_packages SET packages = $1 WHERE id = $2",
			pq.Array(&bothPkgs), app.DBID)
		if rows != nil {
			rows.Close()
		}
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
	err := db.QueryRow("SELECT permissions FROM app_perms WHERE id = $1", app.DBID).
		Scan(pq.Array(&dbPerms))
	if err != nil {
		if err != sql.ErrNoRows {
			return err
		}
		rows, err := db.Query("INSERT INTO app_perms VALUES ($1, $2)",
			app.DBID, pq.Array(&sPerms))
		if rows != nil {
			rows.Close()
		}
		if err != nil {
			return err
		}
	} else {
		bothPerms := util.UniqAppend(sPerms, dbPerms)
		rows, err := db.Query("UPDATE app_perms SET permissions = $1 WHERE id = $2",
			pq.Array(&bothPerms), app.DBID)
		if rows != nil {
			rows.Close()
		}
		if err != nil {
			return err
		}
	}

	return nil
}

// SetIcon is a function that sets the icon field of the DB.
func SetIcon(id int64, icon string) error {
	if !useDB || id == 0 {
		return nil
	}

	rows, err := db.Query("UPDATE app_versions SET icon = $1 WHERE id = $2", icon, id)
	if rows != nil {
		rows.Close()
	}
	return err
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
		if rows != nil {
			rows.Close()
		}
		if err != nil {
			return err
		}
	} else {
		bothHosts := util.UniqAppend(hosts, dbHosts)
		rows, err := db.Query("UPDATE app_hosts SET hosts = $1 WHERE id = $2",
			pq.Array(&bothHosts), app.DBID)
		if rows != nil {
			rows.Close()
		}
		if err != nil {
			return err
		}
	}

	return nil
}

// SetReflect sets the value of uses_reflect for an app version
func SetReflect(id int64, val bool) error {
	rows, err := db.Query("UPDATE app_versions SET uses_reflect = $1 WHERE id = $2", val, id)
	if rows != nil {
		rows.Close()
	}
	return err
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
		&appVer)

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
	if rows != nil {
		defer rows.Close()
	}
	if err != nil {
		return []Developer{}, err
	}
	ret := make([]Developer, 0, num)
	for i := 0; rows.Next(); i++ {
		ret = append(ret, Developer{})
		rows.Scan(
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
	if rows != nil {
		defer rows.Close()
	}
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
	if rows != nil {
		defer rows.Close()
	}
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
	if rows != nil {
		defer rows.Close()
	}
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

// GetAltApps takes an app's DB ID and returns a collection of
// alternative apps for the specified app
func GetAltApps(appID string) ([]AltApp, error) {
	rows, err := db.Query("SELECT * FROM alt_apps alt WHERE alt.id = $1", appID)
	if rows != nil {
		defer rows.Close()
	}
	if err != nil {
		return []AltApp{}, err
	}
	var result []AltApp
	for i := 0; rows.Next(); i++ {
		var altApp AltApp
		rows.Scan(&altApp.ID, &altApp.Title, &altApp.URL)
		result = append(result, altApp)
	}

	if rows.Err() != sql.ErrNoRows && rows.Err() != nil {
		fmt.Println("Database err", rows.Err())
	}

	return result, nil
}

// precentifyArray produces a Postgres compatable 'like any' string intended
// to be used as part of a larger query.
func percentifyArray(arr []string) string {

	partQuery := " (Array[ "

	for i, param := range arr {
		partQuery += "LOWER('%" + param + "%')"
		if i < len(arr)-1 {
			partQuery += ","
		}

	}

	partQuery += "])"

	return partQuery
}

// QuickQuery depricates all of dean's queries.
//
//
func QuickQuery(
	fullDetails bool, appStore string, limit string, offset string, developers []string,
	genres []string, permissions []string, appIDs []string, titles []string,
) ([]AppVersion, error) {

	var storestart string
	storestart = "SELECT " +
		"a.id," +
		"a.title," +
		"a.summary," +
		"a.description," +
		"a.store_url," +
		"a.price," +
		"a.free," +
		"a.rating," +
		"a.num_reviews," +
		"a.genre," +
		"a.family_genre," +
		"a.min_installs," +
		"a.max_installs," +
		"a.updated," +
		"a.android_ver," +
		"a.content_rating," +
		"a.recent_changes," +
		"v.app," +
		"v.store," +
		"v.region," +
		"v.version," +
		"v.icon," +
		"d.email," +
		"d.name," +
		"d.store_site," +
		"d.site," +
		"h.hosts," +
		"p.permissions," +
		"pkg.packages"
		//"app_perms.permissions," + "packages" +

	tableQuery := " FROM " +
		appStore +
		" a FULL OUTER JOIN app_versions v ON (a.id = v.id) " +
		" FULL OUTER JOIN developers d ON (a.developer = d.id) " +
		" FULL OUTER JOIN app_hosts h ON (a.id = h.id) " +
		" FULL OUTER JOIN app_perms p ON (a.id = p.id) " +
		" FULL OUTER JOIN app_packages pkg  ON (a.id = pkg.id) "

	//Table Join Appends
	//+ "NATURAL JOIN app_perms " + "NATURAL JOIN app_packages"
	structuredQuery := storestart + tableQuery +
		" WHERE LOWER(a.title) LIKE any " + percentifyArray(titles) +
		" AND LOWER(d.name) LIKE any " + percentifyArray(developers) +
		" AND LOWER(a.genre) LIKE any " + percentifyArray(genres) +
		//" AND LOWER(app_perms.permissions) like any " + percentifyArray(permissions) + //TODO: s a array so need to check the arrays...
		" AND LOWER(v.app) LIKE any " + percentifyArray(appIDs) +
		" LIMIT " + limit + " OFFSET " + offset

	println(structuredQuery)

	rows, err := db.Query(structuredQuery)

	defer rows.Close()

	if err != nil {
		return []AppVersion{}, err
	}

	var result []AppVersion
	for i := 0; rows.Next(); i++ {

		var appData AppVersion
		var playInf PlayStoreInfo

		//Potential null values
		var summ sql.NullString
		var desc sql.NullString
		var genre sql.NullString
		var famGenre sql.NullString
		var video sql.NullString
		var icon sql.NullString
		app := sql.NullString{}
		store := sql.NullString{}
		region := sql.NullString{}
		ver := sql.NullString{}
		devStoreSite, devSite := sql.NullString{}, sql.NullString{}
		//var perms []string
		//var packages []string
		hosts, perms, pkgs := []sql.NullString{}, []sql.NullString{}, []sql.NullString{}
		var err error
		//Cannot just cast straight into types because of the postgre type conversion
		err = rows.Scan(
			&appData.ID,
			&playInf.Title,
			&summ,
			&desc,
			&playInf.StoreURL,
			&playInf.Price,
			&playInf.Free,
			&playInf.Rating,
			&playInf.NumReviews,
			&genre,
			&famGenre,
			&playInf.Installs.Min,
			&playInf.Installs.Max,
			&playInf.Updated,
			&playInf.AndroidVer,
			&playInf.ContentRating,
			pq.Array(&playInf.RecentChanges),
			&app,
			&store,
			&region,
			&ver,
			&icon, //&appData.Icon,
			pq.Array(&appData.Dev.Emails),
			&appData.Dev.Name,
			&devStoreSite,    //&appData.Dev.StoreSite,
			&devSite,         //&appData.Dev.Site,
			pq.Array(&hosts), //pq.Array(&appData.Hosts),
			pq.Array(&perms), //pq.Array(&appData.Perms),
			pq.Array(&pkgs))  //pq.Array(&appData.Packages))
		if err != nil {
			fmt.Println("Database Query", err)
		} else {

			//info.FamilyGenre = famGenre.Valid ? famGenre.String : ""
			//TODO: ight be able to get away with nul being "" after the scan stage
			playInf.Summary = summ.String
			playInf.Description = desc.String
			playInf.Genre = genre.String
			playInf.Video = video.String
			playInf.FamilyGenre = famGenre.String
			appData.Icon = icon.String
			appData.StoreInfo = playInf
			appData.Dev.StoreSite = devStoreSite.String
			appData.Dev.Site = devSite.String
			appData.App = app.String
			appData.Store = store.String
			appData.Region = region.String
			appData.Ver = ver.String
			for _, host := range hosts {
				appData.Hosts = append(appData.Hosts, host.String)
			}
			for _, perm := range perms {
				appData.Perms = append(appData.Perms, perm.String)
			}
			for _, pkg := range pkgs {
				appData.Packages = append(appData.Packages, pkg.String)
			}
			result = append(result, appData)
		}
	}

	if rows.Err() != sql.ErrNoRows && rows.Err() != nil {
		fmt.Println("Database err", rows.Err())
	}

	return result, nil

}

// GetAppsToAnalyze returns a list of up to 10 apps that have analyzed=False and
// downloaded=True for the analyzer.
func GetAppsToAnalyze() ([]AppVersion, error) {
	rows, err := db.Query("SELECT v.id, v.app,v.store, v.region, v.version, v.screen_flags, v.icon FROM app_versions v FULL OUTER JOIN playstore_apps p ON (v.id = p.id) WHERE v.analyzed = False AND v.downloaded = True ORDER BY v.last_analyze_attempt NULLS FIRST, p.max_installs USING> LIMIT 10")
	if rows != nil {
		defer rows.Close()
	}
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
	if rows != nil {
		rows.Close()
	}
	return err
}

// SetAnalyzed sets analyzed=True for a given app.
func SetAnalyzed(id int64) error {
	rows, err := db.Query("UPDATE app_versions SET analyzed = True WHERE id = $1", id)
	if rows != nil {
		rows.Close()
	}
	return err
}
