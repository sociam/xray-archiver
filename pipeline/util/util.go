package util

import (
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"os"
	"os/exec"
	"path"
	"time"
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
}

// Permission Struct represents the permission information found
// in an APK
type Permission struct {
	ID        string `xml:"name,attr"`
	MaxSdkVer string `xml:"maxSdkVersion,attr"`
}

// NewApp Constructs a new app. initialising values based on
// the parameters passed.
func NewApp(dbID int64, id, store, region, ver string) *App {
	return &App{DBID: dbID, ID: id, Store: store, Region: region, Ver: ver}
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
	return path.Join(Cfg.AppDir, app.ID, app.Store, app.Region, app.Ver)
}

// ApkPath creates a string that represents the location of the APK
// on disk. Used to populate the Path string of an App object.
func (app *App) ApkPath() string {
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
			app.UnpackDir, err = ioutil.TempDir(Cfg.UnpackDir, path.Base(app.Path))
			if err != nil {
				// maybe do something else?
				log.Fatal("Failed to create temp dir in ", Cfg.UnpackDir, ": ", err)
			}
		} else {
			app.UnpackDir = path.Join(Cfg.UnpackDir, app.ID, app.Store, app.Region, app.Ver)
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

//Logger For Systemd logs
var prefixes = []string{"<0>", "<1>", "<2>", "<3>", "<4>", "<5>", "<6>", "<7>"}

type DebugLevel int

const (
	EMERG DebugLevel = iota
	ALERT
	CRIT
	ERR
	WARNING
	NOTICE
	INFO
	DEBUG
)

// func log(level DebugLevel, args ...string) {
// 	fmt.Println(prefixes[level], args)
// }

type logger struct{}

var Log = logger{}

func (l logger) Log(str string, args ...interface{}) {
	fmt.Printf(str, args...)
}

func (l logger) Emerg(str string, args ...interface{}) {
	fmt.Printf(prefixes[EMERG]+str, args...)
}

func (l logger) ALERT(str string, args ...interface{}) {
	fmt.Printf(prefixes[ALERT]+str, args...)
}

func (l logger) CRIT(str string, args ...interface{}) {
	fmt.Printf(prefixes[CRIT]+str, args...)
}

func (l logger) ERR(str string, args ...interface{}) {
	fmt.Printf(prefixes[ERR]+str, args...)
}

func (l logger) WARNING(str string, args ...interface{}) {
	fmt.Printf(prefixes[WARNING]+str, args...)
}

func (l logger) NOTICE(str string, args ...interface{}) {
	fmt.Printf(prefixes[NOTICE]+str, args...)
}

func (l logger) INFO(str string, args ...interface{}) {
	fmt.Printf(prefixes[INFO]+str, args...)
}

func (l logger) DEBUG(str string, args ...interface{}) {
	fmt.Printf(prefixes[DEBUG]+str, args...)
}
