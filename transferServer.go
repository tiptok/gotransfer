package main

import (
	"fmt"

	"time"

	comm "github.com/tiptok/gotransfer/comm"
	conn "github.com/tiptok/gotransfer/conn"
)

var srv conn.TcpServer
var srvUdp conn.UpdServer
var exit chan int
var param comm.Param
var tcpTransferList map[string](map[string]*conn.Connector)
var udpTransList map[string](*conn.UpdClient)

func main() {
	tcpTransferList = make(map[string](map[string]*conn.Connector))
	udpTransList = make(map[string](*conn.UpdClient))
	exit = make(chan int, 1)
	//F:/App/Go/WorkSpace/src/github.com/tiptok/gotransfer/
	err := comm.ReadConfig(&param, "param.json")
	if err != nil {
		fmt.Println(err.Error())
	}
	fmt.Println(param)
	//启动tcp服务
	// go func() {
	// 	srv.NewTcpServer(param.Tcpsrv.Port, param.Tcpsrv.Sendsize, param.Tcpsrv.Recvsize)
	// 	srv.Start(&transferSvrHandler{})
	// }()

	//启动upd服务
	go func() {
		srvUdp.NewUpdServer(param.Tcpsrv.Port, param.Tcpsrv.Sendsize, param.Tcpsrv.Recvsize)
		srvUdp.Start(&transferUdpSvrHandler{})
	}()

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
	for _, srv := range param.Transfers {
		var tcpClient conn.TcpClient
		tcpClient.NewTcpClient(srv.Ip, srv.Port, 500, 500)
		if !tcpClient.Start(&transferClinetHandler{}) {
			continue //client 启动失败
		}
		//tcpAddress +"-1" as tcp key
		//已存在
		key := c.RemoteAddress + "-1"
		if _, ok := tcpTransferList[key]; !ok {
			tcpTransferList[key] = make(map[string]*conn.Connector)
		}
		tcpTransferList[key][tcpClient.LocalAdr] = tcpClient.Conn
	}
	return true
}
func (trans *transferSvrHandler) OnReceive(c *conn.Connector, d conn.TcpData) bool {
	//fmt.Println(TimeFormat(time.Now()), c.RemoteAddress, " OnReceive Data.", string(d.Bytes()))

	sKey := c.RemoteAddress + "-1"
	if _, ok := tcpTransferList[sKey]; ok {
		//分发数据
		for _, value := range tcpTransferList[sKey] {
			value.SendChan <- d
		}
	}
	return true
}
func (trans *transferSvrHandler) OnClose(c *conn.Connector) {
	//从列表移除
}

/*
连接事件
tpc client 使用
*/
type transferClinetHandler struct {
}

func (trans *transferClinetHandler) OnConnect(c *conn.Connector) bool {
	return true
}
func (trans *transferClinetHandler) OnReceive(c *conn.Connector, d conn.TcpData) bool {
	sKey := SearchDstClient((*c.Conn).LocalAddr().String())
	if _, ok := srv.Online[sKey]; ok {
		//fmt.Println(TimeFormat(time.Now()), (*c.Conn).LocalAddr().String(), "Send To->", sKey, ":", string(d.Bytes()))
		//下方数据到 目标客户端
		srv.Online[sKey].SendChan <- d
	}
	return true
}
func (trans *transferClinetHandler) OnClose(c *conn.Connector) {
	//从列表移除
}

/*
连接事件
udp server 使用
*/
type transferUdpSvrHandler struct {
}

func (trans *transferUdpSvrHandler) OnConnect(c *conn.Connector) bool {
	// for _, srv := range param.Transfers {
	// 	var tcpClient conn.UpdClient
	// 	tcpClient.NewTcpClient(srv.Ip, srv.Port, 500, 500)
	// 	if !tcpClient.Start(&transferUpdClinetHandler{}) {
	// 		continue //client 启动失败
	// 	}
	// 	//tcpAddress +"-1" as tcp key
	// 	//已存在
	// 	key := c.RemoteAddress + "-2"
	// 	if _, ok := tcpTransferList[key]; !ok {
	// 		tcpTransferList[key] = make(map[string]*conn.Connector)
	// 	}
	// 	tcpTransferList[key][tcpClient.LocalAdr] = tcpClient.Conn
	// }
	return true
}
func (trans *transferUdpSvrHandler) OnReceive(c *conn.Connector, d conn.TcpData) bool {
	//fmt.Println(TimeFormat(time.Now()), c.RemoteAddress, "Udp OnReceive Data.", string(d.Bytes()))

	//创建客户端
	if len(udpTransList) == 0 {
		for _, srv := range param.Transfers {
			var c conn.UpdClient
			c.NewUpdClient(srv.Ip, srv.Port, 500, 500)
			if !c.Start(&transferUpdClinetHandler{}) {
				continue //client 启动失败
			}
			key := c.LocalAdr
			if _, ok := udpTransList[key]; !ok {
				udpTransList[key] = &c
			}
		}
	}

	//sKey := c.RemoteAddress + "-2"
	for _, value := range udpTransList {
		value.Conn.SendChan <- d
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
	// sKey := SearchDstClient((*c.Conn).LocalAddr().String())
	// if _, ok := srvUdp.Online[sKey]; ok {
	// 	fmt.Println(time.Now(), (*c.Conn).LocalAddr().String(), "Send To->", sKey, ":", string(d.Bytes()))
	// 	//下方数据到 目标客户端
	// 	srv.Online[sKey].SendChan <- d
	// }
	return true
}
func (trans *transferUpdClinetHandler) OnClose(c *conn.Connector) {
	//从列表移除
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
