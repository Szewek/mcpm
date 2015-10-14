package modes

import (
	"fmt"
)

func getclient(mo *ModeOptions) {
	fmt.Println("This mode should download Minecraft launcher...")
}

func init() {
	modelist["get-client"] = getclient
}
