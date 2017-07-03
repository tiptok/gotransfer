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

func (c *TcpClient) Start(handler TcpHandler) {

	c.Handler = handler
	sAddr, err := net.ResolveTCPAddr("tcp4", c.ServerIp+":"+strconv.Itoa(c.ServerPort))
	//sAddr := tcpServer.Ip + ":" + strconv.Itoa(tcpServer.Port)
	client, err := net.Dial("tcp", sAddr.String())
	if err != nil {
		log.Println("Client Start Error." + err.Error())
		return
	}
	log.Println(sAddr.String(), "Start Client.")

	connector := NewConn(&client, nil)
	c.Handler.OnConnect(connector)
	go connector.ProcessRecv()
	go connector.DataHandler()
	go connector.ProcessSend()
}
