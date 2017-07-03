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
}

//new tcpServer
func (updServer *UpdServer) NewTcpServer(port, sSize, rSize int) {
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
	sAddr, err := net.ResolveUDPAddr("udp", "127.0.0.1"+":"+strconv.Itoa(tcpServer.Port))
	//sAddr := tcpServer.Ip + ":" + strconv.Itoa(tcpServer.Port)
	listen, err := net.ListenUDP("udp", sAddr)
	defer listen.Close()
	if err != nil {
		log.Println("udp Server Start Error." + err.Error())
	}
	log.Println(sAddr.String(), "Start Listen.")

	for {
		// data := make([]byte, 1024)
		// n, rAddr, err := listen.ReadFromUDP(data[0:])
		// if err != nil {
		// 	fmt.Print("Recv Conn Error.", err.Error())
		// } else {
		//connector := NewConn(&conn, tcpServer.Handler, *tcpServer.config)

		//新连接 添加到在线列表里面
		// if _, exists := tcpServer.Online[connector.RemoteAddress]; !exists {
		// 	tcpServer.Online[connector.RemoteAddress] = connector
		// }

		// tcpServer.Handler.OnConnect(connector)
		// //go connector.ProcessRecv()
		// go connector.DataHandler()
		// go connector.ProcessSend()
		//}
	}

}
