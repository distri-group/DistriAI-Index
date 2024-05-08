package common

import (
	"github.com/fsnotify/fsnotify"
	"github.com/spf13/viper"
	"os"
)

type Config struct {
	Server   Server
	Database Database
	Redis    Redis
	Mailbox  Mailbox
	Aws      Aws
	Chain    Chain
}

type Server struct {
	Mode string
	Port string
}

type Database struct {
	Host     string
	Port     string
	Username string
	Password string
	Database string
}

type Redis struct {
	Addr     string
	Password string
	DB       int
}

type Mailbox struct {
	Host     string
	Port     int
	Username string
	Password string
}

type Aws struct {
	Region          string
	AccessKeyId     string
	SecretAccessKey string
}

type Chain struct {
	Rpc              string
	ProgramId        string
	FaucetPrivateKey string
	Dist             string
	DistDecimals     uint8
	DistFaucetAmount uint64
}

func initConfig() {
	viper.SetConfigType("yaml")
	configEnv := os.Getenv("GO_ENV")
	switch configEnv {
	case "dev":
		viper.SetConfigFile("config/config-dev.yml")
	default:
		viper.SetConfigFile("config/config.yml")
	}

	err := viper.ReadInConfig()
	if err != nil {
		panic(err)
	}
	err = viper.Unmarshal(&Conf)
	if err != nil {
		panic(err)
	}
	viper.WatchConfig()
	viper.OnConfigChange(func(in fsnotify.Event) {
		err = viper.Unmarshal(&Conf)
		if err != nil {
			panic(err)
		}
	})
}
