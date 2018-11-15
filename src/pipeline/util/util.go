package util

import (
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"net"
	"net/http"
	"net/url"
	"os"
	"os/exec"
	"path"
	"strings"
	"time"
	"bytes"
	"math"
//	"sync"
)

// Unit for maps in data.go
type Unit struct{}

var unit Unit

// App Struct for holding of information extracted from the APK
type App struct {
	DBID                   int64
	ID, Store, Region, Ver string
	Path, UnpackDir        string
	Perms                  []Permission
	Hosts                  []string
	Packages               []string
	Icon                   string
	UsesReflect            bool
	APKLocationUUID        string
	APKLocationPath        string
	APKLocationRoot        string
}

// Permission Struct represents the permission information found
// in an APK
type Permission struct {
	ID        string `xml:"name,attr"`
	MaxSdkVer string `xml:"maxSdkVersion,attr"`
}

// NewApp Constructs a new app. initialising values based on
// the parameters passed.
func NewApp(dbID int64, id, store, region, ver, apkLocationPath, apkLocationRoot, apkLocationUUID string) *App {
	return &App{
		DBID:            dbID,
		ID:              id,
		Store:           store,
		Region:          region,
		Ver:             ver,
		APKLocationPath: apkLocationPath,
		APKLocationRoot: apkLocationRoot,
		APKLocationUUID: apkLocationUUID}
}

// AppByPath returns an App object with the Path value initialised.
func AppByPath(path string) *App {
	return &App{Path: path}
}

// AppDir returns the directory of the apk and other misc files.
func (app *App) AppDir() string {
	if app.Path != "" {
		return path.Dir(app.Path)
	}
	return path.Dir(app.ApkPath())
}

// getUUIDMountPath uses linux 'findmnt' program to try and find the mount point of the provided uuid.
func getUUIDMountPath(UUID string) string {
	out, err := exec.Command("findmnt", "-rn", "-S", "UUID="+UUID, "-o", "TARGET").Output()
	if err != nil {
		log.Fatal(err)
	}
	return string(out[:])
}

// ApkPath will try and find the APK on any mounted disks in potential locations
// starting with the location specified in the DB, falling down to checking
// locations with the root of the path specidied in the DB substituted,
// follewed by checking each location in the config forming a path from the
// app version details.
func (app *App) ApkPath() string {
	fmt.Println("Getting APK Path for App:", app.ID)

	apkLocation := path.Join(app.APKLocationPath, app.ID+".apk")
	fmt.Println("Checking if APK is at: ", apkLocation)

	if _, err := os.Stat(apkLocation); err == nil {
		fmt.Println("App Found in DB specified location: ", apkLocation)
		return path.Join(path.Clean(app.APKLocationPath), app.ID+".apk")
	}

	// Check if the Root is wrong, get the filesystem mount path
	// and replace the APKLocationRoot of APKLocationPath with the new root.
	uuidMount := getUUIDMountPath(app.APKLocationUUID)
	apkLocation = path.Join(strings.Replace(app.APKLocationPath, app.APKLocationRoot, uuidMount, 1), app.ID+".apk")

	fmt.Println("Checking if APK is at: ", apkLocation)
	if _, err := os.Stat(apkLocation); err == nil {
		fmt.Println("App Found on DB specified Device. UUID: ", app.APKLocationUUID, "APK Path:", apkLocation)
		return path.Join(path.Clean(apkLocation), app.ID+".apk")
	}

	// if the app cannot be found in the new mount location for whatever UUID, go through
	// each storage location in the config and check there.
	for _, location := range Cfg.StorageConfig.APKDownloadDirectories {
		apkLocation = path.Join(strings.Replace(app.APKLocationPath, app.APKLocationRoot, location.Path, 1), app.ID+".apk")
		fmt.Println("Checking if APK is at: ", apkLocation)

		if _, err := os.Stat(apkLocation); err == nil {
			fmt.Println("Searched for APK and found it in: ", apkLocation)
			return path.Join(path.Clean(apkLocation), app.ID+".apk")
		}
	}

	if app.Path != "" {
		return app.Path
	}
	return path.Join(app.AppDir(), app.ID+".apk")
}

