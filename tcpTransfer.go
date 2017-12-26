package main

import "github.com/tiptok/gotransfer/comm"

type tcpTransfer struct {
	TransContext comm.DataContext
}

/*
连接事件
tpc server 使用
*/
// type transferSvrHandler struct {
// }

// func (trans *transferSvrHandler) OnConnect(c *conn.Connector) bool {
// 	defer func() {
// 		conn.MyRecover()
// 	}()
// 	if param.NeedTransfer {
// 		for _, srv := range param.Transfers {
// 			var tcpClient conn.TcpClient
// 			tcpClient.NewTcpClient(srv.Ip, srv.Port, 500, 500)
// 			if !tcpClient.Start(&TransferClientHandler{}) {
// 				continue //client 启动失败
// 			}
// 			//tcpAddress +"-1" as tcp key
// 			//已存在
// 			key := c.RemoteAddress
// 			if _, ok := tcpTransferList[key]; !ok {
// 				tcpTransferList[key] = make(map[string]conn.TcpClient)
// 			}
// 			tcpTransferList[key][tcpClient.LocalAdr] = tcpClient
// 		}
// 	}

// 	return true
// }
// func (trans *transferSvrHandler) OnReceive(c *conn.Connector, d conn.TcpData) bool {
// 	//fmt.Println(TimeFormat(time.Now()), c.RemoteAddress, " OnReceive Data.", string(d.Bytes()))
// 	defer func() {
// 		conn.MyRecover()
// 	}()
// 	sKey := c.RemoteAddress + "-1"
// 	if _, ok := tcpTransferList[sKey]; ok {
// 		//分发数据
// 		for _, value := range tcpTransferList[sKey] {
// 			(*value.Conn).SendChan <- d
// 		}
// 	}
// 	//echo
// 	// if c.IsConneted {
// 	// 	c.SendChan <- d
// 	// }
// 	return true
// }
// func (trans *transferSvrHandler) OnClose(c *conn.Connector) {
// 	defer func() {
// 		conn.MyRecover()
// 	}()
// 	//从列表移除
// 	DoColse(c)
// 	//log.Println((*c.Conn).RemoteAddr(), "On Close.")
// 	log.Println(c.RemoteAddress, "On Close.")
// }
// func DoColse(c *conn.Connector) {
// 	sKey := c.RemoteAddress
// 	tempTrans := ""
// 	if _, ok := srv.Online[sKey]; ok {
// 		delete(srv.Online, sKey) /*从在线列表移除*/
// 		if _, ok2 := tcpTransferList[sKey+"-1"]; ok2 {
// 			for _, v := range tcpTransferList[sKey+"-1"] {
// 				(*v.Conn).Close() /*关闭分发的连接*/
// 				tempTrans += (*v.Conn).RemoteAddress + " "
// 			}
// 			delete(tcpTransferList, sKey+"-1") /*从分发列表移除*/
// 			log.Println(sKey, "Remote Do cloese:")
// 		}
// 	}
// }
