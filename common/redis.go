package common

import "github.com/redis/go-redis/v9"

func initRedis() {
	Rdb = redis.NewClient(&redis.Options{
		Addr:     Conf.Redis.Addr,
		Password: Conf.Redis.Password,
		DB:       Conf.Redis.DB,
	})
}
