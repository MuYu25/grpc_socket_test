package helper

import (
	"fmt"
	"time"
)

func GetOrderIDTime() (orderID string) {
	currentTime := time.Now().Nanosecond()
	orderID = fmt.Sprintf("%d", currentTime)
	return
}
