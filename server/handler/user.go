package handler

import (
	"crypto/md5"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/game/server/common"
	"github.com/game/server/config"
	"github.com/game/server/crypto"
	"github.com/game/server/database"
	"github.com/game/server/db"
	"github.com/pkg/errors"
	"io"
	"os"
	"path/filepath"
	"strconv"
	"time"
)

type EditUserInfo struct {
	Name           string `json:"user_name"`
	Password       string `json:"password"`
	RepeatPassword string `json:"repeat_password"`
}

type RegistUser struct {
	db.User
	RepeatPassword string `json:"repeat_password"`
}

type UserResponse struct {
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
	Name      string    `json:"name"` //用户名
	Type      int       `json:"type"` //用户类型(0:超级管理员， 1: 普通用户)
}

func generateToken() string {
	crutime := time.Now().UnixNano()
	h := md5.New()
	io.WriteString(h, strconv.FormatInt(crutime, 10))
	return fmt.Sprintf("%x", h.Sum(nil))
}

//Login the root user
func Login(c *gin.Context) {
	type LoginInfo struct {
		User     string `json:"user_name"`
		Password string `json:"password"`
	}
	var info LoginInfo
	err := common.UnmarshalBody(c, &info)
	if err != nil {
		common.ErrorResponse(c, common.E_ParseBodyFailed, err.Error())
		return
	}

	var userName = info.User
	var password = info.Password
	password = crypto.PasswdEncryMD5(password)
	logger.Debug("====user: %v", info)
	var address = crypto.GetRemoteAddr(c)
	data := make(map[string]interface{})
	if dbuser, err := db.GetUserByName(userName); err != nil {
		errStr := fmt.Sprintf("user name or password error")
		common.ErrorResponse(c, common.E_UserNameOrPasswordError, errStr)
	} else if password != dbuser.Password {
		errStr := fmt.Sprintf("user name or password error")
		common.ErrorResponse(c, common.E_UserNameOrPasswordError, errStr)
	} else if !dbuser.IsValid {
		dbuser.Password = ""
		data["user"] = getUserResponse(*dbuser)
		data["token"] = ""
		err := errors.New("user name is invalid")
		common.Response(c, err, common.E_UserNameInvalid, data)
	} else {
		token := generateToken()
		if err := db.InsertToken(userName, token, address); err != nil {
			errStr := fmt.Sprintf("insert token failed!")
			common.ErrorResponse(c, common.E_UserTokenGetFailed, errStr)
		}
		dbuser.Password = ""
		data["user"] = getUserResponse(*dbuser)
		data["token"] = token
		common.Response(c, nil, common.E_Success, data)
	}
}

//add the new user to the table
func AddNewUser(c *gin.Context) {
	token := common.Header(c, "token")
	logger.Debug("====token: %s", token)
	tokenInfo, err := db.GetToken(token)
	if err != nil {
		errStr := fmt.Sprintf("get user from token: %s,err : %#v", token, err.Error())
		common.ErrorResponse(c, common.E_UserTokenGetFailed, errStr)
		return
	}
	logger.Debug("====tokenInfo: %v", tokenInfo)

	dbuser, err := db.GetUserByName(tokenInfo.Name)
	if err != nil {
		errStr := fmt.Sprintf("to get user name %s err : %#v", tokenInfo.Name, err.Error())
		common.ErrorResponse(c, common.E_UserGetFailed, errStr)
		return
	}

	if dbuser.Type != db.SUPER_ADMIN {
		errStr := fmt.Sprintf("the current user is no permission")
		common.ErrorResponse(c, common.E_UserHasNoPermission, errStr)
		return
	}

	var userInfo RegistUser
	err = common.UnmarshalBody(c, &userInfo)
	if err != nil {
		common.ErrorResponse(c, common.E_ParseBodyFailed, err.Error())
		return
	}

	userInfo.IsValid = true
	userInfo.Type = 1

	if err, errCode := checkAddInfo(userInfo); err != nil {
		common.ErrorResponse(c, errCode, err.Error())
		return
	}

	user := userInfo.User
	if err := db.CheckUser(&user); err != nil {
		common.ErrorResponse(c, common.E_UserNameIsExist, err.Error())
		return
	}


	if err := db.InsertUser(&user); err != nil {
		common.ErrorResponse(c, common.E_InsertFailed, err.Error())
		return
	}
	common.Response(c, nil, common.E_Success, nil)
}

//delete the specific user from the user's table
func DelUser(c *gin.Context) {
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

	if dbuser.Type != db.SUPER_ADMIN {
		errStr := fmt.Sprintf("the current user is no permission")
		common.ErrorResponse(c, common.E_UserHasNoPermission, errStr)
		return
	}
	delName := c.Query("name")
	if delName == "" {
		errStr := fmt.Sprintf("the username of the deleted user is empty")
		common.ErrorResponse(c, common.E_UserNameInvalid, errStr)
		return
	} else if delName == config.GetServerConfig().Config.AdminUser {
		errStr := fmt.Sprintf("the root user cannot be deleted")
		common.ErrorResponse(c, common.E_UserNameInvalid, errStr)
		return
	}

	_, err = db.GetUserByName(delName)
	if err != nil {
		errStr := fmt.Sprintf("the user name cannot find, delete failed")
		common.ErrorResponse(c, common.E_UserGetFailed, errStr)
		return
	}

	if err := db.DeleteUserByName(delName); err != nil {
		errStr := fmt.Sprintf("delete user failed")
		common.ErrorResponse(c, common.E_DeleteFailed, errStr)
		return
	}
	common.Response(c, nil, common.E_Success, nil)
}

