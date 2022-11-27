package logger

import (
	"github.com/gin-gonic/gin"
	"github.com/game/server/common"
	"github.com/game/server/logger/logging"
)

func GetServiceLog(c *gin.Context) {
	level := logging.GetLoggerLevel(c.Query("module"))
	common.Response(c, nil, common.E_Success, level)
}

func SetServiceLog(c *gin.Context) {
	var myLogger logging.Logger
	if common.CheckError(c, common.UnmarshalBody(c, &myLogger), common.E_JsonUnmarshalError) {
		return
	}

	logging.SetLevel(myLogger.Module, myLogger.Level)
	common.Response(c, nil, common.E_Success, "成功")
}
