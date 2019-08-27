package conn

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"log"
	"net"
	"sync"
	"sync/atomic"
	"time"

	"github.com/tiptok/gotransfer/comm"
)

func assertIConnectorImplementation(){
	var _  IConnector= (*Connector)(nil)
}

const (
	G_Recv =1
	G_Send =2
	G_Handle=4
)

//Connector  conn
type Connector struct {
	//srv           *TcpServer
	Conn          *net.Conn
	handler       TcpHandler
	SendChan      chan *TcpData
	RecChan       chan *TcpData
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
	/*对象池*/
	Pool *comm.SyncPool

	stopFlag        int32
	connectedFlag   int32

	gWait    sync.WaitGroup
	gClosed  int32
}

//NewConn new Connector
func NewConn(tcpconn *net.Conn, h TcpHandler, config Conifg) *Connector { //, srv *TcpServer
	c := &Connector{
		//srv:           srv,
		Conn:          tcpconn,
		SendChan:      make(chan *TcpData, config.SendSize),
		RecChan:       make(chan *TcpData, config.ReceiveSize),
		ExitChan:      make(chan interface{},1),
		RemoteAddress: (*tcpconn).RemoteAddr().String(),
		handler:       h,
		IsConneted:    true,
		HeartTime:     time.Now(),
		Config:        config,
		Leftbuf:       bytes.NewBuffer([]byte{}),
		Pool:          comm.NewSyncPool(config.PoolMinSize, config.PoolMaxSize, 2),
	}
	return c
}

func(Connector *Connector)Status()string{
	rspMap:=make(map[string] interface{})
	rspMap["GroutineClosed"]=Connector.gClosed
	data,err:= json.Marshal(rspMap)
	if err!=nil {
		log.Println(err)
		return ""
	}
	return string(data)
}

//连接事件
//接收数据
func (connector *Connector) ProcessRecv(ctx context.Context) {
	defer func() {
		if err := recover(); err != nil {
			log.Println(err)
		}
	}()
	defer func(){
		connector.gWait.Wait()
		connector.handler.OnClose(connector) //所有都执行完以后触发结束
	}()
	connector.gWait.Add(1)
	defer func(){
		atomic.AddInt32(&connector.gClosed,G_Recv)
		connector.gWait.Done()
		connector.Close()
	}()

	conn := *(connector.Conn)
	for {
		select {
		case <-ctx.Done():
			return
		default:
		}
		data := connector.Pool.Alloc(connector.Config.PackageSize)
		length, err := conn.Read(data)
		if err != nil{//io.EOF
			//处理关闭
			if connector.IsConneted {
				connector.SendChan <- &TcpData{singnal:Singnal_Kill} //向写数据发送一个nil告诉即将关闭
			}
			return
		}
		connector.RefreshTime() //刷新最新时间
		err = connector.parsePart(data[:length])
	}
}

//DataHandler 处理数据
func (connector *Connector) DataHandler(ctx context.Context) {
	defer func() {
		if err := recover(); err != nil {
			log.Println(err)
		}
	}()
	connector.gWait.Add(1)
	defer func(){
		atomic.AddInt32(&connector.gClosed,G_Handle)
		connector.gWait.Done()
	} ()

	for {
		select {
		case <-ctx.Done():
			return
		case p, IsClose := <-connector.RecChan:
			if !IsClose {
				//connector.Close()
				return
			}
			if !(connector.handler.OnReceive(connector, p)) {
				p.Free()
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
	connector.gWait.Add(1)
	defer func(){
		atomic.AddInt32(&connector.gClosed,G_Send)
		connector.gWait.Done()
	}()
	conn := *(connector.Conn)
	for {
		select {
		case <-ctx.Done():
			return
		case p, IsClose := <-connector.SendChan:
			if !IsClose {
				return
			}
			if p.buffer == nil && p.singnal==Singnal_Kill{
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
		atomic.CompareAndSwapInt32(&c.stopFlag,0,1)
		(*c.Conn).Close()
		atomic.CompareAndSwapInt32(&c.connectedFlag,0,1)
		//TODO:exitchan
		//c.ExitChan <- 1
	})
}

//公共方法
func(c *Connector)SetSyncPool(pool *comm.SyncPool)*Connector{
	if pool==nil{
		pool=comm.NewSyncPool(c.Config.PoolMinSize, c.Config.PoolMaxSize, 2)
	}
	c.Pool = pool
	return c
}
func(c *Connector)SetProtocol(p  Protocol)*Connector{
	c.P = p
	return c
}
func(c *Connector)SetCancleFunc(cancel  context.CancelFunc)*Connector{
	c.cancelFunc= cancel
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

func(connector *Connector)Connected()bool{
	return connector.connectedFlag==1
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
		packdata, leftdata, e := connector.P.ParseMsg(connector.Leftbuf.Bytes(), connector)
		connector.Leftbuf.Reset() //LeftBuf置空
		if e != nil {
			err =e
			return
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
		if leftdata != nil && len(leftdata) > 0 {
			_, err = connector.WriteLeftData(leftdata)
		}
		return
	}
	err = connector.parseToEntity(data)
	return
}

//parseToEntity 解析数据到实体
func (connector *Connector) parseToEntity(data []byte) (err error) {
	tcpData := &TcpData{buffer: data,pool:connector.Pool}
	connector.RecChan <- tcpData //发送给接收
	//connector.handler.OnReceive(connector, p)
	return err
}

//ParseToEntity  parse data to entity
func (connector *Connector) ParseToEntity(data []byte) (entity interface{}, err error) {
	if connector.P == nil {
		err = errors.New("ParseToEntity Error:未定义协议")
		return nil, err
	}
	defer connector.Pool.Free(data)       //回收到对象池
	entity, err = connector.P.Parse(data) //entity
	return entity, err
}
