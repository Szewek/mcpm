package modes

import (
	"encoding/json"
	"fmt"

	"github.com/Szewek/mcpm/database"
	"github.com/Szewek/mcpm/util"
)

type (
	downloadInfo struct {
		CreatedAt  string `json:"Created_at"`
		Name, Type string
		Id         int
	}
	packageInfo struct {
		Authors  []string
		Title    string
		Versions map[string][]downloadInfo
	}
)

func info(mo *ModeOptions) {
	pn := mo.Args[0]
	db := database.GetDatabase()
	if pkg := db.Packages().GetByName(pn); pkg != nil {
		pr, hte := util.DownloadPackageInfo(pkg.Type, pkg.ID, pkg.Name)
		util.Must(hte)
		defer pr.Close()

		var pi packageInfo
		jd := json.NewDecoder(pr)
		util.Must(jd.Decode(&pi))

		fmt.Printf("Name: %s\nAuthors: %s\nVersions:\n", pi.Title, pi.Authors)
		for k, v := range pi.Versions {
			fmt.Printf("  %s:\n", k)
			for i := 0; i < len(v); i++ {
				fmt.Printf("    %s; %s (%d) (%s)\n", v[i].Name, v[i].Type, v[i].Id, v[i].CreatedAt)
			}
		}
	}
}

func init() {
	registerMode("info", info)
}
