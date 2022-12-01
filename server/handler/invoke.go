package handler

import (
	"github.com/game/server/common"
	"github.com/game/server/database"
	"github.com/game/server/db"
	"github.com/gin-gonic/gin"
	"math/rand"
	"strings"
	"time"
)

func Invoke(c *gin.Context) {
	token := common.Header(c, "token")
	tokenInfo, err := db.GetToken(token)
	if err != nil {
		common.ErrorResponse(c, common.E_UserTokenGetFailed, err.Error())
		return
	}

	 // 通过token查询当前登陆的用户信息
	dbuser, err := db.GetUserByName(tokenInfo.Name)
	if err != nil {
		common.ErrorResponse(c, common.E_UserGetFailed, err.Error())
		return
	}

	//  生成1-6的随机数
	rand.Seed(time.Now().UnixNano())
	rnum := rand.Int()%6+1

	invoke := &db.Invoke{
		Name : dbuser.Name,
		Room : 1,       // 房间号为1固定
	}

    // 查询用户操作表是否已经有该用户记录
	userInvoke, err := db.GetInvokeByName(dbuser.Name)
	if err != nil && !strings.Contains(err.Error(), "record not found") {
		common.ErrorResponse(c, common.E_UserGetFailed, err.Error())
		return
	}

	// 默认设置起始位置为0
    userInvokeOriginalPos := 0

    // 如果表里已经有记录，就把起始位置设置成库里的记录
    if userInvoke != nil {
		userInvokeOriginalPos = userInvoke.CurrentPos
	}

	// （没有该用户数据，说明该用户没有投过骰子， 就插入该用入记录到invoke表中）
	if database.RecordNotFound(err) {
		db.InsertInvoke(invoke)
	}

	// 如果有错误，而且错误不是"record not found, 就返回错误信息"
	if err != nil &&!database.RecordNotFound(err) {
		common.ErrorResponse(c, common.E_UserGetFailed, err.Error())
		return
	}

	// 通过房间号，查询出在这个房间的两个用户
	users, err  := db.GetInvokesByRoom(userInvoke.Room)
	if err != nil {
		common.ErrorResponse(c, common.E_UserGetFailed, err.Error())
		return
	}

	// 设置当前位置 = 起始位置+随机数
	invoke.CurrentPos = userInvoke.CurrentPos + rnum
	if invoke.CurrentPos > 100 { // 如果当前位置>100 , 就产生回退
		invoke.CurrentPos =  100 - (invoke.CurrentPos % 100)
	}
	if invoke.CurrentPos == 100 {// 如果当前位置== 100，获得胜利
		invoke.IsWin = true
		invoke.CurrentPos = 100
	}

	var CurUser, NoCurUser *db.Invoke
	if len(users) == 1 { // 如果invoke表里只有1号房间的一条记录，说明只有一个用户投过骰子，把当前位置直接存入数据库
		err := db.UpdateInvokePosByname(users[0].Name, invoke.CurrentPos)
		if err != nil {
			common.ErrorResponse(c, common.E_UpdateFailed, err.Error())
			return
		}
	}else if len(users) == 2 {// 如果有两条记录，需要判断是否发生蛇吻
		for _, v := range users {
			if v.Name == dbuser.Name {
				CurUser = v
			}else{
				NoCurUser = v
			}
		}
		if CurUser == nil || NoCurUser == nil {
			common.ErrorResponse(c, common.E_UserGetFailed, "")
            return
		}

		if CurUser.CurrentPos == NoCurUser.CurrentPos { // 发生蛇吻，回退到原来位置
			invoke.CurrentPos = userInvokeOriginalPos
		}
		err := db.UpdateInvokePosByname(CurUser.Name, invoke.CurrentPos) // 更新数据库位置
		if err != nil {
			common.ErrorResponse(c, common.E_UpdateFailed, err.Error())
			return
		}
		err = db.UpdateInvokeWinnerByname(CurUser.Name, invoke.IsWin) // 更新数据库获胜者
		if err != nil {
			common.ErrorResponse(c, common.E_UpdateFailed, err.Error())
			return
		}

	}else {
		common.ErrorResponse(c, common.E_UserGetFailed, "")
		return
	}

	common.Response(c, nil, common.E_Success, invoke)
}