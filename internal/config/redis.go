package config

import (
	"github.com/redis/go-redis/v9"
	"github.com/spf13/viper"
)

func NewRedis(config *viper.Viper) *redis.Client {
	addr := config.GetString("redis.address")
	password := config.GetString("redis.password")
	db := config.GetInt("redis.db")

	rdb := redis.NewClient(&redis.Options{
		Addr:     addr,
		Password: password,
		DB:       db,
	})

	return rdb
}
