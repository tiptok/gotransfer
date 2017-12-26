package main

import (
	"fmt"

	"time"

	"log"

	"net"

	comm "github.com/tiptok/gotransfer/comm"
	conn "github.com/tiptok/gotransfer/conn"
)

var Alive bool
var srv conn.TcpServer
var srvUdp conn.UpdServer
var exit chan int
var param comm.Param
var tcpTransferList map[string](map[string]conn.TcpClient)
var udpTransList map[string](*conn.UpdClient)

var udpFrom map[string]([]*conn.UpdClient)
var udpTo map[string]*conn.Connector

//timer work 定时清理
var timerClear *time.Timer

//timer work 重连
var timerReconn *time.Timer

func main() {
	//Test()
	Alive = true
	tcpTransferList = make(map[string](map[string]conn.TcpClient))
	udpTransList = make(map[string](*conn.UpdClient))
	exit = make(chan int, 1)
	udpFrom = make(map[string]([]*conn.UpdClient))
	udpTo = make(map[string]*conn.Connector)

	//F:/App/Go/WorkSpace/src/github.com/tiptok/gotransfer/
	err := comm.ReadConfig(&param, "param.json")
	if err != nil {
		fmt.Println(err.Error())
	}
	fmt.Println(param)
	/*定时清理*/
	timerClear = time.NewTimer(time.Second * 60)
	go TimerClear()
	/*重连*/
	timerReconn = time.NewTimer(time.Second * 10)
	/*日志模块*/
	comm.Config("E:\\Logs\\goTransferLog", 3)

	//go TimerReconn()
	//启动tcp服务
	go func() {
		srv.NewTcpServer(param.Tcpsrv.Port, param.Tcpsrv.Sendsize, param.Tcpsrv.Recvsize)
		srv.Start(&transferSvrHandler{})
	}()

	//启动upd服务
	// go func() {
	// 	srvUdp.NewUpdServer(param.Tcpsrv.Port, param.Tcpsrv.Sendsize, param.Tcpsrv.Recvsize)
	// 	srvUdp.Start(&transferUdpSvrHandler{})
	// }()

	<-exit
	fmt.Println("App Stop.")
}

/*
连接事件
tpc server 使用
*/
type transferSvrHandler struct {
}

func (trans *transferSvrHandler) OnConnect(c *conn.Connector) bool {
	defer func() {
		//conn.MyRecover()
	}()
	log.Println(c.RemoteAddress, "On OnConnect.")
	if param.NeedTransfer {
		for _, srv := range param.Transfers {
			var tcpClient conn.TcpClient
			tcpClient.NewTcpClient(srv.Ip, srv.Port, 500, 500)
			if !tcpClient.Start(&TransferClientHandler{}) {
				continue //client 启动失败
			}
			//tcpAddress +"-1" as tcp key
			//已存在
			key := c.RemoteAddress + "-1"
			if _, ok := tcpTransferList[key]; !ok {
				tcpTransferList[key] = make(map[string]conn.TcpClient)
			}
			tcpTransferList[key][tcpClient.LocalAdr] = tcpClient
		}
	}
	return true
}
func (trans *transferSvrHandler) OnReceive(c *conn.Connector, d conn.TcpData) bool {
	fmt.Println(TimeFormat(time.Now()), c.RemoteAddress, " OnReceive Data.", string(d.Bytes()))
	comm.Debugf(" OnReceive Data.", string(d.Bytes()))
	defer func() {
		conn.MyRecover()
	}()
	// sKey := c.RemoteAddress + "-1"
	// if _, ok := tcpTransferList[sKey]; ok {
	// 	//分发数据
	// 	for _, value := range tcpTransferList[sKey] {
	// 		(*value.Conn).SendChan <- d
	// 	}
	// }
	//echo

	if c.IsConneted {
		c.SendChan <- d
	}
	return true
}
func (trans *transferSvrHandler) OnClose(c *conn.Connector) {
	defer func() {
		//conn.MyRecover()
	}()
	//从列表移除
	//DoColse(c)

	log.Println(c.RemoteAddress, "On Close.")
}
func DoColse(c *conn.Connector) {
	sKey := c.RemoteAddress
	tempTrans := ""
	if _, ok := srv.Online[sKey]; ok {
		delete(srv.Online, sKey) /*从在线列表移除*/
		if _, ok2 := tcpTransferList[sKey+"-1"]; ok2 {
			for _, v := range tcpTransferList[sKey+"-1"] {
				(*v.Conn).Close() /*关闭分发的连接*/
				tempTrans += (*v.Conn).RemoteAddress + " "
			}
			delete(tcpTransferList, sKey+"-1") /*从分发列表移除*/
			log.Println(sKey, "Remote Do cloese:")
		}
	}
}

/*
连接事件
tpc client 使用
*/
type TransferClientHandler struct {
}

