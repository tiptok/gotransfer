package conn

import (
	"context"
	"errors"
	"fmt"
	"log"
	"net"
	"strconv"
)

var (
	Error_SendDataLen = errors.New("gotransfer:SendEntity Error PackData Lenght <= 0")
)

//server
type TcpServer struct {
	Port    int
	Conn    []Connector
	Handler TcpHandler
	Config  *Conifg
	P       Protocol
	//Online  map[string]*Connector
}

//new tcpServer
func (tcpServer *TcpServer) NewTcpServer(port, sSize, rSize int) {
	tcpServer.Port = port
	tcpServer.Config = &Conifg{
		SendSize:    500,
		ReceiveSize: 500,
		PackageSize: 1024,
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
			connector := NewConn(&conn, tcpServer.Handler, *tcpServer.Config)
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

/*统一Server接口*/
type ITcpServer interface {
	Start() bool
	Stop() bool
}

type TcpServerBase struct {
	Server TcpServer
}

func SendEntity(e interface{}, c *Connector) (data []byte, err error) {
	data, err = c.P.PacketMsg(e)
	if err == nil {
		if len(data) > 0 {
			send := NewTcpData(data)
			c.SendChan <- send
		} else {
			err = Error_SendDataLen
		}
	}
	return data, err
}
