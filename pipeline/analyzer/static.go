package main

import (
	"encoding/json"
	"encoding/xml"
	"github.com/sociam/xray-archiver/pipeline/util"
	"io/ioutil"
	"os"
	"os/exec"
	"path"
	"path/filepath"
	"strings"
)

// AndroidManifest is a struct representing the interesting parts of the
// AndroidManifest.xml in APKs
type AndroidManifest struct {
	Package    string            `xml:"package,attr"`
	Perms      []util.Permission `xml:"uses-permission"`
	Sdk23Perms []util.Permission `xml:"uses-permission-sdk-23"`
}

func parseManifest(app *util.App) (*AndroidManifest, error) {
	manifest := AndroidManifest{}
	manifestFile, err := os.Open(path.Join(app.OutDir(), "AndroidManifest.xml"))
	if err != nil {
		return nil, err
	}
	bytes, err := ioutil.ReadAll(manifestFile)
	if err != nil {
		return nil, err
	}
	err = xml.Unmarshal(bytes, &manifest)
	if err != nil {
		return nil, err
	}

	if manifest.Package != "" {
		app.ID = manifest.Package
	}
	return &manifest, nil
}

func (manifest *AndroidManifest) getPerms() []util.Permission {
	return append(manifest.Perms, manifest.Sdk23Perms...)
}

type company struct {
	ID           string   `json:"id"`
	Name         string   `json:"company"`
	Domains      []string `json:"domains"`
	Founded      string   `json:"founded"`
	Acquired     string   `json:"acquired in"`
	CType        string   `json:"type"`
	TypeTag      string   `json:"typetag"`
	Jurisdiction string   `json:"jurisdiction_code"`
	Parent       string   `json:"parent"`
	Capital      string   `json:"capital"`
	Equity       string   `json:"equity"`
	Size         string   `json:"size"`
	DataSource   string   `json:"data source"`
	Description  string   `json:"description"`
}

func simpleAnalyze(app *util.App) ([]string, error) {
	//TODO: fix error handling

	//TODO: replace with DB calls
	var companies map[string]company
	companyFile, err := os.Open(path.Join(util.Cfg.DataDir, "company_details.json"))
	if err != nil {
		return []string{}, err
	}
	bytes, err := ioutil.ReadAll(companyFile)
	if err != nil {
		return []string{}, err
	}
	err = json.Unmarshal(bytes, &companies)
	if err != nil {
		return []string{}, err
	}

	for name := range companies {
		if _, ok := trackers[name]; !ok {
			delete(companies, name)
		}
	}

	// getDomainCo := func(host string) *string {
	// 	for _, company := range companies {
	// 		for _, domain := range company.domains {
	// 			if strings.Contains(host, domain) {
	// 				return &company.id
	// 			}
	// 		}
	// 	}
	// 	return nil
	// }

	cmd := exec.Command("grep", "-Erho", "\"https?://[^ >]+\"",
		path.Join(app.OutDir(), "smali", "**", "*.smali"))
	out, err := cmd.Output()
	if err != nil {
		return []string{}, err
	}
	urls := strings.Split(string(out), "\n")

	// var appTrackers []string

	// irrelevant := []string{"app", "identity", "n/a", "other", "", "library"}
	// for name, company := range companies {
	// 	for _, domain := range company.Domains {
	// 		if strings.Contains(string(urls), string(domain)) {
	// 			toAppend := true
	// 			for _, cat := range irrelevant {
	// 				if companies[name].TypeTag == cat {
	// 					toAppend = false
	// 					break
	// 				}
	// 			}
	// 			if toAppend {
	// 				appTrackers = append(appTrackers, name)
	// 			}
	// 		}
	// 	}
	// }

	return urls, nil
}

func findPackages(app *util.App) ([]string, error) {
	// TODO: fix error handling
	paths := make(map[string]util.Unit)
	err := os.Chdir(path.Join(app.OutDir(), "smali"))
	if err != nil {
		return []string{}, err
	}

	err = filepath.Walk(".",
		func(fname string, info os.FileInfo, err error) error {
			if err != nil {
				return err
			}

			if path.Ext(fname) == ".smali" {
				paths[path.Dir(fname)] = unit
			}

			return nil
		})

	pkgs := make([]string, 0, 20)
	for path := range paths {
		//pkg := strings.Replace(path, string(os.PathSeparator), ".", -1)
		pkg := strings.Map(func(ch rune) rune {
			if ch == os.PathSeparator {
				return '.'
			}
			return ch
		}, path)
		pkgs = append(pkgs, pkg)
	}

	return pkgs, err
}
