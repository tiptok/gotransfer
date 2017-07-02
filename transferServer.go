package main

import (
	"fmt"

	comm "github.com/tiptok/gotransfer/comm"
	conn "github.com/tiptok/gotransfer/conn"
)

var srv conn.TcpServer
var exit chan int

func main() {
	exit = make(chan int, 1)
	var param comm.Param
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

type transferSvrHandler struct {
}

func (trans *transferSvrHandler) OnConnect(c *conn.Connector) bool {
	return true
}

func (trans *transferSvrHandler) OnReceive(c *conn.Connector, d conn.TcpData) bool {
	fmt.Println("OnReceive Data.", len(d.Bytes()))
	return true
}

func (trans *transferSvrHandler) OnClose(c *conn.Connector) {

}
