package conn

import "github.com/tiptok/gotransfer/comm"

//tcpData
type TcpData struct {
	buffer []byte
	pool *comm.SyncPool
}

func NewTcpData(d []byte) TcpData {
	return TcpData{buffer: d}
}
func (d TcpData) Bytes() []byte {
	return d.buffer
}

func (d TcpData) Lenght() int {
	return len(d.buffer)
}

func(d TcpData)Free(){
	d.pool.Free(d.buffer)
}

type TcpServerBase struct {
	Server TcpServer
}

func SendEntity(e interface{}, c *Connector) (data []byte, err error) {
	data, err = c.P.PacketMsg(e)
	if err == nil {
		if len(data) > 0 {
			send := NewTcpData(data)
			c.SendChan <- send
		} else {
			err = Error_SendDataLen
		}
	}
	return data, err
}
