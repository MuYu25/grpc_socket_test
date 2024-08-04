package websocket

import (
	"fmt"
	"mygowebsockt/helper"
	"mygowebsockt/lib/cache"
	"mygowebsockt/models"
	"sync"
	"time"
)

type ClientManager struct {
	Clients     map[*Client]bool   // 全部的连接
	ClientsLock sync.RWMutex       // 读写锁
	Users       map[string]*Client // 登陆的用户 // appID+uuid
	UsersLock   sync.RWMutex       // 读写锁
	Register    chan *Client       // 连接连接处理
	Login       chan *login        // 登陆处理
	Unregister  chan *Client       // 断开连接处理
	Broadcast   chan []byte        // 广播处理
}

func NewClientManager() (clientManager *ClientManager) {
	clientManager = &ClientManager{
		Clients:    make(map[*Client]bool),
		Users:      make(map[string]*Client),
		Register:   make(chan *Client),
		Login:      make(chan *login),
		Unregister: make(chan *Client),
		Broadcast:  make(chan []byte),
	}
	return
}

func GetUserKey(appID uint32, userID string) (key string) {
	key = fmt.Sprintf("%d_%s", appID, userID)
	return
}

func (manager *ClientManager) InClient(client *Client) (ok bool) {
	manager.ClientsLock.RLock()
	defer manager.ClientsLock.RUnlock()

	// 连接存在，再添加
	_, ok = manager.Clients[client]
	return
}

func (manager *ClientManager) GetClients() (clients map[*Client]bool) {
	clients = make(map[*Client]bool)
	manager.ClientsRange(func(client *Client, value bool) (result bool) {
		clients[client] = value
		return true
	})
	return
}

// ClientRange 遍历
func (manager *ClientManager) ClientsRange(f func(client *Client, value bool) (result bool)) {
	manager.ClientsLock.RLock()
	defer manager.ClientsLock.RUnlock()
	for key, value := range manager.Clients {
		result := f(key, value)
		if !result {
			return
		}
	}
}

func (manager *ClientManager) GetClientLen() (clientsLen int) {
	clientsLen = len(manager.Clients)
	return
}

func (manager *ClientManager) AddClients(client *Client) {
	manager.ClientsLock.Lock()
	defer manager.ClientsLock.Unlock()
	manager.Clients[client] = true
}

func (manager *ClientManager) DelClients(client *Client) {
	manager.ClientsLock.Lock()
	defer manager.ClientsLock.Unlock()
	if _, ok := manager.Clients[client]; ok {
		delete(manager.Clients, client)
	}
}
func (manaager *ClientManager) GetUserClient(appID uint32, userID string) (client *Client) {
	manaager.UsersLock.RLock()
	defer manaager.UsersLock.RUnlock()
	userKey := GetUserKey(appID, userID)
	if value, ok := manaager.Users[userKey]; ok {
		client = value
	}
	return
}

func (manager *ClientManager) GetUserLen() (userLen int) {
	userLen = len(manager.Users)
	return
}

func (manager *ClientManager) AddUsers(key string, client *Client) {
	manager.UsersLock.Lock()
	defer manager.UsersLock.Unlock()
	manager.Users[key] = client
}

func (manager *ClientManager) DelUsers(client *Client) (result bool) {
	manager.UsersLock.Lock()
	defer manager.UsersLock.Unlock()
	key := GetUserKey(client.AppID, client.UserID)
	if value, ok := manager.Users[key]; ok {
		if value.Addr != client.Addr {
			return
		}
		delete(manager.Users, key)
		result = true
	}
	return
}

func (manager *ClientManager) GetUserKeys() (userKeys []string) {
	userKeys = make([]string, 0, len(manager.Users))
	for key := range manager.Users {
		userKeys = append(userKeys, key)
	}
	return
}
func (manager *ClientManager) GetUserList(appID uint32) (userList []string) {
	userList = make([]string, 0)
	manager.UsersLock.RLock()
	defer manager.UsersLock.RUnlock()
	for _, v := range manager.Users {
		if v.AppID == appID {
			userList = append(userList, v.UserID)
		}
	}
	fmt.Println("GetUserList len:", len(manager.Users))
	return
}

func (manager *ClientManager) GetUserLists() (clients []*Client) {
	clients = make([]*Client, 0)
	manager.UsersLock.RLock()
	defer manager.UsersLock.RUnlock()
	for _, v := range manager.Users {
		clients = append(clients, v)
	}
	return
}

