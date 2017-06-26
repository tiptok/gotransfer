package conn

import (
	"fmt"
	"log"
	"net"
	"strconv"
)

//server
type TcpServer struct {
	Ip      string
	Port    int
	Conn    []Connector
	Handler TcpHandler
	config  *Conifg
}

func (tcpServer *TcpServer) NewTcpServer(ip string, port int) {
	tcpServer.Ip = ip
	tcpServer.Port = port
	tcpServer.config = &Conifg{
		SendSize:    500,
		ReceiveSize: 500,
	}
	tcpServer.Conn = make([]Connector, 500)
}

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
			//添加到连接列表里面
			//tcpconn := NewConn(*net.TCPConn(conn), tcpServer)
		}
		//go tcpServer.ProcessRecv(conn)

	}

}

func (tcpServer *TcpServer) ProcessRecv(conn net.Conn) {
	defer conn.Close()
	fmt.Println("New Conn->RemoteAddr:", conn.RemoteAddr())

	for {
		buf := make([]byte, 1024)
		length, err := conn.Read(buf)
		if err != nil {
			fmt.Println(conn.RemoteAddr(), " Error reading.", err.Error())
			return
		}
		tcpServer.RecChan <- TcpData{buffer: buf[:length]}
		fmt.Println("Receive data from client:", string(buf[:length]))
		fmt.Println("Recv Queen Size:", len(tcpServer.RecChan))
	}
}

//server config
type Conifg struct {
	SendSize    uint32
	ReceiveSize uint32
}

type Connector struct {
	srv      *TcpServer
	conn     *net.TCPConn
	handler  TcpHandler
	SendChan chan TcpData
	RecChan  chan TcpData
}

type TcpHandler interface {
	OnConnect(*Connector) bool

	OnReceive(*Connector) bool

	OnClose(*Connector)
}

type TcpData struct {
	buffer []byte
}

func NewConn(tcpconn *net.TCPConn, srv *TcpServer) *Conn {
	return &Connector{
		srv:      srv,
		conn:     tcpconn,
		SendChan: make(chan TcpData, srv.config.SendSize),
		RecChan:  make(chan TcpData, srv.config.ReceiveSize),
	}
}
