package conn

import (
	"fmt"
	"log"
	"net"
	"strconv"
)

//server
type TcpServer struct {
	Ip       string
	Port     int
	Conn     []Connector
	Handler  TcpHandler
	config   *Conifg
	protocol Protocol
	Online   map[string]*Connector
}

//new tcpServer
func (tcpServer *TcpServer) NewTcpServer(ip string, port int) {
	tcpServer.Ip = ip
	tcpServer.Port = port
	tcpServer.config = &Conifg{
		SendSize:    500,
		ReceiveSize: 500,
	}
	//tcpServer.Conn = make([]Connector, 500)
	tcpServer.Online = make(map[string]*Connector)
}

//启动tcp服务
func (tcpServer *TcpServer) Start(handler TcpHandler) {

	tcpServer.Handler = handler
	sAddr := tcpServer.Ip + ":" + strconv.Itoa(tcpServer.Port)
	listen, err := net.Listen("tcp", sAddr)
	if err != nil {
		log.Fatal("Server Start Error." + err.Error())
	}
	log.Fatal(tcpServer.Ip, ":", tcpServer.Port, "Start Listen.")

	for {
		conn, err := listen.Accept()
		if err != nil {
			fmt.Print("Recv Conn Error.", err.Error())
		} else {
			//tcpconn ? netconn?
			//tcpServer.Conn = append(tcpServer.Conn, conn)
			connector := NewConn(&conn, tcpServer)

			//新连接 添加到在线列表里面
			if _, exists := tcpServer.Online[connector.RemoteAddress]; !exists {
				tcpServer.Online[connector.RemoteAddress] = connector
			}

			tcpServer.Handler.OnConnect(connector)
			go connector.ProcessRecv()
			go connector.DataHandler()
			go connector.ProcessSend()
		}
	}

}

//回调处理事件
type TcpHandler interface {
	OnConnect(*Connector) bool

	OnReceive(*Connector, TcpData) bool

	OnClose(*Connector)
}

//tcpData
type TcpData struct {
	buffer []byte
}
