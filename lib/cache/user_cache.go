package cache

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"mygowebsockt/lib/redislib"
	"mygowebsockt/models"

	"github.com/redis/go-redis/v9"
)

const (
	userOnlinePrefix    = "acc:user:online:" // 用户在线状态
	userOnlineCacheTime = 24 * 60 * 60
)

func getUserOnlineKey(userKey string) (key string) {
	key = fmt.Sprintf("%s%s", userOnlinePrefix, userKey)
	return
}

func GetUserOnlineInfo(userKey string) (userOnline *models.UserOnline, err error) {
	redisClient := redislib.GetClient()
	key := getUserOnlineKey(userKey)
	data, err := redisClient.Get(context.Background(), key).Bytes()
	if err != nil {
		if errors.Is(err, redis.Nil) {
			fmt.Println("GetUserOnlineInfo", userKey, err)
			return
		}
		fmt.Println("GetUserOnlineInfo", userKey, err)
		return
	}
	userOnline = &models.UserOnline{}
	err = json.Unmarshal(data, userOnline)
	if err != nil {
		fmt.Println("获取用户在线数据 json Unmarshal", userKey, err)
		return
	}
	fmt.Println("获取用户在线数据", userKey, "time", userOnline.LoginTime, userOnline.HeartbeatTime,
		userOnline.AccIp, userOnline.IsLogoff)
	return
}

// 设置用户在线数据
func SetUserOnlineInfo(userKey string, userOnline *models.UserOnline) (err error) {
	redisClinet := redislib.GetClient()
	key := getUserOnlineKey(userKey)
	valueByte, err := json.Marshal(userOnline)
	if err != nil {
		fmt.Println("设置用户在线数据 json Marshal", userKey, err)
		return
	}
	_, err = redisClinet.Do(context.Background(), "setEx", key, userOnlineCacheTime, string(valueByte)).Result()
	if err != nil {
		fmt.Println("设置用户在线数据", userKey, err)
		return
	}
	return
}
