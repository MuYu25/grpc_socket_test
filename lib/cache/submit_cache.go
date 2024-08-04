package cache

import (
	"context"
	"fmt"
	"mygowebsockt/lib/redislib"
)

const (
	submitAgainPrefix = "acc:submit:again" // 数据不重复提交
)

func getSubmitAgainKey(from string, value string) (key string) {
	key = fmt.Sprintf("%s%s:%s", submitAgainPrefix, from, value)
	return
}

func submitAgain(from string, second int, value string) (isSubmitAgain bool) {
	// 默认重复提交
	isSubmitAgain = true
	key := getSubmitAgainKey(from, value)
	redisClient := redislib.GetClient()
	number, err := redisClient.Do(context.Background(), "setNx", key, "1").Int()
	if err != nil {
		fmt.Println("submitAgain", key, number, err)
		return
	}
	if number != 1 {
		return
	}
	isSubmitAgain = false
	redisClient.Do(context.Background(), "Expire", key, second)
	return
}

func SeqDupliCates(seq string) (result bool) {
	result = submitAgain("seq", 12*60*60, seq)
	return
}
