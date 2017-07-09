package conn

import (
	"log"
	"net"
	"strconv"
)

type TcpClient struct {
	ServerPort int
	ServerIp   string
	Handler    TcpHandler
	config     *Conifg
	protocol   Protocol
	LocalAdr   string
	Conn       *Connector
}

//new tcpClient
func (tcpClient *TcpClient) NewTcpClient(ip string, port, sSize, rSize int) {
	tcpClient.ServerIp = ip
	tcpClient.ServerPort = port
	tcpClient.config = &Conifg{
		SendSize:    500,
		ReceiveSize: 500,
	}
}

//启动tcp服务

func (c *TcpClient) Start(handler TcpHandler) bool {
	defer func() {
		MyRecover()
	}()
	c.Handler = handler
	sAddr, err := net.ResolveTCPAddr("tcp4", c.ServerIp+":"+strconv.Itoa(c.ServerPort))
	//sAddr := tcpServer.Ip + ":" + strconv.Itoa(tcpServer.Port)
	client, err := net.Dial("tcp", sAddr.String())
	if err != nil {
		log.Println(err.Error())
		return false
	}
	c.LocalAdr = client.LocalAddr().String()
	if err != nil {
		log.Println("Client Start Error." + err.Error())
		return false
	}
	log.Println(sAddr.String(), "Start Client.")

	connector := NewConn(&client, handler, *c.config)
	c.Conn = connector
	c.Handler.OnConnect(connector)
	go connector.ProcessRecv()
	go connector.DataHandler()
	go connector.ProcessSend()
	return true
}
