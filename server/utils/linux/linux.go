package linux

import (
	"bytes"
	"crypto/md5"
	"errors"
	"fmt"
	"io"
	"net"
	"os"
	"os/exec"
	"strings"
	"time"
)

// 内存信息
type Memery struct {
	Total     string //系统总的可用物理内存大小
	Used      string //已被使用的物理内存的大小
	Free      string //还有多少物理内存可用
	Shared    string //被共享使用的物理内存大小
	BuffCache string //被buffer和cache使用的物理内存大小
	Available string //还可以被应用程序使用的物理内存大小
}

// GetLocalIP 获取本地IP地址
func GetLocalIP() (string, error) {
	addrs, err := net.InterfaceAddrs()
	if err != nil {
		return "", err
	}

	for _, addr := range addrs {
		if ipnet, ok := addr.(*net.IPNet); ok && !ipnet.IP.IsLoopback() {
			if nil != ipnet.IP.To4() {
				return ipnet.IP.String(), nil
			}
		}
	}

	return "", fmt.Errorf("no valid host IP")
}

// GetLocalUser 获取本地用户
func GetLocalUser() (string, error) {

	res, _, err := ExecWithShellTimeout("whoami", 10)
	if err != nil {
		return "", err
	}
	return res, nil
}

// ExecWithShellTimeout ...
func ExecWithShellTimeout(cmd string, timeout time.Duration) (string, string, error) {
	var stdout bytes.Buffer
	var stderr bytes.Buffer

	runningCmd := exec.Command("sh", "-c", cmd)
	//runningCmd.SysProcAttr = &syscall.SysProcAttr{Setpgid: true}
	runningCmd.Stdout = &stdout
	runningCmd.Stderr = &stderr
	err := runningCmd.Start()
	if err != nil {
		return "", "", err
	}
	execErrChan := make(chan error, 1)
	//defer close(execErrChan)
	go func() {
		execErrChan <- runningCmd.Wait()
	}()

	select {
	case err = <-execErrChan:
		stdoutStr := strings.TrimSpace(string(bytes.TrimRight(stdout.Bytes(), "\x00")))
		stderrStr := strings.TrimSpace(string(bytes.TrimRight(stderr.Bytes(), "\x00")))
		if err != nil {
			//close(execErrChan)
			return stdoutStr, stderrStr, err
		}
		//close(execErrChan)
		return stdoutStr, stderrStr, nil
	case <-time.After(time.Second * timeout):
		//close(execErrChan)
		var errMsg string
		errMsg = "timeout"
		// if err := syscall.Kill(-runningCmd.Process.Pid, syscall.SIGKILL); err != nil {
		// 	errMsg = fmt.Sprintf("wait [%s] timeout [%ds], kill pid [%d] error: [%s]", cmd, timeout, runningCmd.Process.Pid, err)
		// } else {
		// 	errMsg = fmt.Sprintf("wait [%s] timeout [%ds], kill pid [%d]", cmd, timeout, runningCmd.Process.Pid)
		// }
		//close(execErrChan)
		return "", "", errors.New(errMsg)
	}
}

// CheckStringHaveWE ...
func CheckStringHaveWE(res string) error {
	errKeyWord := []string{"ERROR", "[WARNING]"}
	for _, keyWord := range errKeyWord {
		if err := strings.Contains(res, keyWord); err {
			return errors.New(res)
		}
	}

	return nil
}

// PortTelnetCheck 检查端口存活性
func PortTelnetCheck(IP, port string) error {
	address := net.JoinHostPort(IP, port)
	conn, err := net.DialTimeout("tcp", address, 1*time.Second)
	if err != nil {
		return err
	}
	if conn != nil {
		_ = conn.Close()
	}

	return nil
}

// PathExists ...
func PathExists(path string) (bool, error) {
	_, err := os.Stat(path)
	if err == nil {
		return true, nil
	}
	if os.IsNotExist(err) {
		return false, nil
	}
	return false, err
}

// DeleteFile ...
func DeleteFile(filePath string) (string, error) {
	is, err := PathExists(filePath)
	if err != nil {
		return "", err
	}
	if !is {
		return fmt.Sprintf("not found %s, skip do delete", filePath), nil
	}
	if err := os.Remove(filePath); err != nil {
		return "", err
	}
	return fmt.Sprintf("remove %s", filePath), nil
}

// CreateDir ...
func CreateDir(path string) error {
	f, err := os.Stat(path)
	if err != nil {
		err = os.MkdirAll(path, os.ModePerm)
		if err != nil {
			return err
		}
		return nil
	}
	if f.IsDir() {
		return nil
	}
	return nil
}

