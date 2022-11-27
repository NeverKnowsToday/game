package handler

import (
	"github.com/game/server/common"
	"github.com/game/server/db"
	"github.com/gin-gonic/gin"
)

func Invoke(c *gin.Context) {
	token := common.Header(c, "token")
	tokenInfo, err := db.GetToken(token)
	if err != nil {
		common.ErrorResponse(c, common.E_UserTokenGetFailed, err.Error())
		return
	}

	dbuser, err := db.GetUserByName(tokenInfo.Name)
	if err != nil {
		common.ErrorResponse(c, common.E_UserGetFailed, err.Error())
		return
	}

	database.CheckInvoke()








	common.Response(c, nil, common.E_Success, nil)

}