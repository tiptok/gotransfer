package conn

import (
	"log"
	"net"
	"strconv"
)

type UpdClient struct {
	ServerPort int
	ServerIp   string
	Handler    TcpHandler
	config     *Conifg
	protocol   Protocol
	LocalAdr   string
	Conn       *Connector
}

//new tcpClient
func (udpClient *UpdClient) NewTcpClient(ip string, port, sSize, rSize int) {
	udpClient.ServerIp = ip
	udpClient.ServerPort = port
	udpClient.config = &Conifg{
		SendSize:    500,
		ReceiveSize: 500,
	}
}

//启动tcp服务

func (c *UpdClient) Start(handler TcpHandler) bool {

	c.Handler = handler
	sAddr, err := net.ResolveUDPAddr("udp", c.ServerIp+":"+strconv.Itoa(c.ServerPort))
	//sAddr := tcpServer.Ip + ":" + strconv.Itoa(tcpServer.Port)
	client, err := net.Dial("udp", sAddr.String())
	c.LocalAdr = client.LocalAddr().String()
	if err != nil {
		log.Println("udp Client Start Error." + err.Error())
		return false
	}
	log.Println(sAddr.String(), "udp Client Start Client.")

	connector := NewConn(&client, handler, *c.config)
	c.Conn = connector
	c.Handler.OnConnect(connector)
	go connector.ProcessRecv()
	go connector.DataHandler()
	go connector.ProcessSend()
	return true
}