func GetUsers(c *gin.Context) {
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

	if dbuser.Type != db.SUPER_ADMIN {
		errStr := fmt.Sprintf("the current user is no permission")
		common.ErrorResponse(c, common.E_UserHasNoPermission, errStr)
		return
	}
	users, err := db.GetUsers()
	if err != nil {
		errStr := fmt.Sprintf("get users failed")
		common.ErrorResponse(c, common.E_UserGetFailed, errStr)
		return
	}
	userResponses := make([]UserResponse, 0)
	for _, user := range users {
		userResponses = append(userResponses, getUserResponse(*user))
	}
	common.Response(c, nil, common.E_Success, userResponses)
}

func EditUser(c *gin.Context) {
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

	if dbuser.Type != db.SUPER_ADMIN {
		errStr := fmt.Sprintf("the current user is no permission")
		common.ErrorResponse(c, common.E_UserHasNoPermission, errStr)
		return
	}

	var editUserInfo EditUserInfo
	err = common.UnmarshalBody(c, &editUserInfo)
	if err != nil {
		common.ErrorResponse(c, common.E_ParseBodyFailed, err.Error())
		return
	}

	if editUserInfo.Name == "" {
		errStr := fmt.Sprintf("the username is empty")
		common.ErrorResponse(c, common.E_UserNameInvalid, errStr)
		return
	} else if editUserInfo.Password == "" || editUserInfo.RepeatPassword == "" || editUserInfo.Password != editUserInfo.RepeatPassword {
		errStr := fmt.Sprintf("the password or repeat_password is error")
		common.ErrorResponse(c, common.E_UserPasswordInvalid, errStr)
		return
	} else if editUserInfo.Name == config.GetServerConfig().Config.AdminUser {
		errStr := fmt.Sprintf("the root user cannot be changed")
		common.ErrorResponse(c, common.E_UserHasNoPermission, errStr)
		return
	}

	getUserInfo, err := db.GetUserByName(editUserInfo.Name)
	if err != nil {
		errStr := fmt.Sprintf("the user name cannot find, edit failed")
		common.ErrorResponse(c, common.E_UserGetFailed, errStr)
		return
	}

	if err := database.UpdateAttr(db.Db, getUserInfo, "password", crypto.PasswdEncryMD5(editUserInfo.Password)); err != nil {
		common.ErrorResponse(c, common.E_UpdateFailed, err.Error())
		return
	}
	common.Response(c, nil, common.E_Success, nil)
}

func checkAddInfo(info RegistUser) (error, string) {
	if info.Name == "" {
		return errors.New("name is nil"), common.E_UserNameInvalid
	} else if info.Password != info.RepeatPassword {
		return errors.New("password does not the same as repeat"), common.E_UserPasswordInvalid
	} else if info.Password == "" {
		return errors.New("password is nil"), common.E_UserPasswordInvalid
	} else if info.RepeatPassword == "" {
		return errors.New("repeat password is nil"), common.E_UserPasswordInvalid
	}
	return nil, common.E_Success
}

func saveFile(content, path, fileName string) (error, string) {
	if err := os.MkdirAll(path, os.ModePerm); err != nil {
		return fmt.Errorf("Mkdir %s failed : %s ", path, err.Error()), common.E_MkdirFailed
	}

	filePath := filepath.Join(path, fileName)
	logger.Debug("*************  filePath  ****************")
	logger.Debug(filePath)
	logger.Debug("*************  filePath  ****************")
	f, err := os.Create(filePath)
	defer f.Close()
	if err != nil {
		return fmt.Errorf("Create or open file %s failed : %s ", filePath, err.Error()), common.E_CreateFileFailed
	}
	_, err = f.WriteString(content)
	if err != nil {
		return err, common.E_WriteFileFailed
	}
	return nil, common.E_Success
}


func getUserByToken(c *gin.Context) (string, string, error) {
	token := common.Header(c, "token")
	tokenInfo, err := db.GetToken(token)
	if err != nil {
		return "", common.E_UserTokenGetFailed, err
	}

	dbuser, err := db.GetUserByName(tokenInfo.Name)
	if err != nil {
		return "", common.E_UserGetFailed, err
	}

	return dbuser.Name, common.E_Success, nil
}

func getUserResponse(user db.User) UserResponse {
	userResponse := UserResponse{}
	userResponse.Name = user.Name
	userResponse.Type = user.Type
	userResponse.UpdatedAt = user.UpdatedAt
	userResponse.CreatedAt = user.CreatedAt
	return userResponse
}
