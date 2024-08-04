package cache

import (
	"context"
	"encoding/json"
	"fmt"
	"mygowebsockt/lib/redislib"
	"mygowebsockt/models"
	"strconv"
)

const (
	serverHashKey       = "acc:hash:servers" // 全部的服务器
	serverHashCacheTime = 2 * 60 * 60        // key过期时间
	serverHahsTimeout   = 3 * 60             // 超时时间
)

func getServerHashKey() (key string) {
	key = fmt.Sprintf("%s", serverHashKey)
	return
}

// 设置服务器信息
func SetServerInfo(server *models.Server, currentTime uint64) (err error) {
	key := getServerHashKey()
	value := fmt.Sprintf("%d", currentTime)
	redisClient := redislib.GetClient()
	number, err := redisClient.Do(context.Background(), "hSet", key, server.String(), value).Int()
	if err != nil {
		fmt.Println("SetServerInfo", key, number, err)
		return
	}
	redisClient.Do(context.Background(), "Expire", key, serverHashCacheTime)
	return
}

// 下线服务器
func DelServerInfo(server *models.Server) (err error) {
	key := getServerHashKey()
	redisClient := redislib.GetClient()
	number, err := redisClient.Do(context.Background(), "hDel", key, server.String()).Int()
	if err != nil {
		fmt.Println("DelServerInfo", key, number, err)
		return
	}
	if number != 1 {
		return
	}
	redisClient.Do(context.Background(), "Expire", key, serverHashCacheTime)
	return
}

// 获取所有服务器
func GetServerAll(currentTime uint64) (servers []*models.Server, err error) {
	servers = make([]*models.Server, 0)
	key := getServerHashKey()
	redisClient := redislib.GetClient()
	val, err := redisClient.Do(context.Background(), "hGetAll", key).Result()
	valueByte, _ := json.Marshal(val)
	fmt.Println("GetServerAll", key, string(valueByte))
	serverMap, err := redisClient.HGetAll(context.Background(), key).Result()
	if err != nil {
		fmt.Println("SetServerInfo", key, err)
		return
	}
	for key, value := range serverMap {
		valueUint64, err := strconv.ParseUint(value, 10, 64)
		if err != nil {
			fmt.Println("GetServerAll", key, err)
			return nil, err
		}
		if valueUint64+serverHahsTimeout < currentTime {
			continue
		}
		server, err := models.StringToServer(key)
		if err != nil {
			fmt.Println("GetServerAll", key, err)
			return nil, err
		}
		servers = append(servers, server)
	}
	return
}
