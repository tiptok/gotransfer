package conn

import (
	"bytes"
	"errors"
	"net"
	"sync"
	"time"
)

//连接事件
//接收数据
func (connector *Connector) ProcessRecv() {
	defer func() {
		MyRecover()
	}()

	conn := *(connector.Conn)
	//fmt.Println("New Conn->RemoteAddr:", conn.RemoteAddr())

	for {

		// select {
		// case <-connector.ExitChan:
		// 	connector.ExitChan <- 1
		// 	return
		// default:
		// }
		data := make([]byte, 1024)
		length, err := conn.Read(data)
		if err != nil {
			//?如何处理关闭
			//log.Println("ProcessRecv Error:", err.Error())
			if connector.IsConneted {
				connector.SendChan <- TcpData{} //向写数据发送一个nil告诉即将关闭
			}
			return
		}
		connector.HeartTime = time.Now() //刷新最新时间
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
		if connector != nil {
			connector.Close()
		}
	}()

	for {
		select {
		//修改注释
		// case <-connector.ExitChan:
		// 	connector.ExitChan <- 1
		// 	return
		case p, IsClose := <-connector.RecChan:
			if !IsClose {
				connector.Close()
				return
			}
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
		case <-connector.ExitChan:
			//connector.ExitChan <- 1
			return
		case p, IsClose := <-connector.SendChan:
			if !IsClose {
				//log.Println("Connnector Send Chan Closed...")
				//connector.Close()
				return
			}
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
		c.ExitChan <- 1
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
	HeartTime     time.Time
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
		HeartTime:     time.Now(),
	}
}
