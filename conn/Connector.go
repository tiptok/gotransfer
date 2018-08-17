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
	rb := new(bytes.Buffer)
	data := make([]byte, connector.Config.PackageSize)
	for {

		select {
		case <-ctx.Done():
			return
		default:
		}
		length, err := conn.Read(data)
		//length, err := rb.ReadFrom(conn)
		if err != nil {
			//?如何处理关闭
			if connector.IsConneted {
				connector.SendChan <- TcpData{} //向写数据发送一个nil告诉即将关闭
			}
			connector.Close()
			return
		}
		connector.RefreshTime() //刷新最新时间
		err = connector.parsePart(data[:length])
		rb.Reset()
	}
}

//DataHandler 处理数据
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

//ProcessSend 发送数据
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

//ReadFullData 读取数据
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

//Close connector on close
func (c *Connector) Close() {
	c.CloseOnce.Do(func() {
		c.cancelFunc()
		close(c.SendChan)
		close(c.RecChan)
		c.IsConneted = false
		c.handler.OnClose(c)
		(*c.Conn).Close()
		c.ExitChan <- 1
	})
}

//Conifg connector config
type Conifg struct {
	SendSize        uint32
	ReceiveSize     uint32
	PackageSize     int
	IsParsePartMsg  bool //是否对接收到的数据进行分包处理
	IsParseToEntity bool //是否解析为实体
}

//Connector  conn
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
	Config        Conifg
	//取消
	cancelFunc context.CancelFunc
	/*剩余包数据*/
	Leftbuf *bytes.Buffer
	P       Protocol
}

//NewConn new Connector
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
		Config:        config,
		Leftbuf:       bytes.NewBuffer([]byte{}),
	}
	// c.ctx = context.Background()
	// ctx, cancel := context.WithCancel(c.ctx)
	// c.ctx = ctx
	// c.cancelFunc = cancel
	return c
}

//LocalAddr	连接的本地Address
func (c *Connector) LocalAddr() (net.Addr, error) {
	if !c.IsConneted {
		return nil, errors.New("Connector Closed.")
	}
	return (*c.Conn).LocalAddr(), nil
}

/*写入剩余数据*/
func (c *Connector) WriteLeftData(leftdata []byte) (int, error) {
	return c.Leftbuf.Write(leftdata)
}

//RefreshTime 刷新心跳时间
func (connector *Connector) RefreshTime() {
	connector.HeartTime = time.Now() //刷新最新时间
}

//parsePart 解析分包数据
func (connector *Connector) parsePart(data []byte) (err error) {
	isParsePart := connector.Config.IsParsePartMsg
	if connector.P == nil {
		isParsePart = false //未定义协议,不进行分包处理
	}
	//需要解析分包
	if isParsePart {
		//将新接收到的数据追加的上一次解析剩余的
		connector.WriteLeftData(data)
		//AddLeft := connector.Leftbuf
		packdata, leftdata, err := connector.P.ParseMsg(connector.Leftbuf.Bytes(), connector)
		connector.Leftbuf.Reset() //LeftBuf置空
		if err != nil {
			return err
		}

		/*解析完整包*/
		if len(packdata) > 0 {
			for i := 0; i < len(packdata); i++ {
				if len(packdata[i]) <= 0 {
					continue
				}
				err = connector.parseToEntity(packdata[i]) //发送给接收
			}
		}
		// else {
		// 	err = errors.New("parsePart 未解析出数据包")
		// }
		if leftdata != nil && len(leftdata) > 0 {
			_, err = connector.WriteLeftData(leftdata)
			//log.Println("Left Data:", comm.BinaryHelper.ToBCDString(leftdata, 0, int32(len(leftdata))), connector.Leftbuf.Len())
		}
	} else {
		err = connector.parseToEntity(data)
	}
	return err
}

//parseToEntity 解析数据到实体
func (connector *Connector) parseToEntity(data []byte) (err error) {
	// isParseToEntity := connector.Config.IsParseToEntity
	// if connector.P == nil {
	// 	isParseToEntity = false //未定义协议,不进行分包处理
	// }
	//是否解析为实体
	// if isParseToEntity {
	// 	_, err = connector.P.Parse(data) //entity
	// 	//传递 entity
	// } else {
	// 	tcpData := TcpData{buffer: data}
	// 	connector.RecChan <- tcpData //发送给接收
	// }
	tcpData := TcpData{buffer: data}
	connector.RecChan <- tcpData //发送给接收
	return err
}

//ParseToEntity  parse data to entity
func (connector *Connector) ParseToEntity(data []byte) (entity interface{}, err error) {
	if connector.P == nil {
		err = errors.New("ParseToEntity Error:未定义协议")
		return nil, err
	}
	entity, err = connector.P.Parse(data) //entity
	return entity, err
}
