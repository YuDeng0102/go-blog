package main

import (
	"server/core"
	"server/global"
)

func main() {
	global.Config = core.InitConf()
	global.Log = core.InitLogger()

	core.RunServer()
}
