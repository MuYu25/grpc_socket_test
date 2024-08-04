package grpcserver

import (
	"context"
	"fmt"
	"log"
	"mygowebsockt/models"
	"mygowebsockt/protobuf"
	"mygowebsockt/servers/websocket"
	"mygowebsockt/setting"
	"net"

	"github.com/spf13/viper"
	"google.golang.org/grpc"
	"google.golang.org/protobuf/proto"
)

type server struct {
	protobuf.UnimplementedAccServerServer
}

func setErr(rsp proto.Message, code uint32, message string) {
	message = setting.GetErrorMessage(code, message)
	switch v := rsp.(type) {
	case *protobuf.QueryUserOnlineRsp:
		v.RetCode = code
		v.ErrMsg = message
	case *protobuf.SendMsgRsp:
		v.RetCode = code
		v.ErrMsg = message
	case *protobuf.SendMsgAllRsp:
		v.RetCode = code
		v.ErrMsg = message
	case *protobuf.GetUserListRsp:
		v.RetCode = code
		v.ErrMsg = message
	default:
	}
}

func (s *server) QueryUserOnline(c context.Context, req *protobuf.QueryUserOnlineReq) (rsp *protobuf.QueryUserOnlineRsp, err error) {
	fmt.Println("grpc_request 查询用户是否在线", req.String())
	rsp = &protobuf.QueryUserOnlineRsp{}
	online := websocket.CheckUserOnline(req.GetAppID(), req.GetUserID())
	setErr(req, setting.OK, "")
	rsp.Online = online
	return rsp, nil
}

func (s *server) SendMsg(c context.Context, req *protobuf.SendMsgReq) (rsp *protobuf.SendMsgRsp, err error) {
	fmt.Println("grpc_request 给本机用户发送消息", req.String())
	rsp = &protobuf.SendMsgRsp{}
	if req.GetIsLocal() {
		setErr(rsp, setting.ParameterIllegal, "")
		return
	}
	data := models.GetMsgData(req.GetUserID(), req.GetSeq(), req.GetCms(), req.GetMsg())
	sendResults, err := websocket.SendUserMessageLocal(req.GetAppID(), req.GetUserID(), data)
	if err != nil {
		fmt.Println("系统错误", err)
		setErr(rsp, setting.ServerError, "")
		return rsp, nil
	}
	if !sendResults {
		fmt.Println("发送失败", err)
		setErr(rsp, setting.OperationFailure, "")
		return rsp, nil
	}
	setErr(rsp, setting.OK, "")
	fmt.Println("grpc_response 给本机用户发送消息成功", rsp.String())
	return
}

func (s *server) SendMsgAll(c context.Context, req *protobuf.SendMsgAllReq) (rsp *protobuf.SendMsgAllRsp, err error) {
	fmt.Println("grpc_request 给本机所有用户发送消息", req.String())
	rsp = &protobuf.SendMsgAllRsp{}
	data := models.GetMsgData(req.GetUserID(), req.GetSeq(), req.GetCms(), req.GetMsg())
	websocket.AllSendMessages(req.GetAppID(), req.GetUserID(), data)
	setErr(rsp, setting.OK, "")
	fmt.Println("grpc_response 给本机所有用户发送消息成功", rsp.String())
	return
}

func (s *server) GetUserList(c context.Context, req *protobuf.GetUserListReq) (rsp *protobuf.GetUserListRsp, err error) {
	fmt.Println("grpc_request 获取本机用户列表", req.String())
	appID := req.GetAppID()
	rsp = &protobuf.GetUserListRsp{}
	userList := websocket.GetUserList(appID)
	setErr(rsp, setting.OK, "")
	rsp.UserID = userList
	fmt.Println("grpc_response 获取本机用户列表成功", rsp.String())
	return
}

func Init() {
	rpcProt := viper.GetString("app.rpcProt")
	fmt.Println("grpc server 启动", rpcProt)
	lis, err := net.Listen("tcp", ":"+rpcProt)
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}
	s := grpc.NewServer()
	protobuf.RegisterAccServerServer(s, &server{})
	if err := s.Serve(lis); err != nil {
		log.Fatalf("failed to serve: %v", err)
	}
}
