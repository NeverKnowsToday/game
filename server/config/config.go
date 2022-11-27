package config

import (
	"github.com/game/server/logger/logging"
	"github.com/pkg/errors"
	"github.com/spf13/viper"
	"strings"
)

var myViper = viper.New()

var logger = logging.GetLogger("config", logging.DEFAULT_LEVEL)

const (
	cmdRoot = "server"
)

var (
	serverCfg *ServerConfig
)

type DefaultConfig struct {
	ExpirationTime          string
	PushIntervalTime        int
	SyncInvokeTime          int
	GetDataFromWisChainTime int
	JudgeReconnectTime      int
	WordRule                string
	WordFormatRule          string
	WordLength              int
	ContextTime             int
	MachineSecret			string
}

type DbConfig struct {
	Address  string
	Port     int
	Type     string
	Name     string
	User     string
	Password string
	PageSize int
}


type Config struct {
	AdminUser          string
	AdminPassWord      string
}

type Crypto struct {
	Family    string
	Algorithm string
	Hash      string
}

type ServerConfig struct {
	Default DefaultConfig
	Db      DbConfig
	Port    int
	Config  Config
	Crypto  Crypto
}


// InitConfig ...
// initConfig reads in config file
func InitConfig(configFile string) error {
	return InitConfigWithCmdRoot(configFile, cmdRoot)
}

// InitConfigWithCmdRoot reads in a config file and allows the
// environment variable prefixed to be specified
func InitConfigWithCmdRoot(configFile string, cmdRootPrefix string) error {
	myViper.SetEnvPrefix(cmdRootPrefix)
	myViper.AutomaticEnv()
	replacer := strings.NewReplacer(".", "_")
	myViper.SetEnvKeyReplacer(replacer)

	if configFile != "" {
		// create new viper
		myViper.SetConfigFile(configFile)
		// If a config file is found, read it in.
		err := myViper.ReadInConfig()

		if err == nil {
			logger.Infof("Using config file: %s", myViper.ConfigFileUsed())
		} else {
			return errors.Wrap(err, "Fatal error config file")
		}
	}

	serverCfg = &ServerConfig{}
	// Unmarshal the config into 'serverCfg'
	err := myViper.Unmarshal(serverCfg)
	if err != nil {
		return errors.Wrapf(err, "Incorrect format in file '%s'", configFile)
	}
	logger.Debugf("************************ config **************************\n")
	logger.Debugf("%#v\n", serverCfg)
	logger.Debugf("************************ config **************************\n")
	return nil
}

func GetServerConfig() *ServerConfig {
	return serverCfg
}
