package main

import (
	"fmt"

	"time"

	comm "github.com/tiptok/gotransfer/comm"
	conn "github.com/tiptok/gotransfer/conn"
)

var srv conn.TcpServer
var exit chan int
var param comm.Param
var tcpTransferList map[string][]conn.Connector

func main() {
	tcpTransferList = make(map[string][]conn.Connector)
	exit = make(chan int, 1)
	//F:/App/Go/WorkSpace/src/github.com/tiptok/gotransfer/
	err := comm.ReadConfig(&param, "param.json")
	if err != nil {
		fmt.Println(err.Error())
	}
	fmt.Println(param)
	//启动服务
	srv.NewTcpServer(param.Tcpsrv.Port, param.Tcpsrv.Sendsize, param.Tcpsrv.Recvsize)
	srv.Start(&transferSvrHandler{})
	<-exit
	fmt.Println("App Stop.")
}

/*
连接事件
tpc server 使用
*/
type transferSvrHandler struct {
}

func (trans *transferSvrHandler) OnConnect(c *conn.Connector) bool {
	for _, srv := range param.Transfers {
		var tcpClient conn.TcpClient
		tcpClient.NewTcpClient(srv.Ip, srv.Port, 500, 500)

		key :=c.RemoteAddress+"1"
		if _,ok :=map[key];ok{
			
		}
		else{
			map[key] =make([]conn.Connector,1)
		}
	}
	return true
}
func (trans *transferSvrHandler) OnReceive(c *conn.Connector, d conn.TcpData) bool {
	fmt.Println(time.Now(), c.RemoteAddress, " OnReceive Data.", string(d.Bytes()))
	return true
}
func (trans *transferSvrHandler) OnClose(c *conn.Connector) {
}
