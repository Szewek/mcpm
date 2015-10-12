package modes

import "github.com/Szewek/mcpm/database"

func update(mo *ModeOptions) {
	database.UpdateDatabase(mo.Verbose)
}

func init() {
	modelist["update"] = update
}
