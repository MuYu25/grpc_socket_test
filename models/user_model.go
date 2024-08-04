package models

import (
	"fmt"
	"time"
)

const (
	heartbeatTimeout = 3 * 60 // 用户心跳超时时间
)

type UserOnline struct {
	AccIp         string `json:"acc_ip"`           // acc ip
	AccPort       string `json:"acc_port"`         // acc端口
	AppID         uint32 `json:"appID"`            // appid
	UserID        string `json:"userID"`           // 用户ids
	ClientID      string `json:"clientID"`         // 客户端id
	ClientPort    string `json:"clientPort"`       // 客户端端口
	LoginTime     uint64 `json:"loginTime"`        // 用户上次登陆时间
	HeartbeatTime uint64 `json:"heartbeatTimeout"` // 用户上次心跳时间
	LogOutTime    uint64 `json:"logOutTime"`       // 用户登出时间
	Qua           string `json:"qua"`              // qua
	DeviceInfo    string `json:"deviceInfo"`       // 设备信息
	IsLogoff      bool   `json:"isLogoff"`         // 是否登出
}

func UserLogin(accIp, accPort string, appID uint32, userID string, addr string, loginTime uint64) (userOnline *UserOnline) {
	return &UserOnline{
		AccIp:         accIp,
		AccPort:       accPort,
		AppID:         appID,
		UserID:        userID,
		ClientID:      addr,
		ClientPort:    addr,
		LoginTime:     loginTime,
		HeartbeatTime: loginTime,
		IsLogoff:      false,
	}
}

func (u *UserOnline) Heartbeat(currentTime uint64) {
	u.HeartbeatTime = currentTime
	u.IsLogoff = false
}

func (u *UserOnline) Logout() {
	u.LogOutTime = uint64(time.Now().Unix())
	u.IsLogoff = true
}

func (u *UserOnline) IsOnline() (onlie bool) {
	if u.IsLogoff {
		return
	}
	currentTime := uint64(time.Now().Unix())
	if u.HeartbeatTime < (currentTime - heartbeatTimeout) {
		fmt.Println("用户是否在线 心跳超时", u.AppID, u.UserID, u.HeartbeatTime, currentTime)
		return
	}
	if u.IsLogoff {
		fmt.Println("用户是否在线 用户登出", u.AppID, u.UserID)
		return
	}
	return true
}

func (u *UserOnline) UserIsLocal(localIp, localPort string) bool {
	if u.AccIp == localIp && u.AccPort == localPort {
		return true
	}
	return false
}
