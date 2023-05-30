package main

import (
	"flag"
	"fmt"
	"github.com/game/server/config"
	"github.com/game/server/db"
	"github.com/game/server/logger/logging"
	"github.com/game/server/router"
	"runtime"
)

var (
	configPath = flag.String("configPath", "./config.yaml", "config path")
	logger     = logging.GetLogger("main", logging.DEFAULT_LEVEL)
)

func checkArg(val interface{}, name string) {
	switch val.(type) {
	case *string:
		if val == nil || *val.(*string) == "" {
			logger.Panicf("The flag %s is empty", name)
		}
	case *int:
		if val == nil || *val.(*int) == 0 {
			logger.Panicf("The flag %s is empty", name)
		}
	default:
		logger.Debug("Unknown argument type")
	}
}

func main(){
	// parse init param
	logger.Debug("Usage : ./server -configPath=")
	flag.Parse()
	checkArg(configPath, "configPath")
	err := config.InitConfig(*configPath)
	if err != nil {
		panic(err)
	}
	//err = common.InitErrorCue("./ErrorCue.json")
	//if err != nil {
	//	panic(err)
	//}
	logger.Debugf("config parser :===============test")


	//初始化数据库
	if err = db.InitDb(); err != nil {
		panic(err)
	}



	// 设置使用系统最大CPU
	runtime.GOMAXPROCS(runtime.NumCPU())
	// 构造路由器
	r := router.GetRouter()

	// 运行服务
	//address := fmt.Sprintf("0.0.0.0:%d", config.GetServerConfig().Port)

	address := fmt.Sprintf("127.0.0.1:%d", 8080)
	r.Run(address)


}