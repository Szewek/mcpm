package modes

import (
	"fmt"
)

func getserver(mo *ModeOptions) {
	fmt.Println("This mode should download Minecraft Server...")
}

func init() {
	modelist["get-server"] = getserver
}
