package conn

type Packet interface {
	Serialize() []byte
}

/*
	打包分包协议
*/
type Protocol interface {
	/*
		打包分包处理
		obj 打包的数据结构
	*/
	PacketMsg(obj interface{}) (data []byte, err error)
	/*
		打包数据体
		obj 数据体
	*/
	Packet(obj interface{}) (packdata []byte, err error)
	/*
		分包处理
		packdata 解析出一个完整包
		leftdata 解析剩余报文的数据
		err 	 分包错误
	*/
	ParseMsg(data []byte, c *Connector) (packdata [][]byte, leftdata []byte, err error)
	/*
		解析数据
		obj 解析出对应得数据结构
		err 解析出错
	*/
	Parse(packdata []byte) (obj interface{}, err error)
}

// /*协议解析配置*/
// type ProtocolBaseConfig struct {
// }

type DefaultProtocol struct {
}

type DefaultTcpData struct {
	/*头标识*/
	BEGIN byte

	/*消息头*/
	/*消息类型编号*/
	MsgTypeId int16
	/*序号*/
	Id int16
	/*数据体长度*/
	Length int32
	/*分包属性*/
	PackagesProperty int32

	//Body

	/*校验字*/
	Valid byte
	/*尾标识*/
	END byte
}
