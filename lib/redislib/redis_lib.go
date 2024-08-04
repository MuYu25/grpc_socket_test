package redislib

import (
	"context"
	"fmt"

	"github.com/redis/go-redis/v9"
	"github.com/spf13/viper"
)

var (
	client *redis.Client
)

func NewClinet() {
	client = redis.NewClient(&redis.Options{
		Addr:         viper.GetString("redis.addr"),
		Password:     viper.GetString("redis.password"),
		DB:           viper.GetInt("redis.db"),
		PoolSize:     viper.GetInt("redis.poolsize"),
		MinIdleConns: viper.GetInt("redis.minidleconns"),
	})
	pong, err := client.Ping(context.Background()).Result()
	fmt.Println("初始化redis:", pong, err)
}

func GetClient() *redis.Client {
	return client
}
