package tcp_conn

import (
	"net"
)

type NetConnection interface {
	Connection() net.Conn
}
