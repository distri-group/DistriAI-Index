package common

import (
	"github.com/redis/go-redis/v9"
	"gorm.io/gorm"
)

var (
	Conf *Config
	Db   *gorm.DB
	Rdb  *redis.Client
)

func Init() {
	initConfig()
	initDatabase()
	initRedis()
}
