package conn

//Conifg connector config
type Conifg struct {
	ListenPort      int //监听端口
	SendSize        uint32
	ReceiveSize     uint32
	PackageSize     int  //每次接收读取数据包大小
	IsParsePartMsg  bool //是否对接收到的数据进行分包处理
	IsParseToEntity bool //是否解析为实体
	PackageMinSize  int  //数据包最小容量  单位 B
	PackageMaxSize  int  //数据包最大容量  单位 B
}


