package handler

import (
	"github.com/gin-gonic/gin"
	"github.com/game/server/common"
	"github.com/game/server/config"
	"github.com/game/server/database"
	"github.com/game/server/db"
	"strconv"
	"time"
)


func Filter() gin.HandlerFunc {
	return func(c *gin.Context) {
		logger.Debug("filter")
		//text, err := ioutil.ReadAll(c.Request.Body)
		//敏感词处理todo
		//c.Request.Body = ioutil.NopCloser(bytes.NewBuffer(text))
		token := c.Request.Header.Get("token")
		userinfo := db.User{}
		if token == "" {
			common.ErrorResponse(c, common.E_UserTokenInvalid, "invalid token!")
			c.Abort()
			return
		}
		tokeninfo := &db.Token{}
		errToken := database.Model(db.Db, &db.Token{}).Where("token = ?", token).Find(&tokeninfo)
		if errToken != nil || tokeninfo == nil {
			common.ErrorResponse(c, common.E_UserTokenInvalid, "invalid login status!")
			c.Abort()
			return
		}
		logger.Debug(tokeninfo)
		errUser := database.Model(db.Db, &db.User{}).Where("name = ?", tokeninfo.Name).Find(&userinfo)
		if errUser != nil {
			common.ErrorResponse(c, common.E_UserStateInvalid, "invalid user!")
			c.Abort()
			return
		}

		// 添加 header 权限验证
		defaultConfig := config.GetServerConfig().Default
		timestamp := common.Header(c, "timestamp")
		timeSign := common.Header(c, "timesign")
		if timestamp == "" || timeSign == "" {
			common.ErrorResponse(c, common.E_InvalidOperation, "请校验 timestamp")
			c.Abort()
			return
		} else if sign := common.PasswdEncryMD5(defaultConfig.MachineSecret + timestamp); timeSign != sign {
			logger.Errorf("校验 timestamp 失败; BCAP_MACHINE_SECRET = %s; sign = %s; 传入的timesign = %s", defaultConfig.MachineSecret, sign, timeSign)
			common.ErrorResponse(c, common.E_InvalidOperation, "校验 timestamp 失败")
			c.Abort()
			return
		}

		//设置超时时间
		expiration := config.GetServerConfig().Default.ExpirationTime
		if expiration != "" {
			common.EXPIRATION_TIME, _ = strconv.Atoi(expiration)
		} else {
			common.EXPIRATION_TIME = 1800
		}

		var cstZone = time.FixedZone("CST", 8*3600)
		if time.Now().In(cstZone).Unix()-tokeninfo.UpdatedAt.In(cstZone).Unix() > int64(common.EXPIRATION_TIME) {
			db.DeleteToken(token)
			common.ErrorResponse(c, common.E_UserTokenInvalid, "token expired")
			c.Abort()
			return
		}
		tokeninfo.UpdatedAt = time.Now()

		database.UpdateAttr(db.Db, tokeninfo, "updated_at", tokeninfo.UpdatedAt)

		c.Next()
	}
}