// OutDir specifies where Apps should be unpacked to. it also creates
// the directory structure for that path and returns the path as a
// string.
func (app *App) OutDir() string {
	if app.UnpackDir == "" {
		if app.Path != "" {
			var err error
			app.UnpackDir, err = ioutil.TempDir(Cfg.StorageConfig.APKUnpackDirectory, path.Base(app.Path))
			if err != nil {
				// maybe do something else?
				log.Fatal("Failed to create temp dir in ", Cfg.StorageConfig.APKUnpackDirectory, ": ", err)
			}
		} else {
			app.UnpackDir = path.Join(Cfg.StorageConfig.APKUnpackDirectory, app.ID, app.Store, app.Region, app.Ver)
			if err := os.MkdirAll(app.UnpackDir, 0755); err != nil {
				log.Fatalf("Failed to create temp dir in %s: %s", app.UnpackDir, err.Error())
			}
		}
	}
	return app.UnpackDir
}

// Unpack passes an app to apktool to disassemble an APK. the contents are
// stored in the path specified by OutDir.
func (app *App) Unpack() error {
	apkPath, outDir := app.ApkPath(), app.OutDir()
	if _, err := os.Stat(apkPath); err != nil {
		if os.IsNotExist(err) {
			return err
		}
		return fmt.Errorf("couldn't open apk %s: %s", apkPath, err.Error())
	}

	if err := os.MkdirAll(path.Dir(outDir), 0755); err != nil {
		return os.ErrPermission
	}
	now := time.Now()
	if err := os.Chtimes(path.Dir(outDir), now, now); err != nil {
		return os.ErrPermission
	}

	cmd := exec.Command("apktool", "d", "-s", apkPath, "-o", outDir, "-f")
	out, err := cmd.CombinedOutput()
	if err != nil {
		return fmt.Errorf("%s unpacking apk; output below:\n%s",
			err.Error(), string(out))
	}
	return nil
}

// Cleanup removes all directories specifed in an app object's OutDir.
func (app *App) Cleanup() error {
	return os.RemoveAll(app.OutDir())
}

// CheckDir verifies that a Dir is a Dir and exists.
func CheckDir(dir, name string) {
	fif, err := os.Stat(dir)
	if err != nil {
		if os.IsNotExist(err) {
			err = os.MkdirAll(dir, 0644)
			if err != nil {
				//TODO: something else
				panic(fmt.Sprintf("Couldn't create %s: %s", name, err.Error()))
			}
		} else {
			//TODO: something else
			panic(err)
		}
	} else if !fif.IsDir() {
		panic(fmt.Sprintf("%s isn't a directory!", name))
	}
}

// UniqAppend takes the contents of one array and adds any content
// not present in another array.
func UniqAppend(a []string, b []string) []string {
	ret := make([]string, 0, len(a)+len(b))
	for _, e := range a {
		ret = append(ret, e)
	}
	// set ret = a

	for _, be := range b {
		add := true
		for _, ae := range a {
			if ae == be {
				add = false
				break
			}
		}
		if add {
			ret = append(ret, be)
		}
	}
	return ret
}

/*
func uniqAppend(a []interface{}, b []interface{}) []interface{} {
	eMap = map[interface{}]Unit
	for _, e := range a {
		eMap[e] := unit
	}
	for _, e := range b {
		eMap[e] := unit
	}

	ret := make([]interface{}, 0, len(eMap))
	for e, _ := range eMap {
		ret := append(ret, e)
	}
	return ret
}
*/

