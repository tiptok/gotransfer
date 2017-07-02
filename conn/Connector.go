package conn

import (
	"bytes"
	"errors"
	"fmt"
	"net"
)

//连接事件
//接收数据
func (connector *Connector) ProcessRecv() {
	defer func() {
		recover()
		//关闭
	}()

	conn := *(connector.Conn)
	defer conn.Close()
	fmt.Println("New Conn->RemoteAddr:", conn.RemoteAddr())

	for {
		// data := make([]byte, 1024)
		// length, err := conn.Read(data)
		// tcpData := TcpData{buffer: data[:length]}
		tcpData, err := connector.ReadFullData()
		if err != nil && tcpData.Lenght() > 0 {
			connector.RecChan <- tcpData
		}
	}
}

//处理数据
func (connector *Connector) DataHandler() {
	defer func() {
		recover()
		//close
	}()

	for {
		select {
		case p := <-connector.RecChan:
			if !(connector.srv.Handler.OnReceive(connector, p)) {
				return
			}
		}
	}
}

//发送数据
func (connector *Connector) ProcessSend() {
	defer func() {
		recover()
		//close
	}()
	conn := *(connector.Conn)
	for {
		select {
		case p := <-connector.SendChan:
			if len(p.buffer) == 0 {
				return
			}
			if _, err := conn.Write(p.buffer); err != nil {
				return
			}
		}
	}
}

func (connector *Connector) ReadFullData() (TcpData, error) {

	conn := *(connector.Conn)
	defer conn.Close() //关闭

	buf := bytes.NewBuffer([]byte{})
	var tcpData TcpData
	for {
		data := make([]byte, 1024)
		length, err := conn.Read(data)
		if err != nil {
			//fmt.Println(conn.RemoteAddr(), " Error reading.", err.Error())
			return tcpData, err
		}
		if length == 0 {
			if buf.Len() > 0 {
				tcpData = TcpData{buffer: buf.Bytes()}
				return tcpData, nil
			} else {
				return tcpData, errors.New("Read Data None")
			}
		}
		buf.Write(data[:length])
	}
	//return TcpData{buffer: buf.Bytes()}, nil
}

//Connector config
type Conifg struct {
	SendSize    uint32
	ReceiveSize uint32
}

type Connector struct {
	srv           *TcpServer
	Conn          *net.Conn
	handler       TcpHandler
	SendChan      chan TcpData
	RecChan       chan TcpData
	RemoteAddress string
}

func NewConn(tcpconn *net.Conn, srv *TcpServer) *Connector {
	return &Connector{
		srv:           srv,
		Conn:          tcpconn,
		SendChan:      make(chan TcpData, srv.config.SendSize),
		RecChan:       make(chan TcpData, srv.config.ReceiveSize),
		RemoteAddress: (*tcpconn).RemoteAddr().String(),
	}
}