// FormatFileSize ...
func FormatFileSize(fileSize int) (size string) {
	if fileSize < 1024 {
		//return strconv.FormatInt(fileSize, 10) + "B"
		return fmt.Sprintf("%.2fB", float64(fileSize)/float64(1))
	} else if fileSize < (1024 * 1024) {
		return fmt.Sprintf("%.2fKB", float64(fileSize)/float64(1024))
	} else if fileSize < (1024 * 1024 * 1024) {
		return fmt.Sprintf("%.2fMB", float64(fileSize)/float64(1024*1024))
	} else if fileSize < (1024 * 1024 * 1024 * 1024) {
		return fmt.Sprintf("%.2fGB", float64(fileSize)/float64(1024*1024*1024))
	} else if fileSize < (1024 * 1024 * 1024 * 1024 * 1024) {
		return fmt.Sprintf("%.2fTB", float64(fileSize)/float64(1024*1024*1024*1024))
	} else { //if fileSize < (1024 * 1024 * 1024 * 1024 * 1024 * 1024)
		return fmt.Sprintf("%.2fEB", float64(fileSize)/float64(1024*1024*1024*1024*1024))
	}
}

// GetCurrentTimeStr ...
func GetCurrentTimeStr() string {
	now := time.Now()
	currentTimeStr := now.Format("2006-01-02 15:04:05")

	return currentTimeStr
}

// GetCurrentTimeStr ...
func GetCurrentTimeStr2() string {
	now := time.Now()
	currentTimeStr := now.Format("20060102_150405")

	return currentTimeStr
}

// StrToMd5Str ...
func StrToMd5Str(str string) string {
	has := md5.Sum([]byte(str))

	return fmt.Sprintf("%x", has)
}

// UpperPath 返回上一级目录，/app/opengauss/brmsoft 返回 /app/opengauss/
func UpperPath(path string) string {
	aa := strings.Split(path, "/")

	res := ""
	for c, v := range aa {
		if c+1 != len(aa) {
			res += v + "/"
		}
	}

	return res
}

// LastLevelPath 返回最后一级目录，/app/opengauss/brmsoft 返回 brmsoft
func LastLevelPath(path string) string {
	aa := strings.Split(path, "/")

	res := ""
	for c, v := range aa {
		if c+1 == len(aa) {
			res = v
		}
	}

	return res
}

type tempFile struct {
	Path string //临时文件路径
}

func NewTempFile(path string) *tempFile {
	return &tempFile{
		Path: path,
	}

}

func (t *tempFile) CreateTempFile() (string, error) {
	f, err := os.Create(t.Path)
	if err != nil {
		return "", fmt.Errorf("CreateTmpFile [%s] failed, err: [%s]", t.Path, err)
	}
	defer f.Close()
	return t.Path, nil
}

func (t *tempFile) AddSlice2File(data []string) (n int, err error) {

	var newConf string

	for _, v := range data {
		newConf += v + "\n"
	}

	f, err := os.OpenFile(t.Path, os.O_RDWR, 0755)
	if err != nil {
		return 0, fmt.Errorf("AddParam2File [%s] failed, err: [%s]", t.Path, err)
	}

	defer f.Close()
	return f.WriteString(newConf)
}

// ReadFile reads the named file and returns the contents.
// A successful call returns err == nil, not err == EOF.
// Because ReadFile reads the whole file, it does not treat an EOF from Read
// as an error to be reported.
func ReadFile(name string) ([]byte, error) {
	f, err := os.OpenFile(name, os.O_RDONLY, 0)
	if err != nil {
		return nil, err
	}
	defer f.Close()

	var size int
	if info, err := f.Stat(); err == nil {
		size64 := info.Size()
		if int64(int(size64)) == size64 {
			size = int(size64)
		}
	}
	size++ // one byte for final read at EOF

	// If a file claims a small size, read at least 512 bytes.
	// In particular, files in Linux's /proc claim size 0 but
	// then do not work right if read in small pieces,
	// so an initial read of 1 byte would not work correctly.
	if size < 512 {
		size = 512
	}

	data := make([]byte, 0, size)
	for {
		if len(data) >= cap(data) {
			d := append(data[:cap(data)], 0)
			data = d[:len(data)]
		}
		n, err := f.Read(data[len(data):cap(data)])
		data = data[:len(data)+n]
		if err != nil {
			if err == io.EOF {
				err = nil
			}
			return data, err
		}
	}
}