// Dedup deduplicates a slice
func Dedup(a []string) []string {
	length := len(a) - 1
	for i := 0; i < length; i++ {
		for j := i + 1; j <= length; j++ {
			if a[i] == a[j] {
				a[j] = a[length]
				a = a[:length]
				length--
				j--
			}
		}
	}
	return a
}

// Combine puts together two maps of string keys and unit values.
func Combine(a, b map[string]Unit) map[string]Unit {
	ret := a
	for e := range b {
		ret[e] = unit
	}
	return ret
}

// StrMap creates a map of strings and units.
func StrMap(args ...string) map[string]Unit {
	ret := make(map[string]Unit)
	for _, e := range args {
		ret[e] = unit
	}

	return ret
}

// WriteJSON writes and encodes json dat.
func WriteJSON(w io.Writer, data interface{}) error {
	enc := json.NewEncoder(w)
	enc.SetEscapeHTML(false)
	enc.SetIndent("", "")
	return enc.Encode(data)
}

// WriteDEAN Writes and Encodes a 'Nah Mate'.
func WriteDEAN(w io.Writer, data interface{}) error {
	w.Write([]byte("Nah\n"))
	WriteJSON(w, data)
	w.Write([]byte("mate."))
	return nil
}

// GetJSON from valid url string gets json
func GetJSON(url string, target interface{}) error {
	client := &http.Client{Timeout: 10 * time.Second}
	r, err := client.Get(url)
	if err != nil {
		return err
	}
	if r.StatusCode != http.StatusOK {
		return fmt.Errorf("Got status %d while attempting to get GeoIP data", r.StatusCode)
	}
	defer r.Body.Close()
	return json.NewDecoder(r.Body).Decode(target)
}

