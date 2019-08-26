package conn

import (
	"context"
	"errors"
	"fmt"
	"log"
	"net"
	"strconv"
	"sync/atomic"
)

var (
	Error_SendDataLen = errors.New("gotransfer:SendEntity Error PackData Lenght <= 0")
)

func assertITcpServermplementation(){
	var _ ITcpServer = (*TcpServer)(nil)
}

//server
type TcpServer struct {
	Conn    []Connector
	Handler TcpHandler
	Config  *Conifg
	P       Protocol
	stopFlag        int32
	connectedFlag   int32
	listener net.Listener
}

//new tcpServer
func NewTcpServer(config *Conifg,handler TcpHandler,protocol Protocol)*TcpServer {
	return &TcpServer{
		Config:config,
		Handler:handler,
		P:protocol,
	}
}

//启动tcp服务
func (tcpServer *TcpServer) Start() error{
	defer func() {
		if err := recover(); err != nil {
			log.Println(err)
		}
	}()
	svrAddr := ":" + strconv.Itoa(tcpServer.Config.ListenPort)
	sAddr, err := net.ResolveTCPAddr("tcp4", svrAddr)
	//sAddr, err := net.ResolveTCPAddr("tcp4", "127.0.0.1"+":"+strconv.Itoa(tcpServer.Port))
	tcpServer.listener, err = net.Listen("tcp", sAddr.String())
	if err != nil {
		log.Println("Server Start Error." + err.Error())
	}
	log.Println("tcp listen ",sAddr.String())

	for {
		conn, err := tcpServer.listener.Accept()
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

//启动tcp服务
func (tcpServer *TcpServer) Stop() error{
	if !atomic.CompareAndSwapInt32(&tcpServer.stopFlag,0,1){
		return nil
	}
	tcpServer.listener.Close()
	return nil
}