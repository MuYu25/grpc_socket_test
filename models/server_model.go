package models

import (
	"fmt"
	"strings"
)

type Server struct {
	Ip   string `json:"ip"`   // IP地址
	Port string `json:"port"` // 端口号
}

func NewServer(ip string, port string) *Server {
	return &Server{
		Ip:   ip,
		Port: port,
	}
}

func (s *Server) String() (str string) {
	if s == nil {
		return
	}
	str = fmt.Sprintf("%s:%s", s.Ip, s.Port)
	return str
}

func StringToServer(str string) (server *Server, err error) {
	list := strings.Split(str, ":")
	if len(list) != 2 {
		err = fmt.Errorf("invalid server string: %s", str)
		return
	}
	server = &Server{
		Ip:   list[0],
		Port: list[1],
	}
	return
}
