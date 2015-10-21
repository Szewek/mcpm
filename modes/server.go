package modes

import (
	"fmt"
)

const srvurl = "http://s3.amazonaws.com/Minecraft.Download/versions/%s/minecraft_server.%s.jar"

func getserver(mo *ModeOptions) {
	fmt.Println("This mode should download Minecraft Server...")
}

func init() {
	registerMode("get-server", getserver)
}
