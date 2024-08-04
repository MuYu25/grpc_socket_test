package websocket

import (
	"fmt"
	"runtime/debug"

	"github.com/gorilla/websocket"
)

const (
	heartbeatExpirationTime = 6 * 60
)

type login struct {
	AppID  uint32
	UserID string
	Client *Client
}

func (l *login) GetKey() (key string) {
	key = GetUserKey(l.AppID, l.UserID)
	return
}

type Client struct {
	Addr          string          // 客户端地址
	Socket        *websocket.Conn // 用户连接
	Send          chan []byte     // 待发送的数据
	AppID         uint32          // 登陆的平台 ID app/web/ios
	UserID        string          // 用户 ID，用户登陆后才有
	FirstTime     uint64          // 用户首次连接时间
	HeartbeatTime uint64          // 用户心跳时间
	LoginTime     uint64          // 用户登陆时间
}

func NewClient(addr string, socket *websocket.Conn, firstTime uint64) (client *Client) {
	client = &Client{
		Addr:          addr,
		Socket:        socket,
		Send:          make(chan []byte, 100),
		FirstTime:     firstTime,
		HeartbeatTime: firstTime,
	}
	return
}

func (c *Client) GetKey() (key string) {
	key = GetUserKey(c.AppID, c.UserID)
	return
}

func (c *Client) read() {
	defer func() {
		if r := recover(); r != nil {
			fmt.Println("write stop", string(debug.Stack()), r)
		}
	}()
	defer func() {
		fmt.Println("读取客户端数据，关闭send", c)
		close(c.Send)
	}()
	for {
		_, message, err := c.Socket.ReadMessage()
		if err != nil {
			fmt.Println("读取客户端数据 错误", c.Addr, err)
			return
		}
		fmt.Println("读取客户端数据 处理:", string(message))
		ProcessData(c, message)
	}
}

func (c *Client) write() {
	defer func() {
		if r := recover(); r != nil {
			fmt.Println("write stop", string(debug.Stack()), r)
		}
	}()
	defer func() {
		clientManager.Unregister <- c
		_ = c.Socket.Close()
		fmt.Println("Clinent发送数据 defer", c)
	}()
	for {
		select {
		case message, ok := <-c.Send:
			if !ok {
				fmt.Println("client 发送数据 关闭连接", c.Addr, "ok", ok)
				return
			}
			_ = c.Socket.WriteMessage(websocket.TextMessage, message)
		}
	}
}

func (c *Client) SendMsg(msg []byte) {
	if c == nil {
		return
	}
	defer func() {
		if r := recover(); r != nil {
			fmt.Println("SendMsg stop", string(debug.Stack()))
		}
	}()
	c.Send <- msg
}

func (c *Client) Close() {
	close(c.Send)
}

func (c *Client) Login(appID uint32, userID string, loginTime uint64) {
	c.AppID = appID
	c.UserID = userID
	c.LoginTime = loginTime
	// 登陆成功=心跳一次
	c.Heartbeat(loginTime)
}

func (c *Client) Heartbeat(heartbeatTime uint64) {
	c.HeartbeatTime = heartbeatTime
}

func (c *Client) IsHeartbeatTimeout(currentTime uint64) (timeout bool) {
	if c.HeartbeatTime+heartbeatExpirationTime <= currentTime {
		timeout = true
	}
	return
}

func (c *Client) IsLogin() (isLogin bool) {
	if c.UserID != "" {
		isLogin = true
		return
	}
	return
}
