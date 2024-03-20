package common

import (
	"distriai-index-solana/common/s3actions"
	"github.com/redis/go-redis/v9"
	"gorm.io/gorm"
)

var (
	Conf        *Config
	Db          *gorm.DB
	Rdb         *redis.Client
	S3Presigner s3actions.Presigner
)

func Init() {
	initConfig()
	initDatabase()
	initRedis()
	initS3()
}
