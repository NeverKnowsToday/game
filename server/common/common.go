package common

import (
	"crypto/md5"
	"encoding/json"
	"fmt"
	"github.com/gin-gonic/gin"
	//"github.com/game/server/config"
	"github.com/game/server/logger/logging"
	"github.com/pkg/errors"
	"io/ioutil"
	"net/http"
	"net/url"
	"strings"
	"sync"
)

var logger = logging.GetLogger("common", logging.DEFAULT_LEVEL)
var EXPIRATION_TIME int
var RWMutex sync.RWMutex

const (
	USER_CERT       = "user.pem"
	USER_KEY        = "user.key"
	CertPath        = "./cert"
	WordLibraryPath = "./wordLibrary.txt"
)

var ErrCueMaps = map[string]string{}

func InitErrorCue(configFile string) error {
	bytes, err := ioutil.ReadFile(configFile)
	if err != nil {
		return err
	}
	err = json.Unmarshal(bytes, &ErrCueMaps)
	if err != nil {
		return err
	}
	return nil
}

func ReadBody(c *gin.Context) ([]byte, error) {
	body, err := ioutil.ReadAll(c.Request.Body)
	logger.Debug("***********body info******************************")
	logger.Debug(string(body))
	logger.Debug("*****************************************")
	return body, err
}

func UnmarshalBody(c *gin.Context, info interface{}) error {
	body, err := ReadBody(c)
	if err != nil {
		return err
	}
	if err = json.Unmarshal(body, info); err != nil {
		logger.Error(err.Error())
	} else {
		logger.Debug("***********struct info******************************")
		logger.Debugf("%v", info)
		logger.Debug("*****************************************")
	}
	return err
}

func Query(c *gin.Context, key string) string {
	val := c.Query(key)
	logger.Debugf(">>>>Query:{%s: %s}", key, val)
	return val
}

func Param(c *gin.Context, key string) string {
	val := c.Param(key)
	logger.Debugf(">>>>Param:{%s: %s}", key, val)
	return val
}

func Header(c *gin.Context, key string) string {
	val := c.Request.Header.Get(key)
	if key == "user" && strings.Contains(val, "%") {
		val, _ = url.QueryUnescape(val)
	}
	logger.Debugf(">>>>Header:{%s,%s}", key, val)
	return val
}

type ResponseInfo struct {
	ErrCode string      `json:"err_code"` // 错误码
	ErrCue  string      `json:"err_cue"`  // 错误提示语
	ErrMsg  string      `json:"err_msg"`  // 错误原因
	Data    interface{} `json:"data"`
}

// Response http response
func Response(c *gin.Context, err error, errCode string, data interface{}) {
	var resp ResponseInfo
	if (err != nil && errCode == E_UserTokenInvalid) || (err != nil && errCode == E_UserNameInvalid) {
		logger.Errorf("*********************\n%+v*********************\n", err)
		c.Writer.WriteHeader(http.StatusUnauthorized)
		resp.ErrCode = errCode
		resp.ErrCue = ErrCueMaps[errCode]
		resp.ErrMsg = err.Error()
	} else if err != nil {
		logger.Errorf("*********************\n%+v*********************\n", err)
		c.Writer.WriteHeader(http.StatusBadRequest)
		resp.ErrCode = errCode
		resp.ErrCue = ErrCueMaps[errCode]
		resp.ErrMsg = err.Error()
	} else {
		c.Writer.WriteHeader(http.StatusOK)
		resp.ErrCode = E_Success
		resp.ErrCue = ErrCueMaps[E_Success]
	}

	resp.Data = data
	ret, _ := json.Marshal(resp)
	c.Writer.Write(ret)
	c.Writer.Header().Set("Access-Control-Allow-Origin", "*")
	c.Writer.Header().Set("Content-Type","application/json")
}

//response for details
func CheckErrorDetail(c *gin.Context, err error, errCode string, format string, a ...interface{}) bool {
	if err != nil {
		err = errors.Wrapf(err, format, a...)
	}
	return CheckError(c, err, errCode)
}

//check error,if has error return true
func CheckError(c *gin.Context, err error, errCode string) bool {
	if err != nil {
		Response(c, err, errCode, nil)
		return true
	}
	return false
}

func ErrorResponse(c *gin.Context, errCode string, format string, a ...interface{}) {
	Response(c, errors.Errorf(format, a...), errCode, nil)
}

func GetRemoteAddr(c *gin.Context) string {
	address := c.Request.RemoteAddr
	index := strings.Index(address, ":")
	return address[:index]
}

// User password MD5 encryption
func PasswdEncryMD5(str string) string {
	data := []byte(str)
	has := md5.Sum(data)
	md5str := fmt.Sprintf("%x", has)
	return md5str
}