func (trans *TransferClientHandler) OnConnect(c *conn.Connector) bool {
	return true
}
func (trans *TransferClientHandler) OnReceive(c *conn.Connector, d conn.TcpData) bool {
	defer func() {
		conn.MyRecover()
	}()
	sKey := SearchDstClient((*c.Conn).LocalAddr().String())
	if _, ok := srv.Online[sKey]; ok {
		//fmt.Println(TimeFormat(time.Now()), (*c.Conn).LocalAddr().String(), "Send To->", sKey, ":", string(d.Bytes()))
		//下方数据到 目标客户端
		srv.Online[sKey].SendChan <- d
	}
	return true
}
func (trans *TransferClientHandler) OnClose(c *conn.Connector) {
	//从列表移除
}

/*
连接事件
udp server 使用
*/
type transferUdpSvrHandler struct {
}

func (trans *transferUdpSvrHandler) OnConnect(c *conn.Connector) bool {
	address := c.RemoteAddress
	if _, ok := udpFrom[address]; !ok {
		udpFrom[address] = make([]*conn.UpdClient, len(param.Transfers))
	}
	/*发起一个新的udp连接*/
	for index, srv := range param.Transfers {
		var new conn.UpdClient
		new.NewUpdClient(srv.Ip, srv.Port, 500, 500)
		if !new.Start(&transferUpdClinetHandler{}) {
			continue //client 启动失败
		}
		udpFrom[address][index] = &new

		//添加新的映射
		if _, ok := udpTo[new.LocalAdr]; !ok {
			udpTo[new.LocalAdr] = c
		}
	}

	return true
}
func (trans *transferUdpSvrHandler) OnReceive(c *conn.Connector, d conn.TcpData) bool {
	//fmt.Println(TimeFormat(time.Now()), c.RemoteAddress, "Udp OnReceive Data.", string(d.Bytes()))
	defer func() {
		conn.MyRecover()
	}()
	//连接
	address := c.RemoteAddress
	if _, ok := udpFrom[address]; !ok {
		udpFrom[address] = make([]*conn.UpdClient, len(param.Transfers))
		log.Println("当前Udp连接数:", 2*len(udpFrom))
		/*发起一个新的udp连接*/
		for index, srv := range param.Transfers {
			var new conn.UpdClient
			new.NewUpdClient(srv.Ip, srv.Port, 500, 500)
			if !new.Start(&transferUpdClinetHandler{}) {
				continue //client 启动失败
			}
			udpFrom[address][index] = &new

			//添加新的映射
			if _, ok := udpTo[new.LocalAdr]; !ok {
				udpTo[new.LocalAdr] = c
			}
			udpTo[new.LocalAdr].HeartTime = time.Now()
		}
	}

	//分发
	if _, ok := udpFrom[c.RemoteAddress]; ok {
		bNeedRemove := false
		// if to,ok :=udpTo[udpFrom[c.RemoteAddress].LocalAdr]{
		// 	to.HeartTime = time.Now
		// }
		for _, value := range udpFrom[c.RemoteAddress] {
			if (*value.Conn).IsConneted {
				value.Conn.SendChan <- d
			} else {
				bNeedRemove = true
			}
		}
		//删除掉线的 从新连接
		if bNeedRemove {
			for _, value := range udpFrom[c.RemoteAddress] {
				(*value.Conn).Close()
			}
		}
		if bNeedRemove {
			delete(udpFrom, c.RemoteAddress)
		}
	}
	return true
}
func (trans *transferUdpSvrHandler) OnClose(c *conn.Connector) {
	//从列表移除
}

/*
连接事件
tpc client 使用
*/
type transferUpdClinetHandler struct {
}

func (trans *transferUpdClinetHandler) OnConnect(c *conn.Connector) bool {
	return true
}
func (trans *transferUpdClinetHandler) OnReceive(c *conn.Connector, d conn.TcpData) bool {
	defer func() {
		conn.MyRecover()
	}()
	sAddr := (*c.Conn).LocalAddr()
	slocal := sAddr.String()
	if _, ok := udpTo[slocal]; ok {
		c := udpTo[slocal]
		addr, err := net.ResolveUDPAddr("udp", c.RemoteAddress)
		if err != nil {
			log.Println(err.Error())
			return false
		}
		_, err = srvUdp.Conn.WriteToUDP(d.Bytes(), addr) //转发给远程udp
		if err != nil {
			log.Println(err.Error())
		}
		//log.Println(slocal, "Send To->", c.RemoteAddress, " ", conn.ToHex(d.Bytes()))
		// if c.IsConneted {
		// 	c.SendChan <- d
		// }
	}
	return true
}
func (trans *transferUpdClinetHandler) OnClose(c *conn.Connector) {
	defer func() {
		conn.MyRecover()
	}()
	//从列表移除
	log.Println(c.RemoteAddress, "On Close...")
}

//公用方法
//寻找目标
func SearchDstClient(skey string) (dst string) {
	for key1, _ := range tcpTransferList {
		for key2, _ := range tcpTransferList[key1] {
			if key2 == skey {
				sub := []byte(key1)
				length := len(key1) - 2
				dst = string(sub[:length]) //去掉末尾 -1
			}
		}
	}
	return dst
}

func TimeFormat(t time.Time) string {
	tFormat := t.Format("2006-01-02 15:04:05")
	return tFormat
}
