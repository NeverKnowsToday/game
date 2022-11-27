package crypto

import (
	"crypto/md5"
	"encoding/hex"
	"fmt"
	"github.com/gin-gonic/gin"
	"strings"
	"crypto/sha256"
)

func Md5sum(data []byte) string {
	w := md5.New()
	w.Write(data)
	hexData := w.Sum(nil)
	strData := hex.EncodeToString(hexData)
	return strData
}

// User password MD5 encryption
func PasswdEncryMD5(str string) string {
	data := []byte(str)
	has := md5.Sum(data)
	md5str := fmt.Sprintf("%x", has)
	return md5str
}

func GetRemoteAddr(c *gin.Context) string {
	address := c.Request.RemoteAddr
	index := strings.Index(address, ":")
	return address[:index]
}

func Sha256(data []byte) string {
	w := sha256.New()
	w.Write(data)
	hexData := w.Sum(nil)
	strData := hex.EncodeToString(hexData)
	return strData
}