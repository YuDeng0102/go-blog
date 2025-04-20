package main

import (
	"server/core"
	"server/flag"
	"server/global"
	"server/initialize"
	"server/rabbitmq"
)

func main() {
	global.Config = core.InitConf()
	global.Log = core.InitLogger()
	initialize.OtherInit()
	global.DB = initialize.InitGorm()
	global.Redis = initialize.ConnectRedis()
	global.ESClient = initialize.ConnectEs()
	defer global.Redis.Close()
	global.RmqConn = initialize.RabbitmqInit()
	defer global.RmqConn.Close()

	flag.InitFlag()

	initialize.InitCron()

	go rabbitmq.StartESUpdateConsumer()

	core.RunServer()
}
