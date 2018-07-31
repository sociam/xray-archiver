package db

import (
	"database/sql"
	"fmt"
	"strconv"
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

// GetAppHostsByID selects an app host record from the DB using the provided ID
func GetAppHostsByID(id int64) (util.AppHostRecord, error) {
	var appHosts util.AppHostRecord

	util.Log.Debug("Requesting App Host info for App with ID: %d", id)
	db.QueryRow("select id, hosts from app_hosts where id = $1", id).Scan(
		&appHosts.ID,
		pq.Array(&appHosts.HostNames))

	util.Log.Debug("Finished Selecting app_hosts record for id: %d. Returning AppHostsRecord Object.", id)

	return appHosts, nil
}

// IncrementCompanyAppAssociationCount increments the counter on the associated app and company
func IncrementCompanyAppAssociationCount(appID int64, companyName string) error {
	if !HasAppVersionID(appID) {
		util.Log.Warning("App ID: %d Not Found", appID)
		return nil
	}

	if !HasCompanyName(companyName) {
		util.Log.Warning("Company Name %s Not Found", companyName)
		return nil
	}

	if !HasCompanyAppAssociation(appID, companyName) {
		util.Log.Debug("Company-App association between app: %d and company: %s Not Found", appID, companyName)
		return nil
	}

	rows, err := db.Query(
		"update companyappassociations set number_of_associations = number_of_associations + 1 where company_name=$1 and associated_app=$2",
		companyName,
		appID)
	rows.Close()

	if err != nil {
		util.Log.Err("Error incrementing number of associations for companyAppAssociation between Company: %s and app with ID: %d", companyName, appID, err)
	}

	return nil
}

// InsertCompanyAppAssociation inserts an app and company name association into the database.
func InsertCompanyAppAssociation(appID int64, companyName string) error {

	if !HasAppVersionID(appID) {
		util.Log.Warning("App ID: %d Not Found", appID)
		return nil
	}

	if !HasCompanyName(companyName) {
		util.Log.Warning("Company Name %s Not Found", companyName)
		return nil
	}

	if HasCompanyAppAssociation(appID, companyName) {
		util.Log.Debug("Company-App association between app: %d and company: %s already exists. Incrementing Count", appID, companyName)
		return IncrementCompanyAppAssociationCount(appID, companyName)
	}

	rows, err := db.Query("insert into companyAppAssociations(company_name, associated_app, number_of_associations) values($1,$2,1)", companyName, appID)
	rows.Close()

	if err != nil {
		util.Log.Err("Error inserting company-app association for app with id: %d and company with name: %s. Error:", appID, companyName, err)
		return err
	}

	return nil
}

// HasCompanyAppAssociation checks if an association between a given app and company name already exists.
func HasCompanyAppAssociation(appID int64, companyName string) bool {
	var assocCount int
	util.Log.Debug("Counting number of CompanyApp Assocications to see if it already exists.")
	db.QueryRow("select count(*) from companyAppAssociations where company_name=$1 and associated_app=$2", companyName, appID).Scan(
		&assocCount)
	util.Log.Debug("Company-App Associations Counted. %d associations found for app with id: %d and company with name: %s", assocCount, appID, companyName)
	return assocCount > 0
}

// InsertCompanyName inserts the provided company name into the database.
func InsertCompanyName(companyName string) error {
	if HasCompanyName(companyName) {
		util.Log.Warning("Company Name: %s already exists. Exiting method.")
		return nil
	}

	rows, err := db.Query("insert into companyNames(company_name) values( $1 )", companyName)
	rows.Close()

	if err != nil {
		util.Log.Err("Error inserting Company Name: %s into the companyNames table. Error: %s", companyName, err)
		return err
	}

	return nil
}

// HasCompanyName Checks if companyNames table has the provided company name
func HasCompanyName(companyName string) bool {
	var companyCount int
	util.Log.Debug("Requesting CompanyName Row for provided company name: %s", companyName)
	db.QueryRow("select count(*) from companyNames where company_name=$1", companyName).Scan(
		&companyCount)
	util.Log.Debug("CompanyNames counted. %d companies found with name matching %s", companyCount, companyName)
	return companyCount > 0
}

// HasAppVersionID Checks if app_versions table has the provided appversionId
func HasAppVersionID(appID int64) bool {
	var appVerCount int
	util.Log.Debug("Requesting app_version Row for provided app version id: %d", appID)
	db.QueryRow("select count(*) from app_versions where id = $1", appID).Scan(
		&appVerCount)
	util.Log.Debug("app_versions counted. %d app_versions found with ID matching %d", appVerCount, appID)
	return appVerCount > 0
}

// GetAppHostIDs returns an array of  app_version ids found in app_hosts.
func GetAppHostIDs() ([]int64, error) {
	ids := make([]int64, 0, 0)

	util.Log.Debug("About To request all app_host IDs.")
	rows, err := db.Query("SELECT id FROM app_hosts")

	if rows != nil {
		util.Log.Debug("Rows successfully Selected.")
		defer rows.Close()
	}
	if err != nil {
		util.Log.Err("Error Occured querying app_hosts for ids.")
		return []int64{}, err
	}

	util.Log.Debug("Scanning app_host ID rows.")
	for i := 0; rows.Next(); i++ {
		ids = append(ids, int64(i))
		rows.Scan(&ids[i])
	}

	if rows.Err() != sql.ErrNoRows && rows.Err() != nil {
		util.Log.Err("Error querying for App IDs")
		return []int64{}, rows.Err()
	}

	util.Log.Debug("Finished Processing App IDs. Returning Array of AppIDs.")
	return ids, nil
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

	if rows.Err() != sql.ErrNoRows && rows.Err() != nil {
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

// GetAppCompanyFreq Queries the DB for global company stats.
// a generalised get Stat Freq table method should be created where
// the name of a freq table is specified. the tables would take a
// standardised format.
func GetAppCompanyFreq() ([]CompanyCoverage, error) {
	rows, err := db.Query("SELECT * FROM company_app_coverage")
	results := []CompanyCoverage{}
	if rows != nil {
		defer rows.Close()
	}
	if err != nil {
		util.Log.Err("Error in stats table query", err)
		return results, err
	}
	util.Log.Debug("Scanning Rows.")

	for i := 0; rows.Next(); i++ {
		row := CompanyCoverage{}
		rows.Scan(
			&row.Company,
			&row.Type,
			&row.AppCount,
			&row.TotalApps,
			&row.CompanyFreq)
		results = append(results, row)
		//util.Log.Debug("Row Scanned: " + fmt.Sprint(i))
	}
	util.Log.Debug("Rows Scanned")
	if rows.Err() != sql.ErrNoRows && rows.Err() != nil {
		util.Log.Err("Error Scanning Rows ", rows.Err())
		return []CompanyCoverage{}, err
	}
	util.Log.Debug("Returning Rows")
	return results, nil
}

// GetCompanyGenreCoverage queries DB for stats.
func GetCompanyGenreCoverage() ([]CompanyGenreCoverage, error) {
	rows, err := db.Query("SELECT * FROM company_genre_coverage")
	results := []CompanyGenreCoverage{}
	if rows != nil {
		defer rows.Close()
	}
	if err != nil {
		util.Log.Err("Error in selecting from company_genre_coverage")
		return results, err
	}
	util.Log.Debug("Scanning Rows for company_genre_coverage")

	for i := 0; rows.Next(); i++ {
		row := CompanyGenreCoverage{}
		rows.Scan(
			&row.Company,
			&row.CompanyCount,
			&row.Genre,
			&row.GenreTotal,
			&row.CoveragePct)
		results = append(results, row)
	}
	util.Log.Debug("Rows Scanned")
	if rows.Err() != sql.ErrNoRows && rows.Err() != nil {
		util.Log.Err("Error Scanning Rows ", rows.Err())
		return []CompanyGenreCoverage{}, err
	}
	util.Log.Debug("Returning Rows")
	return results, nil
}

// GetAppTypeFreq Querys the DB for rows in the app_types_coverage stats table
func GetAppTypeFreq() ([]CompanyTypeCoverage, error) {
	rows, err := db.Query("SELECT * FROM app_type_coverage")
	results := []CompanyTypeCoverage{}
	if rows != nil {
		defer rows.Close()
	}
	if err != nil {
		util.Log.Err("Error in stats table query", err)
		return results, err
	}
	util.Log.Debug("Scanning Rows.")

	for i := 0; rows.Next(); i++ {
		row := CompanyTypeCoverage{}
		rows.Scan(
			&row.Type,
			&row.AppCount,
			&row.TotalApps,
			&row.TypeFreq)
		results = append(results, row)
		//util.Log.Debug("Row Scanned: " + fmt.Sprint(i))
	}
	util.Log.Debug("Rows Scanned")
	if rows.Err() != sql.ErrNoRows && rows.Err() != nil {
		util.Log.Err("Error Scanning Rows ", rows.Err())
		return []CompanyTypeCoverage{}, err
	}
	util.Log.Debug("Returning Rows")
	return results, nil
}

// GetGenreHostAverages Queries DB for Genre averages stats view.
func GetGenreHostAverages() ([]GenreStats, error) {
	rows, err := db.Query("SELECT * FROM genre_host_averages")
	results := []GenreStats{}
	if rows != nil {
		defer rows.Close()
	}
	if err != nil {
		util.Log.Err("Error in stats table query", err)
		return results, err
	}
	util.Log.Debug("Scanning Rows.")

	for i := 0; rows.Next(); i++ {
		row := GenreStats{}
		rows.Scan(
			&row.Category,
			&row.HostCount,
			&row.AppCount,
			&row.GenreAvg)
		results = append(results, row)
		//util.Log.Debug("Row Scanned: " + fmt.Sprint(i))
	}
	util.Log.Debug("Rows Scanned")
	if rows.Err() != sql.ErrNoRows && rows.Err() != nil {
		util.Log.Err("Error Scanning Rows ", rows.Err())
		return []GenreStats{}, err
	}
	util.Log.Debug("Returning Rows")
	return results, nil
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
			util.Log.Err("Database err:", err)
		} else {
			ret = append(ret, app)
		}
	}

	if rows.Err() != sql.ErrNoRows && rows.Err() != nil {
		util.Log.Err("Databse err", rows.Err())
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
		util.Log.Err("Databse err", rows.Err())
		return []App{}, rows.Err()
	}

	return ret, nil
}

//GetManualAltApps - Returns app ids that are logged as alternatives to the one given.
func GetManualAltApps(appID string) ([]string, error) {
	rows, err := db.Query("SELECT alt_id FROM manual_alts WHERE source_id = $1", appID)

	if rows != nil {
		defer rows.Close()
	}

	if err != nil {
		fmt.Println("Error. returning an empty Alt app. ")
		fmt.Println(err)
		return []string{}, err
	}

	result := []string{}

	for i := 0; rows.Next(); i++ {
		str := sql.NullString{}

		err = rows.Scan(&str)

		if err != nil {
			fmt.Println(err)
		} else {
			result = append(result, str.String)
		}
	}

	if rows.Err() != sql.ErrNoRows && rows.Err() != nil {
		util.Log.Err("Database err", rows.Err())
	}

	fmt.Println(fmt.Sprint(len(result)) + " rows found")
	return result, nil
}

// GetAltApps takes an app's DB ID and returns a collection of
// alternative apps for the specified app - For the API
func GetAltApps(appID string) ([]AltApp, error) {
	rows, err := db.Query("SELECT app_id, alt_app_title, alt_to_url, g_play_url, g_play_id, icon_url, official_site_url, is_scraped FROM alt_apps WHERE app_id = $1", appID)

	if rows != nil {
		defer rows.Close()
	}

	if err != nil {
		fmt.Println("Error. returning an empty Alt app. ")
		fmt.Println(err)
		return []AltApp{}, err
	}

	result := []AltApp{}

	for i := 0; rows.Next(); i++ {
		var altApp AltApp
		// initialising nullable strings for populating from database
		var AppID, AltAppTitle, AltToURL, GPlayURL, GPlayID, IconURL, OfficialSiteURL sql.NullString

		// Scanning from database into alt app object / nullable strings
		err = rows.Scan(
			&AppID,

			&AltAppTitle,
			&AltToURL,
			&GPlayURL,
			&GPlayID,
			&IconURL,
			&OfficialSiteURL,

			&altApp.IsScraped)

		// putting nullable strings into alt app object.
		altApp.AltAppTitle = AltAppTitle.String
		altApp.AltToURL = AltToURL.String
		altApp.GPlayURL = GPlayURL.String
		altApp.GPlayID = GPlayID.String
		altApp.IconURL = IconURL.String
		altApp.OfficialSiteURL = OfficialSiteURL.String
		fmt.Println("Alt App Fetched from DB: " + AltAppTitle.String)
		if err != nil {
			fmt.Println(err)
		} else {
			result = append(result, altApp)
		}
	}

	if rows.Err() != sql.ErrNoRows && rows.Err() != nil {
		util.Log.Err("Database err", rows.Err())
	}

	fmt.Println(fmt.Sprint(len(result)) + " rows found")
	return result, nil
}

func asPGAny(arr *[]string) {
	for i, v := range *arr {
		(*arr)[i] = "%" + v + "%"
	}
}

func asPGStartsWith(arr *[]string) {
	for i, v := range *arr {
		(*arr)[i] = v + "%"
	}
}

var appStoreTable = map[string]string{
	"play": "playstore_apps",
}

func extendWhereQuery(querystr *string, colName string, numParam *int, arr *[]string, hasPrev *bool) {
	util.Log.Debug("Attempting to extend query params %s \n", *arr)

	if *hasPrev {
		*querystr += "AND "
	} else {
		*hasPrev = true
	}

	newQuery := colName + "ILIKE ANY($" + strconv.Itoa(*numParam) + ") "
	*querystr += newQuery

	asPGAny(arr)

	*numParam++
}

// QueryAll depricates all of dean's queries.
func QueryAll(
	onlyAnalyzed bool, appStore string, limit string, offset string, developers []string,
	genres []string, permissions []string, appIDs []string, titles []string, startsWith []string,
) ([]AppVersion, error) {

	var queryStr string
	//TODO: Left as join will be built later
	// queryStr = "SELECT " +
	// 	"a.id, a.title, a.summary, a.description, a.store_url, a.price, a.free, a.rating, " +
	// 	"a.num_reviews, a.genre, a.family_genre, a.min_installs, a.max_installs, a.updated, " +
	// 	"a.android_ver, a.content_rating, a.recent_changes, v.app, v.store, v.region, " +
	// 	"v.version, v.icon, v.analyzed, d.email, d.name, d.store_site, d.site, h.hosts, p.permissions, " +
	// 	"pkg.packages " +
	// 	"FROM " + appStoreTable[appStore] + " a " +
	// 	"FULL OUTER JOIN app_versions v ON (a.id = v.id) " +
	// 	"FULL OUTER JOIN developers d ON (a.developer = d.id) " +
	// 	"FULL OUTER JOIN app_hosts h ON (a.id = h.id) " +
	// 	"FULL OUTER JOIN app_perms p ON (a.id = p.id) " +
	// 	"FULL OUTER JOIN app_packages pkg  ON (a.id = pkg.id) "

	queryStr = "SELECT " +
		"id, title, summary, description, store_url, price, free, rating, " +
		"num_reviews, genre, family_genre, min_installs, max_installs, updated, " +
		"android_ver, content_rating, recent_changes, app, store, region, " +
		"version, icon, analyzed, email, name, store_site, site, hosts, permissions, " +
		"packages " +
		"FROM  mat_view.apps_play_data " +
		"WHERE "

	var args []interface{}

	numParam := 1
	hasPrev := false

	//TODO: future me will fix this later... however considering the horrible *quick*query a better refactor is needed...
	//This better refcator will sepeate each select as it's own component. Get ids then get where matching. Would require
	//multiple joins but if the views are setup should not matter
	util.Log.Debug("Requested where  titles %s \n", titles)
	if len(titles) > 0 {
		extendWhereQuery(&queryStr, "title ", &numParam, &titles, &hasPrev)
		args = append(args, pq.Array(&titles))
	}

	util.Log.Debug("Requested where  developers %s \n", developers)
	if len(developers) > 0 {
		extendWhereQuery(&queryStr, "name ", &numParam, &developers, &hasPrev)
		args = append(args, pq.Array(&developers))
	}

	util.Log.Debug("Requested where genres %s \n", genres)
	if len(genres) > 0 {
		util.Log.Debug("Attempting to query params %s \n", genres)

		if hasPrev {
			queryStr += "AND "
		} else {
			hasPrev = true
		}

		newQuery := "app IN( "

		for i := numParam; i < numParam+len(genres)-1; numParam++ {
			newQuery += "$" + strconv.Itoa(i) + ", "
			util.Log.Debug("Adding arg ", genres[i])
			args = append(args, genres[i])
		}

		newQuery += "$" + strconv.Itoa(numParam) + ") "
		numParam++
		util.Log.Debug("Adding arg ", genres[len(genres)-1])
		args = append(args, genres[len(genres)-1])

		queryStr += newQuery
	}

	util.Log.Debug("Requested where genres %s \n", appIDs)
	if len(appIDs) > 0 {
		util.Log.Debug("Attempting to query params %s \n", appIDs)

		if hasPrev {
			queryStr += "AND "
		} else {
			hasPrev = true
		}

		newQuery := "app IN( "

		for i := numParam; i < numParam+len(appIDs)-1; numParam++ {
			newQuery += "$" + strconv.Itoa(i) + ", "
			util.Log.Debug("Adding arg ", appIDs[i])
			args = append(args, appIDs[i])
		}

		newQuery += "$" + strconv.Itoa(numParam) + ") "
		numParam++
		util.Log.Debug("Adding arg ", appIDs[len(appIDs)-1])
		args = append(args, appIDs[len(appIDs)-1])

		queryStr += newQuery

	}

	util.Log.Debug("Requested where startsWith %s \n", startsWith)
	if len(startsWith) > 0 {
		util.Log.Debug("Attempting to query params %s \n", startsWith)

		if hasPrev {
			queryStr += "AND "
		} else {
			hasPrev = true
		}

		newQuery := "title ILIKE ANY($" + strconv.Itoa(numParam) + ") "
		queryStr += newQuery

		asPGStartsWith(&startsWith)

		numParam++

		args = append(args, pq.Array(&startsWith))
	}

	args = append(args, limit, offset)

	var shouldAnalyze string

	if onlyAnalyzed {
		shouldAnalyze = "AND analyzed = true "
	} else {
		shouldAnalyze = ""
	}

	queryStr += shouldAnalyze + " ORDER BY max_installs using> "
	queryStr += "LIMIT $" + strconv.Itoa(numParam) + " "
	numParam++
	queryStr += "OFFSET $" + strconv.Itoa(numParam)

	fmt.Printf("%s\n", queryStr)

	for _, arg := range args {
		util.Log.Debug("Args to postgres query: %v", arg)
	}

	rows, err := db.Query(queryStr, args...)

	if rows != nil {
		defer rows.Close()
	}
	if err != nil {
		util.Log.Debug("Failed to grab app query", err)
		return []AppVersion{}, err
	}

	util.Log.Debug("Examining rows")

	result := []AppVersion{}

	for i := 0; rows.Next(); i++ {
		util.Log.Debug("Casting Row: ", i)

		var appData AppVersion
		var playInf PlayStoreInfo

		//Potential null values
		var summ, desc, genre, famGenre, video, icon, devStoreSite, devSite sql.NullString
		//var perms []string
		//var packages []string
		hosts, perms, pkgs, recentChanges := []sql.NullString{}, []sql.NullString{}, []sql.NullString{}, []sql.NullString{}
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
			pq.Array(&recentChanges),
			&appData.App,
			&appData.Store,
			&appData.Region,
			&appData.Ver,
			&icon, //&appData.Icon,
			&appData.IsAnalyzed,
			pq.Array(&appData.Dev.Emails),
			&appData.Dev.Name,
			&devStoreSite,    //&appData.Dev.StoreSite,
			&devSite,         //&appData.Dev.Site,
			pq.Array(&hosts), //pq.Array(&appData.Hosts),
			pq.Array(&perms), //pq.Array(&appData.Perms),
			pq.Array(&pkgs))  //pq.Array(&appData.Packages))
		if err != nil {
			util.Log.Err("Database Query", err)
		} else {
			util.Log.Debug("Casting data into correct structures")
			playInf.Summary = summ.String
			playInf.Description = desc.String
			playInf.Genre = genre.String
			playInf.Video = video.String
			playInf.FamilyGenre = famGenre.String
			appData.Icon = icon.String
			appData.Dev.StoreSite = devStoreSite.String
			appData.Dev.Site = devSite.String

			for _, host := range hosts {
				appData.Hosts = append(appData.Hosts, host.String)
			}
			for _, perm := range perms {
				appData.Perms = append(appData.Perms, perm.String)
			}
			for _, pkg := range pkgs {
				appData.Packages = append(appData.Packages, pkg.String)
			}
			for _, change := range recentChanges {
				playInf.RecentChanges = append(playInf.RecentChanges, change.String)
			}

			appData.StoreInfo = playInf

			result = append(result, appData)
		}
	}

	util.Log.Debug("Done examining rows.")

	if rows.Err() != sql.ErrNoRows && rows.Err() != nil {
		util.Log.Err("Database err", rows.Err())
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
			util.Log.Err("Getting apps to analyzed errors", err)
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
