package conn

import (
	"fmt"
	"log"
	"net"
	"strconv"
)

//server
type TcpServer struct {
	Port     int
	Conn     []Connector
	Handler  TcpHandler
	config   *Conifg
	protocol Protocol
	Online   map[string]*Connector
}

//new tcpServer
func (tcpServer *TcpServer) NewTcpServer(port, sSize, rSize int) {
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
	sAddr, err := net.ResolveTCPAddr("tcp4", "127.0.0.1"+":"+strconv.Itoa(tcpServer.Port))
	//sAddr := tcpServer.Ip + ":" + strconv.Itoa(tcpServer.Port)
	listen, err := net.Listen("tcp", sAddr.String())
	if err != nil {
		log.Println("Server Start Error." + err.Error())
	}
	log.Println(sAddr.String(), "Start Listen.")

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

func (d TcpData) Bytes() []byte {
	return d.buffer
}

func (d TcpData) Lenght() int {
	return len(d.buffer)
}
