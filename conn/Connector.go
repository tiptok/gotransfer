package conn

import (
	"bytes"
	"errors"
	"fmt"
	"log"
	"net"
	"sync"
)

//连接事件
//接收数据
func (connector *Connector) ProcessRecv() {
	defer func() {
		MyRecover()
		//connector.Close()
	}()

	conn := *(connector.Conn)
	fmt.Println("New Conn->RemoteAddr:", conn.RemoteAddr())

	for {

		// select {
		// case <-connector.ExitChan:
		// 	return
		// default:
		// }
		data := make([]byte, 1024)
		length, err := conn.Read(data)
		if err != nil {
			//?如何处理关闭
			log.Println(err.Error())
			connector.SendChan <- TcpData{} //向写数据发送一个nil告诉即将关闭
			return
		}
		tcpData := TcpData{buffer: data[:length]}
		if tcpData.Lenght() > 0 {
			connector.RecChan <- tcpData
		}
	}
}

//处理数据
func (connector *Connector) DataHandler() {
	defer func() {
		MyRecover()
		connector.Close()
	}()

	for {
		select {
		// case <-connector.ExitChan:
		// 	return
		case p := <-connector.RecChan:
			if !(connector.handler.OnReceive(connector, p)) {
			}
			//default:
		}
	}
}

//发送数据
func (connector *Connector) ProcessSend() {
	defer func() {
		MyRecover()
	}()
	conn := *(connector.Conn)
	for {
		select {
		// case <-connector.ExitChan:
		// 	return
		case p := <-connector.SendChan:
			if p.buffer == nil {
				//log.Println(err.Error())
				connector.ExitChan <- 1
				//connector.Close()
				return
			}
			if _, err := conn.Write(p.buffer); err != nil {
				return
			}
		}
	}
}

func (connector *Connector) ReadFullData() (TcpData, error) {

	defer func() {
		MyRecover()
		connector.Close()
	}()
	conn := *(connector.Conn)
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

func (c *Connector) Close() {
	c.CloseOnce.Do(func() {
		close(c.SendChan)
		close(c.RecChan)
		(*c.Conn).Close()
		c.handler.OnClose(c)
		c.IsConneted = false
	})
}

//Connector config
type Conifg struct {
	SendSize    uint32
	ReceiveSize uint32
}

type Connector struct {
	//srv           *TcpServer
	Conn          *net.Conn
	handler       TcpHandler
	SendChan      chan TcpData
	RecChan       chan TcpData
	RemoteAddress string
	CloseOnce     sync.Once
	IsConneted    bool
	ExitChan      chan interface{}
}

func NewConn(tcpconn *net.Conn, h TcpHandler, config Conifg) *Connector { //, srv *TcpServer
	return &Connector{
		//srv:           srv,
		Conn:          tcpconn,
		SendChan:      make(chan TcpData, config.SendSize),
		RecChan:       make(chan TcpData, config.ReceiveSize),
		ExitChan:      make(chan interface{}),
		RemoteAddress: (*tcpconn).RemoteAddr().String(),
		handler:       h,
		IsConneted:    true,
	}
}
