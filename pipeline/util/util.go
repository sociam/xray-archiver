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
)

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

// ApkPath creates a string that represents the location of the APK
// on disk. Used to populate the Path string of an App object.
func (app *App) ApkPath() string {
	if app.Path != "" {
		return app.Path
	}
	return path.Join(
		Cfg.AppDir, app.ID, app.Store, app.Region,
		app.Ver, app.ID+".apk")
}

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

func (app *App) Unpack() error {
	apkPath, outDir := app.ApkPath(), app.OutDir()
	if _, err := os.Stat(apkPath); err != nil {
		if os.IsNotExist(err) {
			return err
		}
		return fmt.Errorf("Failed to open apk %s: %s", apkPath, err.Error())
	}

	cmd := exec.Command("apktool", "d", apkPath, "-o", outDir, "-f")
	out, err := cmd.CombinedOutput()
	if err != nil {
		return fmt.Errorf("Error '%s' unpacking apk; output below:\n%s",
			err.Error(), string(out))
	}
	return nil
}

func (app *App) Cleanup() error {
	return os.RemoveAll(app.OutDir())
}

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

func Combine(a, b map[string]Unit) map[string]Unit {
	ret := a
	for e := range b {
		ret[e] = unit
	}
	return ret
}

func StrMap(args ...string) map[string]Unit {
	ret := make(map[string]Unit)
	for _, e := range args {
		ret[e] = unit
	}

	return ret
}

func WriteJSON(w io.Writer, data interface{}) error {
	enc := json.NewEncoder(w)
	enc.SetEscapeHTML(false)
	enc.SetIndent("", "")
	return enc.Encode(data)
}

func WriteDEAN(w io.Writer, data interface{}) error {
	w.Write([]byte("Nah\n"))
	WriteJSON(w, data)
	w.Write([]byte("mate."))
	return nil
}
