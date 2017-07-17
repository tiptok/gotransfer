package comm

type Param struct {
	Db           Dbconfig   `json:"DB"`
	Tcpsrv       Tcpconfig  `json:"Server"`
	Transfers    []Transfer `json:"TransferTo"`
	NeedTransfer bool       `json:"NeedTransfer"`
	HasUDP       bool       `json:"HasUDP"`
}
type Dbconfig struct {
	Userid string `json:"UserId"`
	Pwd    string `json:"Pwd"`
}
type Tcpconfig struct {
	Port     int `json:"Port"`
	Recvsize int `json:"RecvSize"`
	Sendsize int `json:"SendSize"`
}
type Transfer struct {
	Ip   string `json:"Ip"`
	Port int    `json:"Port"`
}
