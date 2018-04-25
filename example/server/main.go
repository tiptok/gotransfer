package main

import (
	"encoding/hex"
	"gotransfer/conn"
	"log"
)

type SimpleServerHandler struct {
}

//OnConnect handler
func (trans *SimpleServerHandler) OnConnect(c *conn.Connector) bool {
	log.Println(c.RemoteAddress, "On Connect.")
	return true
}

//OnClose handler
func (trans *SimpleServerHandler) OnClose(c *conn.Connector) {
	log.Println(c.RemoteAddress, "On Close.")
}

var exit chan int

//OnReceive handler
func (trans *SimpleServerHandler) OnReceive(c *conn.Connector, d conn.TcpData) bool {
	log.Printf("%v On Receive Data : %v\n", c.RemoteAddress, hex.EncodeToString(d.Bytes()))
	c.SendChan <- d
	return true
}

func main() {
	//启动tcp服务
	go func() {
		var srv conn.TcpServer
		srv.NewTcpServer(9928, 500, 500)
		srv.Start(&SimpleServerHandler{})
	}()
	//等待退出
	<-exit
}