// sendll 向全部成员（除了自己）发送数据
func (manageer *ClientManager) sendAll(message []byte, ignoreClient *Client) {
	clients := manageer.GetUserLists()
	for _, conn := range clients {
		if conn != ignoreClient {
			conn.SendMsg(message)
		}
	}
}

func (manager *ClientManager) sendAppIDAll(message []byte, appID uint32, ignoreClient *Client) {
	clients := manager.GetUserLists()
	for _, conn := range clients {
		if conn != ignoreClient && conn.AppID == appID {
			conn.SendMsg(message)
		}
	}
}

func (manager *ClientManager) EventRegister(client *Client) {
	manager.AddClients(client)
	fmt.Println("EventRegister:", client.Addr)
	// client.Send <- []byte("连接成功")
}

func (manager *ClientManager) EventLogin(login *login) {
	client := login.Client

	if manager.InClient(client) {
		userKey := login.GetKey()
		manager.AddUsers(userKey, login.Client)
	}
	fmt.Println("EventLogin 用户登陆", client.Addr, login.AppID)
	orderID := helper.GetOrderIDTime()
	_, _ = SendUserMessageAll(login.AppID, login.UserID, orderID, models.MessageCmdEnter, "哈喽～")
}

func (manager *ClientManager) EventUnregister(client *Client) {
	manager.DelClients(client)
	deleteResult := manager.DelUsers(client)
	if !deleteResult {
		return
	}
	userOnline, err := cache.GetUserOnlineInfo(client.GetKey())
	if err == nil {
		userOnline.Logout()
		_ = cache.SetUserOnlineInfo(client.GetKey(), userOnline)
	}
	fmt.Println("EventUnregister 用户断开连接", client.Addr, client.AppID, client.UserID)
	if client.UserID != "" {
		orderID := helper.GetOrderIDTime()
		_, _ = SendUserMessageAll(client.AppID, client.UserID, orderID, models.MessageCmdExit, "用户已经离开～")
	}
}

// 管道处理程序
func (manager *ClientManager) start() {
	for {
		select {
		case conn := <-manager.Register:
			// 建立连接事件
			manager.EventRegister(conn)
		case l := <-manager.Login:
			// 登录事件
			manager.EventLogin(l)
		case conn := <-manager.Unregister:
			// 断开连接事件
			manager.EventUnregister(conn)
		case message := <-manager.Broadcast:
			clients := manager.GetClients()
			for conn := range clients {
				select {
				case conn.Send <- message:
				default:
					close(conn.Send)
				}
			}
		}
	}
}

func GetManagerInfo(isDebug string) (managerInfo map[string]interface{}) {
	managerInfo = make(map[string]interface{})
	managerInfo["clientsLen"] = clientManager.GetClientLen()         // 客户端列表
	managerInfo["usersLen"] = clientManager.GetUserLen()             // 登陆用户列表
	managerInfo["chanRegisterLen"] = len(clientManager.Register)     //  未处理连接事件数
	managerInfo["chanLoginLen"] = len(clientManager.Login)           // 未处理登录事件数
	managerInfo["chanUnregisterLen"] = len(clientManager.Unregister) // 未处理断开连接事件数
	managerInfo["chanBroadcastLen"] = len(clientManager.Broadcast)   // 未处理广播事件数
	if isDebug == "true" {
		addrList := make([]string, 0)
		clientManager.ClientsRange(func(client *Client, value bool) (result bool) {
			addrList = append(addrList, client.Addr)
			return true
		})
		users := clientManager.GetUserKeys()
		managerInfo["clients"] = addrList // 客户端列表
		managerInfo["users"] = users      // 登陆用户列表
	}
	return
}

func GetUserClient(appID uint32, userID string) (client *Client) {
	client = clientManager.GetUserClient(appID, userID)
	return
}

func ClearTimeoutConnections() {
	currentTime := uint64(time.Now().Unix())
	clients := clientManager.GetClients()
	for client := range clients {
		if client.IsHeartbeatTimeout(currentTime) {
			fmt.Println("心跳时间超时，断开连接", client.Addr, client.UserID, client.LoginTime, client.HeartbeatTime)
			_ = client.Socket.Close()
		}
	}
}

func GetUserList(appID uint32) (userList []string) {
	fmt.Println("获取全部用户", appID)
	userList = clientManager.GetUserList(appID)
	return
}

func AllSendMessages(appID uint32, userID string, data string) {
	fmt.Println("广播消息", appID, userID, data)
	ignoreClient := clientManager.GetUserClient(appID, userID)
	clientManager.sendAppIDAll([]byte(data), appID, ignoreClient)
}
