package conn

type Packet interface {
	Serialize() []byte
}

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
