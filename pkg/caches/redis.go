package caches

import (
	"github.com/go-redis/redis/v8"
	"nftshopping-store-api/pkg/config"
)

var (
	redisInstance *redis.Client
)

func GetRedis() (instance *redis.Client, err error) {
	if redisInstance == nil {
		instance, err = NewRedis()
		if err != nil {
			return nil, err
		}
		redisInstance = instance
	}
	return redisInstance, nil
}

func NewRedis() (instance *redis.Client, err error) {
	c, err := config.GetConfig()
	if err != nil {
		return nil, err
	}
	redisConfig := c.Database.Redis
	opt := &redis.Options{
		Addr: redisConfig.Uri,
		DB:   0, // use default DB
	}
	if len(redisConfig.Password) > 0 {
		opt.Password = redisConfig.Password
	}
	return redis.NewClient(opt), nil
}
