package conn

import (
	"context"
	"fmt"
	"log"
	"net"
	"strconv"
)

//server
type TcpServer struct {
	Port    int
	Conn    []Connector
	Handler TcpHandler
	config  *Conifg
	P       Protocol
	//Online  map[string]*Connector
}

//new tcpServer
func (tcpServer *TcpServer) NewTcpServer(port, sSize, rSize int) {
	tcpServer.Port = port
	tcpServer.config = &Conifg{
		SendSize:    500,
		ReceiveSize: 500,
	}
	//tcpServer.Conn = make([]Connector, 500)
	//tcpServer.Online = make(map[string]*Connector)
}

//启动tcp服务

func (tcpServer *TcpServer) Start(handler TcpHandler) {
	defer func() {
		if err := recover(); err != nil {
			log.Println(err)
		}
	}()
	tcpServer.Handler = handler
	svrAddr := ":" + strconv.Itoa(tcpServer.Port)
	sAddr, err := net.ResolveTCPAddr("tcp4", svrAddr)
	//sAddr, err := net.ResolveTCPAddr("tcp4", "127.0.0.1"+":"+strconv.Itoa(tcpServer.Port))
	listen, err := net.Listen("tcp", sAddr.String())
	if err != nil {
		log.Println("Server Start Error." + err.Error())
	}
	log.Println(sAddr.String(), "Tcp Start Listen.")

	for {
		conn, err := listen.Accept()
		if err != nil {
			fmt.Print("Recv Conn Error.", err.Error())
		} else {
			connector := NewConn(&conn, tcpServer.Handler, *tcpServer.config)
			connector.P = tcpServer.P
			//新连接 添加到在线列表里面
			// if _, exists := tcpServer.Online[connector.RemoteAddress]; !exists {
			// 	tcpServer.Online[connector.RemoteAddress] = connector
			// }
			ctx := context.Background()
			ctx, cancel := context.WithCancel(ctx)
			//传递 cancel
			connector.cancelFunc = cancel
			tcpServer.Handler.OnConnect(connector)
			go connector.ProcessRecv(ctx)
			go connector.DataHandler(ctx)
			go connector.ProcessSend(ctx)
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

func NewTcpData(d []byte) TcpData {
	return TcpData{buffer: d}
}
func (d TcpData) Bytes() []byte {
	return d.buffer
}

func (d TcpData) Lenght() int {
	return len(d.buffer)
}
