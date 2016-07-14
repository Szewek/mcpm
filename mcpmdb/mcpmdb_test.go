package mcpmdb

import "testing"

func TestMCPMDatabase(t *testing.T) {
	ndb := &mcpmdbImpl{r: map[string]*MCPMPackage{}, w: map[string]*MCPMPackage{}}
	ndb.open("tmpdb")
	ndb.GetPackage("tinkers-construct")
	ndb.Close()
}
