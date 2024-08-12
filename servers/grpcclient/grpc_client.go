package grpcclient

import (
	"context"
	"fmt"
	"mygowebsockt/models"
	"mygowebsockt/protobuf"
	"mygowebsockt/setting"
	"time"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

func SendMsgAll(server *models.Server, seq string, appID uint32, userID string, cmd string, message string) (sendMsgID string, err error) {
	// conn, err := grpc.Dial(server.String(), grpc.WithTransportCredentials(insecure.NewCredentials()))
	// conn, err := grpc.DialContext(context.Background(), server.String(), grpc.WithTransportCredentials(insecure.NewCredentials()))
	conn, err := grpc.NewClient(server.String(), grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		fmt.Println("连接失败", server.String())
		return
	}
	defer func() { _ = conn.Close() }()
	c := protobuf.NewAccServerClient(conn)
	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()
	req := protobuf.SendMsgAllReq{
		Seq:    seq,
		AppID:  appID,
		UserID: userID,
		Cms:    cmd,
		Msg:    message,
	}
	rsp, err := c.SendMsgAll(ctx, &req)
	if err != nil {
		fmt.Println("给全体用户发送消息", err)
		return
	}

	// 可能存在问题的地方
	if rsp.GetRetCode() != setting.OK {
		fmt.Println("给全体用户发送消息失败", rsp.String())
		err = fmt.Errorf("发送消息失败 code: %d", rsp.GetRetCode())
		return
	}
	sendMsgID = rsp.GetSendMsgID()
	fmt.Println("给全体用户发送消息成功", sendMsgID)
	return
}

func GetUserList(server *models.Server, appID uint32) (userIDs []string, err error) {
	userIDs = make([]string, 0)
	// conn, err := grpc.Dial(server.String(), grpc.WithTransportCredentials(insecure.NewCredentials()))
	conn, err := grpc.NewClient(server.String(), grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		fmt.Println("连接失败", server.String())
		return
	}
	defer func() { _ = conn.Close() }()
	c := protobuf.NewAccServerClient(conn)
	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()
	req := protobuf.GetUserListReq{
		AppID: appID,
	}
	rsp, err := c.GetUserList(ctx, &req)
	if err != nil {
		fmt.Println("获取用户列表失败", err)
		return
	}
	if rsp.GetRetCode() != setting.OK {
		fmt.Println("获取用户列表失败", rsp.String())
		err = fmt.Errorf("获取用户列表失败 code: %d", rsp.GetRetCode())
		return
	}
	userIDs = rsp.GetUserID()
	fmt.Println("获取用户列表成功", userIDs)
	return
}

func SendMsg(server *models.Server, seq string, appID uint32, userID string, cmd string, msgType string, message string) (sendMsgID string, err error) {
	// conn, err := grpc.Dial(server.String(), grpc.WithTransportCredentials(insecure.NewCredentials()))
	conn, err := grpc.NewClient(server.String(), grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		fmt.Println("连接失败", server.String())
		return
	}
	defer func() { _ = conn.Close() }()
	c := protobuf.NewAccServerClient(conn)
	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()
	req := protobuf.SendMsgReq{
		Seq:     seq,
		AppID:   appID,
		UserID:  userID,
		Cms:     cmd,
		Type:    msgType,
		Msg:     message,
		IsLocal: false,
	}
	rsp, err := c.SendMsg(ctx, &req)
	if err != nil {
		fmt.Println("给用户发送消息", err)
		return
	}
	if rsp.GetRetCode() != setting.OK {
		fmt.Println("发送消息", rsp.String())
		err = fmt.Errorf("发送消息失败 code: %d", rsp.GetRetCode())
		return
	}
	sendMsgID = rsp.GetSendMsgID()
	fmt.Println("给用户发送消息成功", sendMsgID)
	return
}
