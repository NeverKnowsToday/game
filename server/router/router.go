package router

import (
	"github.com/gin-gonic/gin"
	"github.com/game/server/handler"
	logHandler "github.com/game/server/logger"
	"sync"
)

// Router 全局路由
var router *gin.Engine
var onceCreateRouter sync.Once

func GetRouter() *gin.Engine {
	onceCreateRouter.Do(func() {
		router = createRouter()
	})

	return router
}

func createRouter() *gin.Engine {
	router := gin.Default()

	router.Use(handler.Cors())

	v1 := router.Group("/v1")
	{
		v1.GET("/log", logHandler.GetServiceLog)
		v1.POST("/log", logHandler.SetServiceLog)
		user := v1.Group("/user")
		{
			//用户
			user.POST("/login", handler.Login)
			user.Use(handler.Filter())
			user.POST("/adduser", handler.AddNewUser)
			user.DELETE("/", handler.DelUser)
			user.PUT("/edituser", handler.EditUser)
			user.GET("/getusers", handler.GetUsers)
		}

		game := v1.Group("/api")
		{
			// 操作游戏
			//game.Use(handler.Filter())
			//game.POST("/start", handler.Start)
			game.GET("/invoke", handler.Invoke)
			game.POST("/insurance/benefit_diff", handler.InsuranceCompare)
			game.GET("/insurance/excel/list", handler.GetExcelList)

		}

	}
	return router
}
