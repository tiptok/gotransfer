package conn

import (
	"context"
	"errors"
	"fmt"
	"github.com/tiptok/gotransfer/comm"
	"log"
	"net"
	"strconv"
	"sync/atomic"
)

var (
	Error_SendDataLen = errors.New("gotransfer:SendEntity Error PackData Lenght <= 0")
)
var (
	DefaultPackageSize = 1024
	DefaultRecChanSize = 100
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
	Pool *comm.SyncPool
}

//new tcpServer
func NewTcpServer(config *Conifg,handler TcpHandler,protocol Protocol)*TcpServer {
	if config.PackageSize==0{
		config.PackageSize = DefaultPackageSize
	}
	if config.ReceiveSize==0{
		config.ReceiveSize = DefaultRecChanSize
		config.SendSize = DefaultRecChanSize
	}
	return &TcpServer{
		Config:config,
		Handler:handler,
		P:protocol,
		Pool:          comm.NewSyncPool(config.PoolMinSize, config.PoolMaxSize, 2),
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
			ctx := context.Background()
			ctx, cancel := context.WithCancel(ctx)
			connector := NewConn(&conn, tcpServer.Handler, *tcpServer.Config).
				SetSyncPool(tcpServer.Pool).
				SetProtocol(tcpServer.P).
				SetCancleFunc(cancel)//传递 cancel
			//TODO:连接管理
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