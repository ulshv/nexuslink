package tcp

import (
	"net"
)

type NetConnection interface {
	Connection() net.Conn
}
