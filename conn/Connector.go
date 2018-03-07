package conn

import (
	"bytes"
	"context"
	"errors"
	"log"
	"net"
	"sync"
	"time"
)

//连接事件
//接收数据
func (connector *Connector) ProcessRecv(ctx context.Context) {
	defer func() {
		if err := recover(); err != nil {
			log.Println(err)
		}
	}()

	conn := *(connector.Conn)
	//fmt.Println("New Conn->RemoteAddr:", conn.RemoteAddr())

	for {

		select {
		case <-ctx.Done():
			return
		default:
		}
		data := make([]byte, 1024)
		length, err := conn.Read(data)
		if err != nil {
			//?如何处理关闭
			//log.Println("ProcessRecv Error:", err.Error())
			if connector.IsConneted {
				connector.SendChan <- TcpData{} //向写数据发送一个nil告诉即将关闭
			}
			connector.Close()
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
func (connector *Connector) DataHandler(ctx context.Context) {
	defer func() {
		if err := recover(); err != nil {
			log.Println(err)
		}
	}()

	for {
		select {
		//修改注释
		// case <-connector.ExitChan:
		// 	connector.ExitChan <- 1
		// 	return
		case <-ctx.Done():
			return
		case p, IsClose := <-connector.RecChan:
			if !IsClose {
				connector.Close()
				return
			}
			if !(connector.handler.OnReceive(connector, p)) {
			}
		}
	}
}

//发送数据
func (connector *Connector) ProcessSend(ctx context.Context) {
	defer func() {
		if err := recover(); err != nil {
			log.Println(err)
		}
	}()
	conn := *(connector.Conn)
	for {
		select {
		case <-ctx.Done():
			return
		case <-connector.ExitChan:
			return
		case p, IsClose := <-connector.SendChan:
			if !IsClose {
				//log.Println("Connnector Send Chan Closed...")
				connector.Close()
				return
			}
			if p.buffer == nil {
				connector.ExitChan <- 1
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
		if err := recover(); err != nil {
			log.Println(err)
		}
	}()
	conn := *(connector.Conn)
	buf := bytes.NewBuffer([]byte{})
	var tcpData TcpData
	for {
		//data := make([]byte, 1024)
		length, err := buf.ReadFrom(conn)
		if err != nil {
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
		//buf.Write(data[:length])
	}
}

func (c *Connector) Close() {
	c.CloseOnce.Do(func() {
		c.cancelFunc()
		close(c.SendChan)
		close(c.RecChan)
		c.IsConneted = false
		err := (*c.Conn).Close()
		c.handler.OnClose(c)
		if err != nil {
			log.Println("Close Exception:", err)
		}
		c.ExitChan <- 1
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
	//取消
	cancelFunc context.CancelFunc
	/*剩余包数据*/
	Leftbuf *bytes.Buffer
	P       Protocol
}

func NewConn(tcpconn *net.Conn, h TcpHandler, config Conifg) *Connector { //, srv *TcpServer
	c := &Connector{
		//srv:           srv,
		Conn:          tcpconn,
		SendChan:      make(chan TcpData, config.SendSize),
		RecChan:       make(chan TcpData, config.ReceiveSize),
		ExitChan:      make(chan interface{}),
		RemoteAddress: (*tcpconn).RemoteAddr().String(),
		handler:       h,
		IsConneted:    true,
		HeartTime:     time.Now(),
		Leftbuf:       bytes.NewBuffer([]byte{}),
	}
	// c.ctx = context.Background()
	// ctx, cancel := context.WithCancel(c.ctx)
	// c.ctx = ctx
	// c.cancelFunc = cancel
	return c
}

/*写入剩余数据*/
func (c *Connector) WriteLeftData(leftdata []byte) (int, error) {
	return c.Leftbuf.Write(leftdata)
}
