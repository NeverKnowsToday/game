package linux
//
//import (
//	"context"
//	"errors"
//	"fmt"
//	"strings"
//	"time"
//
//	"icode.baidu.com/baidu/paas-xcloud/dbpaas-api/pkg/common/constant"
//)
//
////BuildOMMSSHTrust 构建omm账户的origin <-> target SSH互信(isDuplex == true 表示双端互信，false表示origin -> target端的互信)
//func BuildOMMSSHTrust(ctx context.Context, isDuplex bool, originIP, targetIP string) error {
//	user := "omm"
//	return BuildSSHTrust(ctx, isDuplex, user, user, originIP, targetIP)
//}
//
////BuildSSHTrust 构建origin <-> target SSH互信(isDuplex == true 表示双端互信，false表示origin -> target端的互信)
//func BuildSSHTrust(ctx context.Context, isDuplex bool, originUser, targetUser, originIP, targetIP string) error {
//	originCli := fmt.Sprintf(`ssh %s@%s `, originUser, originIP)
//	targetCli := fmt.Sprintf(`ssh %s@%s `, targetUser, targetIP)
//	timeout := time.Duration(2)
//
//	//check
//	originCheckCmd := fmt.Sprintf(`%s "%s \"ls\""`, originCli, targetCli)
//	_, _, err := ExecWithShellTimeout(originCheckCmd, timeout)
//	if err == nil {
//		return nil
//	}
//
//	noStrictKeyChecking := "echo 'StrictHostKeyChecking=no' >> ~/.ssh/config"
//	getPub := "cat ~/.ssh/id_rsa.pub"
//
//	getOriginPubCmd := fmt.Sprintf(`%s "%s; %s"`, originCli, noStrictKeyChecking, getPub)
//	orginPub, _, err := ExecWithShellTimeout(getOriginPubCmd, timeout)
//	if err != nil {
//		return fmt.Errorf("cmd: %s, error: %s", getOriginPubCmd, err.Error())
//	}
//
//	if isDuplex {
//		buildOriginSSHCmd := fmt.Sprintf(`%s "echo '%s' >> ~/.ssh/authorized_keys; %s; %s"`, targetCli, orginPub, noStrictKeyChecking, getPub)
//		targetPub, _, err := ExecWithShellTimeout(buildOriginSSHCmd, timeout)
//		if err != nil {
//			return fmt.Errorf("cmd: %s, error: %s", buildOriginSSHCmd, err.Error())
//		}
//		buildTargetSSHCmd := fmt.Sprintf(`%s "echo '%s' >> ~/.ssh/authorized_keys"`, originCli, targetPub)
//		_, _, err = ExecWithShellTimeout(buildTargetSSHCmd, timeout)
//		if err != nil {
//			return fmt.Errorf("cmd: %s, error: %s", buildTargetSSHCmd, err.Error())
//		}
//	} else {
//		buildOriginSSHCmd := fmt.Sprintf(`%s "echo '%s' >> ~/.ssh/authorized_keys; %s"`, targetCli, orginPub, noStrictKeyChecking)
//		_, _, err := ExecWithShellTimeout(buildOriginSSHCmd, timeout)
//		if err != nil {
//			return fmt.Errorf("cmd: %s, error: %s", buildOriginSSHCmd, err.Error())
//		}
//	}
//
//	//check
//	originCheckCmd = fmt.Sprintf(`%s "%s \"ls\""`, originCli, targetCli)
//	_, _, err = ExecWithShellTimeout(originCheckCmd, timeout)
//	if err != nil {
//		return fmt.Errorf("cmd: %s, error: %s", originCheckCmd, err.Error())
//	}
//	if isDuplex {
//		targetCheckCmd := fmt.Sprintf(`%s "%s \"ls\""`, targetCli, originCli)
//		_, _, err = ExecWithShellTimeout(targetCheckCmd, timeout)
//		if err != nil {
//			return fmt.Errorf("cmd: %s, error: %s", targetCheckCmd, err.Error())
//		}
//	}
//	return nil
//}
//
////GenerateIncrementIntList ...
//func GenerateIncrementIntList(maxLen int) []int {
//	res := make([]int, 0)
//	for i := 1; i <= maxLen; i++ {
//		res = append(res, i)
//	}
//
//	return res
//}
//
//// 两个已经排好序的无重复项列表，获取base列表中current最大值右侧的diff，该函数主要是用来比对两个集群xlog文件diff，恢复
//func SortedSliceStringMaxDiffItem(base, current []string) []string {
//	currentLast := current[len(current)-1]
//
//	diff := []string{}
//	findFlag := -1
//	c := 1
//	for _, v := range base {
//		if v == currentLast {
//			findFlag = 1
//		}
//		if findFlag >= 0 {
//			if c > 1 {
//				diff = append(diff, v)
//			}
//			c++
//		}
//	}
//
//	return diff
//}
//
////SSHObject ...
//type SSHObject struct {
//	Ctx             context.Context
//	IP              string //ssh连接目标机ip地址
//	Username        string //ssh使用账号
//	Password        string //ssh使用密码
//	ConnectUsername string //ssh连接目标机账号
//
//	Target           string //root@20.200.99.169
//	ConnectTarget    string //omm@20.200.99.169
//	SSH              string //ssh路径
//	OptionsIsArrived string // -n -o BatchMode=yes -o TCPKeepAlive=yes -o ConnectTimeout=2
//}
//
////NewSSH ...
//func NewSSH(ctx context.Context, ip, username, password, connectUsername string) *SSHObject {
//	if connectUsername == "" {
//		connectUsername = constant.UserOmmOfOS
//	}
//	ssh := &SSHObject{}
//	ssh.Ctx = ctx
//	ssh.IP = ip
//	ssh.Username = username
//	ssh.Password = password
//	ssh.ConnectUsername = connectUsername
//
//	ssh.Target = fmt.Sprintf("%s@%s", username, ip)
//	ssh.ConnectTarget = fmt.Sprintf("%s@%s", connectUsername, ip)
//	ssh.SSH = "ssh"
//	ssh.OptionsIsArrived = "-n -o BatchMode=yes -o TCPKeepAlive=yes -o ConnectTimeout=2"
//
//	return ssh
//}
//
////ExecCommandBySSHTrust 管控服务器到目标机默认是有互信的，直接ssh 执行命令，而不是用ansible。返回执行命令，标准输出，标准错误，error
//func (s *SSHObject) ExecCommandBySSHTrust(cmd string) (string, string, string, error) {
//	timeout := 5
//	optionsIsArrived := strings.Replace(s.OptionsIsArrived, "-o ConnectTimeout=2", "", 1)
//	command := fmt.Sprintf("%s %s %s true", s.SSH, optionsIsArrived, s.ConnectTarget)
//	stOut, stErr, err := ExecWithShellTimeout(command, time.Duration(timeout))
//	if err != nil {
//		return "", "", "", fmt.Errorf("check ssh trust [%s], stOut: %s, stErr: %s, error: %s",
//			s.ConnectTarget, stOut, stErr, err)
//	}
//	if stErr == "" && stOut == "" {
//		// 互信检查通过
//		command := fmt.Sprintf("%s %s \"%s\"", s.SSH, s.ConnectTarget, cmd)
//		sout, serr, err := ExecWithShellTimeout(command, 60)
//		if err != nil {
//			return command, "", "", fmt.Errorf("execute [%s], stOut: %s, stErr: %s, error: %s",
//				command, sout, serr, err)
//		}
//		return command, sout, serr, nil
//	}
//
//	return "", "", "", fmt.Errorf("failed to check ssh trust [%s], stOut: %s, stErr: %s,",
//		s.ConnectTarget, stOut, stErr)
//}
//
////PrepareMakeAuth ... true是没有互信，false是已经建立互信
//func (s *SSHObject) PrepareMakeAuth() (string, bool) {
//	timeout := 5
//	keyWord := "Permission denied"
//	var resultStr string
//	command := fmt.Sprintf("%s %s %s true", s.SSH, s.OptionsIsArrived, s.ConnectTarget)
//	stOut, stErr, err := ExecWithShellTimeout(command, time.Duration(timeout))
//	if err != nil {
//		resultStr = fmt.Sprintf("stErr: [%s], err: [%s]", stErr, err)
//		if p := strings.Contains(stErr, keyWord); p {
//			return resultStr, true
//		}
//
//		return resultStr, false
//	}
//
//	if stErr == "" && stOut == "" {
//		return "ssh already make auth with manager", false
//	}
//
//	resultStr = fmt.Sprintf("stOut: [%s], stErr: [%s]", stOut, stErr)
//	if p := strings.Contains(stErr, keyWord); p {
//		return resultStr, true
//	}
//	resultStr = fmt.Sprintf("stErr: [%s]", stErr)
//
//	return resultStr, false
//}
//
////MakePubkeyToTarget ex: sshpass -p 'Psbc!2017' ssh root@20.200.99.69 "date"，该函数必须写密码
//func (s *SSHObject) MakePubkeyToTarget(pubkey string) error {
//	targetAuthFile := ".ssh/authorized_keys"
//	targetAuthFilePath := fmt.Sprintf("/root/%s", targetAuthFile)
//	if s.ConnectUsername != constant.UserRootOfOS {
//		targetAuthFilePath = fmt.Sprintf("/home/%s/%s", s.ConnectUsername, targetAuthFile)
//	}
//
//	innerCmd := fmt.Sprintf("echo '%s' >> %s", pubkey, targetAuthFilePath)
//	command := fmt.Sprintf("sshpass -p '%s' ssh %s \"%s\"", s.Password, s.Target, innerCmd)
//	timeout := time.Duration(10)
//
//	execOut, execErr, err := ExecWithShellTimeout(command, timeout)
//	if err != nil {
//		execErr := fmt.Errorf("execute [%s]\n execOut: %s\n execErr: %s \n error: %s", command, execOut, execErr, err)
//		return execErr
//	}
//	err = CheckStringHaveWE(execOut)
//	if err != nil {
//		return err
//	}
//
//	return nil
//}
//
////MakePubkeyToTargetFromOther 互信环境下，对传参IP做互信，与new的ip地址无关
//func (s *SSHObject) MakePubkeyToTargetFromOther(pubkey, ip string) error {
//	targetAuthFile := ".ssh/authorized_keys"
//	targetAuthFilePath := fmt.Sprintf("/root/%s", targetAuthFile)
//	if s.ConnectUsername != constant.UserRootOfOS {
//		targetAuthFilePath = fmt.Sprintf("/home/%s/%s", s.ConnectUsername, targetAuthFile)
//	}
//
//	target := fmt.Sprintf("%s@%s", s.ConnectUsername, ip)
//	innerCmd := fmt.Sprintf("echo '%s' >> %s", pubkey, targetAuthFilePath)
//	command := fmt.Sprintf("ssh %s \"%s\"", target, innerCmd)
//	timeout := time.Duration(5)
//
//	execOut, execErr, err := ExecWithShellTimeout(command, timeout)
//	if err != nil {
//		execErr := fmt.Errorf("execute [%s]\n execOut: %s\n execErr: %s \n error: %s", command, execOut, execErr, err)
//		return execErr
//	}
//	err = CheckStringHaveWE(execOut)
//	if err != nil {
//		return err
//	}
//
//	return nil
//}
//
////StrictHostKeyChecking 在有root互信的情况下，添加或者删除 StrictHostKeyChecking配置。
//func (s *SSHObject) ModifyStrictHostKeyChecking(sw string) error {
//	// sw 是on或者off
//	noStrictKeyCheckingFile := "~/.ssh/config"
//
//	noStrictKeyCheckingAdd := fmt.Sprintf("echo 'StrictHostKeyChecking=no' >> %s", noStrictKeyCheckingFile)
//	noStrictKeyCheckingDrop := fmt.Sprintf("rm -f %s", noStrictKeyCheckingFile)
//
//	cmd := ""
//	if sw == "on" {
//		cmd = noStrictKeyCheckingAdd
//	} else if sw == "off" {
//		cmd = fmt.Sprintf("if [ -e %s ]; then %s; else echo 0;fi", noStrictKeyCheckingFile, noStrictKeyCheckingDrop)
//	}
//
//	command := fmt.Sprintf("ssh %s \"%s\"", s.ConnectTarget, cmd)
//	timeout := time.Duration(10)
//	execOut, execErr, err := ExecWithShellTimeout(command, timeout)
//	if err != nil {
//		execErr := fmt.Errorf("execute [%s]\n execOut: %s\n execErr: %s \n error: %s", command, execOut, execErr, err)
//		return execErr
//	}
//
//	return nil
//}
//
//func (s *SSHObject) GetPublicKey() (string, error) {
//	targetPublicKeyFile := ".ssh/id_rsa.pub"
//	targetGenPublicKeyFile := ".ssh/id_rsa"
//	targetPublicKeyFilePath := fmt.Sprintf("/root/%s", targetPublicKeyFile)
//	targetGenPublicKeyFilePath := fmt.Sprintf("/root/%s", targetGenPublicKeyFile)
//	if s.ConnectUsername != constant.UserRootOfOS {
//		targetPublicKeyFilePath = fmt.Sprintf("/home/%s/%s", s.ConnectUsername, targetPublicKeyFile)
//		targetGenPublicKeyFilePath = fmt.Sprintf("/home/%s/%s", s.ConnectUsername, targetGenPublicKeyFile)
//	}
//
//	generatePublicKey := fmt.Sprintf("ssh-keygen -t rsa -P '' -f %s <<< n", targetGenPublicKeyFilePath)
//	checkCommand := fmt.Sprintf("if [ -e %s ]; then echo 1; else %s;fi", targetPublicKeyFile, generatePublicKey)
//
//	command := fmt.Sprintf("ssh %s \"%s\"", s.ConnectTarget, checkCommand)
//	timeout := time.Duration(10)
//	execOut, execErr, err := ExecWithShellTimeout(command, timeout)
//	if err != nil {
//		execErr := fmt.Errorf("execute [%s]\n execOut: %s\n execErr: %s \n error: %s", checkCommand, execOut, execErr, err)
//		return "", execErr
//	}
//	// 获取key
//	getKeyCommand := fmt.Sprintf("cat %s", targetPublicKeyFilePath)
//	command = fmt.Sprintf("ssh %s \"%s\"", s.ConnectTarget, getKeyCommand)
//	execOut, execErr, err = ExecWithShellTimeout(command, timeout)
//	if err != nil {
//		execErr := fmt.Errorf("execute [%s]\n execOut: %s\n execErr: %s \n error: %s", checkCommand, execOut, execErr, err)
//		return "", execErr
//	}
//
//	return execOut, nil
//}
//
////MakeAuthEachOtherInstance IP列表相互做互信
//func (s *SSHObject) MakeAuthEachOtherInstanceFromManager(targetIP string) (string, error) {
//	timeout := 2
//	keyWord := "Permission denied"
//	if s.Password != "" {
//		return "", fmt.Errorf("not support when ssh target instance with password from manager, must be avoid password")
//	}
//	resultStr := ""
//	source := fmt.Sprintf("%s@%s", s.ConnectUsername, s.IP)
//	target := fmt.Sprintf("%s@%s", s.ConnectUsername, targetIP)
//	command := fmt.Sprintf("%s %s \"%s %s %s true\"", s.SSH, source, s.SSH, s.OptionsIsArrived, target)
//	stOut, stErr, err := ExecWithShellTimeout(command, time.Duration(timeout))
//	if err != nil {
//		resultStr = fmt.Sprintf("source %s to target %s check auth failed: stErr [%s], err [%s], ", source, target, stErr, err)
//		if p := strings.Contains(stErr, keyWord); p {
//			//需要做互信，第一步获取公钥
//			publicKey, err := s.GetPublicKey()
//			if err != nil {
//				resultStr += fmt.Sprintf("get source public key error: %s", err)
//				return "", errors.New(resultStr)
//			}
//			resultStr += fmt.Sprintf("get source public key: [%s], ", publicKey)
//			//第二步拷贝key
//			err = s.MakePubkeyToTargetFromOther(publicKey, targetIP)
//			if err != nil {
//				resultStr += fmt.Sprintf("make ssh auth to %s error: %s", targetIP, err)
//				return "", errors.New(resultStr)
//			}
//			resultStr += "copy public key success"
//
//			return resultStr, nil
//		}
//
//		return "", fmt.Errorf("source %s to target %s check auth unknown: stOut [%s] stErr [%s], err [%s]", source, target, stOut, stErr, err)
//	}
//	if stOut == "" {
//		resultStr = fmt.Sprintf("source %s to target %s check auth pass;", source, target)
//		return resultStr, nil
//	}
//
//	return "", fmt.Errorf("source %s to target %s check auth unknown: stOut [%s] stErr [%s]", source, target, stOut, stErr)
//}