func PostJSON(url string, body []byte, target interface{}) error {
     fmt.Println("Calling POST on url ", url)
        req, err := http.NewRequest("POST", url, bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")
	client := &http.Client{Timeout: 10 * time.Second}
	r, err := client.Do(req)
	if err != nil {  return err}
	if r.StatusCode != http.StatusOK {
	   return fmt.Errorf("Got status %d while attempting to get GeoIP data", r.StatusCode)
	}
	defer r.Body.Close()
	// debugging only 
	// if r.StatusCode == http.StatusOK {
	//    bodyBytes, err2 := ioutil.ReadAll(r.Body)
	//    bodyString := string(bodyBytes)
	//    fmt.Println("response body ", bodyString, err2);
	//    // 
	// }
	return json.NewDecoder(r.Body).Decode(target)
	// return nil; // Error{err:"bad"};	
}


// GeoIPInfo stores apphosts data for geolocation

// Original from freegeoip.net
type OldGeoIPInfo struct {
 	IP          string  `json:"ip"`
 	CountryCode string  `json:"country_code"`
 	CountryName string  `json:"country_name"`
 	RegionCode  string  `json:"region_code"`
 	RegionName  string  `json:"region_name"`
 	City        string  `json:"city"`
 	ZipCode     string  `json:"zip_code"`
 	TimeZone    string  `json:"time_zone"`
 	Latitude    float64 `json:"latitude"`
 	Longitude   float64 `json:"longitude"`
 	MetroCode   int     `json:"metro_code"`
}

// this struct is designed for ip-api.com : 
// please see http://ip-api.com/docs/api:json
type GeoIPInfo struct {
	IP          string  `json:"query"`
	CountryCode string  `json:"countryCode"`
	CountryName string  `json:"country"`
	RegionCode  string  `json:"region"`
	RegionName  string  `json:"regionName"`
	City        string  `json:"city"`
	ZipCode     string  `json:"zip"`
	TimeZone    string  `json:"timezone"`
	Latitude    float64 `json:"lat"`
	Longitude   float64 `json:"lon"`
	Org	    string  `json:"org"`
	
}

type IpAPIQuery struct {
     Query string `json:"query"`
}

// GetHostGeoIP grabs geo location information from hostname
func GetHostGeoIPSub100(host string) ([]GeoIPInfo, error) {

        geoipHost := Cfg.GeoIPEndpoint
	hosts, err := net.LookupHost(host)
	if err != nil {
		return nil, err
	}

	hosts = hosts[:100]
	queries := make([]IpAPIQuery,len(hosts))
	
	for i, host := range hosts {
	    fmt.Println("Marshalling: ", host, i, len(queries), len(hosts))
	    queries[i] = IpAPIQuery{url.PathEscape(host)}
	}
	
	jsons,err := json.Marshal(queries)
	var inf []GeoIPInfo
	err = PostJSON(geoipHost+"/batch", jsons, &inf)

	if err != nil {
	   fmt.Println("Couldn't call batch query ", err.Error())
	   return nil, err
	}
		
	return inf, nil	
}

func MakeRange(min, max int) []int {
    a := make([]int, max-min+1)
    for i := range a {
        a[i] = min + i
    }
    return a
}

func convertGeoStruct(g GeoIPInfo) OldGeoIPInfo {
     return OldGeoIPInfo{
     	    g.IP,
	    g.CountryCode,
	    g.CountryName,
	    g.RegionCode,
	    g.RegionName,
	    g.City,
	    g.ZipCode,
	    g.TimeZone,
	    g.Latitude,
	    g.Longitude,
	    0}
}

// GetHostGeoIP grabs geo location information from hostname
func GetHostGeoIPs(hosts []string) (map[string][]OldGeoIPInfo, error) {

     fmt.Println("Entering GHGIPs", hosts)

     max_hosts := 100
     geoipHost := Cfg.GeoIPEndpoint

     addr_h_map := make(map[string]string)
     addr_geo_map := make(map[string][]OldGeoIPInfo)
     
     addrs := make([]string, 0, len(hosts))

     // STEP 0: resolve hostnames -> addrs
     for _, host := range hosts {
         haddr, err := net.LookupHost(host)
     	 if err != nil {
	    fmt.Println("Error with ", host, " skipping ->");
	 } else {
	   if (len(haddr) > 0) {
	      addr_h_map[haddr[0]] = host
	      addrs = append(addrs, haddr[0])
	   }
	 }
     }

     fmt.Println(" Initial DNS resolve ", addr_h_map, len(addrs) )
     
     ret := make([]GeoIPInfo, 0, len(addrs)) // to return
     queries := make([]IpAPIQuery,len(addrs)) // prepare qureis array
     for i, addr := range addrs {
         queries[i] = IpAPIQuery{url.PathEscape(addr)}
     }

     for _, step := range MakeRange(0,len(addrs)/max_hosts) {
        start_idx := step*max_hosts
	end_idx := int(math.Min(float64((step+1)*max_hosts), float64(len(queries))))
	// fmt.Println(" start::end ", start_idx, "::", end_idx) 
	slice_queries := queries[start_idx:end_idx]
	// fmt.Println(" slice_query ", slice_queries, len(slice_queries)) 
	jsons,err := json.Marshal(slice_queries)
	var inf []GeoIPInfo
	err = PostJSON(geoipHost+"/batch", jsons, &inf)

	if err != nil {
	   fmt.Println("Couldn't call batch query ", err.Error())
	   // return nil, err
	} else {
	  ret = append(ret, inf...)
	}
     }

     // return host to geo map which will make EDI/Refine happy
     // as old style
     for _, r := range ret {
     	 host,ok := addr_h_map[r.IP]
	 if (ok) {
	    prev_rs, gotem := addr_geo_map[host]
	    if (!gotem) {
	       prev_rs = []OldGeoIPInfo{ convertGeoStruct(r) }
	    } else {
	       prev_rs = append(prev_rs, convertGeoStruct(r) )
	    }
	    addr_geo_map[host] = prev_rs
	} else {
	  fmt.Println("not ok ", host);
	 }
     }

     fmt.Println("Returning ", addr_geo_map)
	
     return addr_geo_map, nil
}



