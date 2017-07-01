package conn

import (
	"net"
)

type Packet interface {
	Serialize() []byte
}

type Protocol interface {
	Parse(conn *net.Conn)
	Packet() []byte
}
