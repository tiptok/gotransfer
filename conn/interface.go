package conn

import "context"

/*统一Server接口*/
type ITcpServer interface {
	Start() error
	Stop() error
}

//连接
type IConnector interface {
	ProcessRecv(ctx context.Context)
	DataHandler(ctx context.Context)
	ProcessSend(ctx context.Context)
}


//回调处理事件
type TcpHandler interface {
	OnConnect(*Connector) bool
	OnReceive(*Connector, *TcpData) bool
	OnClose(*Connector)
}

/*
	协议解析
*/
type Protocol interface {
	PacketMsg(obj interface{}) (data []byte, err error)/* 打包分包处理  obj 打包的数据结构*/
	Packet(obj interface{}) (packdata []byte, err error)/* 打包数据体   obj 数据体*/
	ParseMsg(data []byte, c *Connector) (packdata [][]byte, leftdata []byte, err error) /*   分包处理  packdata 解析出一个完整包；  leftdata 解析剩余报文的数据；  err 分包错误*/
	Parse(packdata []byte) (obj interface{}, err error)/* 解析数据  obj 解析出对应得数据结构  err 解析出错 */
}