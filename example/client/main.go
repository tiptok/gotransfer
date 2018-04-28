package main

import (
	"log"
	"time"

	"github.com/tiptok/gotransfer/conn"
)

type TcpClientHandler struct {
}

func (cli *TcpClientHandler) OnConnect(c *conn.Connector) bool {
	log.Println("client On connect", c.RemoteAddress)
	return true
}
func (cli *TcpClientHandler) OnReceive(c *conn.Connector, d conn.TcpData) bool {
	//cli.BizNode.SendToSrc(d)
	log.Println("Rece Data:%v", d)
	return true
}
func (cli *TcpClientHandler) OnClose(c *conn.Connector) {
	//从列表移除
	log.Println("Client OnClose:", c.RemoteAddress)
}

func main() {
	exit := make(chan int, 1)
	cli := &conn.TcpClient{}
	cli.NewTcpClient("127.0.0.1", 9928, 500, 500)
	cli.Start(&TcpClientHandler{})
	for {
		if !cli.Conn.IsConneted {
			cli.ReStart() //重连
		}
		log.Println("cli stop isconnnet:", cli.Conn.IsConneted, &cli)
		time.Sleep(10 * time.Second)
	}
	<-exit
}
