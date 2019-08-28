package conn

import (
	"context"
	"log"
	"net"
	"strconv"
)

type TcpClient struct {
	ServerPort int
	ServerIp   string
	Handler    TcpHandler
	Config     *Conifg
	P          Protocol
	LocalAdr   string
	Conn       *Connector
}

//new tcpClient
func (tcpClient *TcpClient) NewTcpClient(ip string, port int,handler TcpHandler) {
	tcpClient.ServerIp = ip
	tcpClient.ServerPort = port
	tcpClient.Handler = handler
	tcpClient.Config = &Conifg{
		SendSize:    500,
		ReceiveSize: 500,
		PackageSize: 1024,
	}
}

//启动tcp服务

func (c *TcpClient) Start() bool {
	defer func() {
		if err := recover(); err != nil {
			log.Println("On Recover", err)
			//fmt.Println(err)
		}
	}()
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

	connector := NewConn(&client, c.Handler, *c.Config)
	c.Conn = connector
	c.Handler.OnConnect(connector)
	connector.P = c.P

	ctx := context.Background()
	ctx, cancel := context.WithCancel(ctx)
	connector.cancelFunc = cancel
	go connector.ProcessRecv(ctx)
	go connector.DataHandler(ctx)
	go connector.ProcessSend(ctx)
	return true
}

//重新启动
func (c *TcpClient) ReStart() bool {
	return c.Start()
}
