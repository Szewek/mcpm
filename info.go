package main

import (
	"encoding/json"
	"fmt"
)

type (
	_DownloadInfo struct {
		Created_at, Name, Type string
		Id                     int
	}
	_PackageInfo struct {
		Authors  []string
		Title    string
		Versions map[string][]_DownloadInfo
	}
)

func readPackageInfo() {
	pkgn := flagset.Arg(0)
	if pid, ok := (*_PKGS)[pkgn]; ok {
		data := (*_DBASE)[pid]
		pr, hte := downloadPackageInfo(&data)
		must(hte)
		defer pr.Close()
		var js _PackageInfo
		jd := json.NewDecoder(pr)
		must(jd.Decode(&js))

		fmt.Printf("Name: %s\nAuthors: %s\nVersions:\n", js.Title, js.Authors)
		for k, v := range js.Versions {
			fmt.Printf("  %s:\n", k)
			for i := 0; i < len(v); i++ {
				fmt.Printf("    %s; %s (%d) (%s)\n", v[i].Name, v[i].Type, v[i].Id, v[i].Created_at)
			}
		}
	}
}
