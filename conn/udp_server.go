package conn

import (
	"log"
	"net"
	"strconv"
)

//server
type UpdServer struct {
	Port     int
	Handler  TcpHandler
	config   *Conifg
	protocol Protocol
	Online   map[string]*Connector
	Conn     *net.UDPConn
}

//new tcpServer
func (updServer *UpdServer) NewUpdServer(port, sSize, rSize int) {
	updServer.Port = port
	updServer.config = &Conifg{
		SendSize:    500,
		ReceiveSize: 500,
	}
	//tcpServer.Conn = make([]Connector, 500)
	updServer.Online = make(map[string]*Connector)
}

//启动udp服务
//http://www.cnblogs.com/gaopeng527/p/6128290.html

func (tcpServer *UpdServer) Start(handler TcpHandler) {
	tcpServer.Handler = handler
	svrAddr := ":" + strconv.Itoa(tcpServer.Port)
	sAddr, err := net.ResolveUDPAddr("udp", svrAddr)
	//sAddr := tcpServer.Ip + ":" + strconv.Itoa(tcpServer.Port)
	conn, err := net.ListenUDP("udp", sAddr)
	tcpServer.Conn = conn
	defer func() {
		conn.Close()
		panic(err)
	}()
	if err != nil {
		log.Println("udp Server Start Error." + err.Error())
		return
	}
	log.Println(sAddr.String(), "udp Start Listen.")
	//go tcpServer.handleClient(conn)
	//connector := NewConn(conn, handler, *tcpServer.config)
	buf := make([]byte, 1024)
	for {
		n, rAddr, err := conn.ReadFromUDP(buf) //rAddr
		if err != nil {
			return
		}
		//fmt.Println("Receive from client", rAddr.String(), string(buf[0:n]))
		//_, err2 := conn.WriteToUDP([]byte("Welcome client!"), rAddr)
		// if err2 != nil {
		// 	return
		// }
		// //rAddr = "" //远程连接处理
		if n > 0 {
			tcpServer.Handler.OnReceive(&Connector{RemoteAddress: rAddr.String()}, &TcpData{buffer: buf[:n]})
		}
	}
}

func (tcpServer *UpdServer) handleClient(conn *net.UDPConn) {
	buf := make([]byte, 1024)
	for {
		n, rAddr, err := conn.ReadFromUDP(buf) //rAddr
		if err != nil {
			return
		}
		//fmt.Println("Receive from client", rAddr.String(), string(buf[0:n]))
		_, err2 := conn.WriteToUDP([]byte("Welcome client!"), rAddr)
		if err2 != nil {
			return
		}
		//rAddr = "" //远程连接处理
		if n > 0 {
			//tcpServer.Handler.OnReceive(&Connector{RemoteAddress: rAddr.String()}, TcpData{buffer: buf[:n]})
		}
	}
}

/*向upd client 发送数据*/
func (tcpServer *UpdServer) SendToClient(udpAdr string, d TcpData) (int, error) {
	addr, err := net.ResolveUDPAddr("udp", udpAdr)
	if err != nil {
		return 0, err
	}
	n, err := tcpServer.Conn.WriteToUDP(d.Bytes(), addr) //转发给远程udp
	return n, err
}
