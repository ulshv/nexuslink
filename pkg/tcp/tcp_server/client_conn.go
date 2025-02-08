package tcp_server

import (
	"io"
	"net"

	"github.com/ulshv/nexuslink/internal/logger"
)

type clientConnection struct {
	conn     net.Conn
	closedCh chan struct{}
}

func newClientConnection(conn net.Conn) *clientConnection {
	return &clientConnection{
		conn:     conn,
		closedCh: make(chan struct{}),
	}
}

func (c *clientConnection) Write(b []byte) (int, error) {
	logger.DefaultLogger.Debug("Writing to client connection", "addr", c.conn.RemoteAddr(), "bytes", len(b))
	n, err := c.conn.Write(b)
	if err != nil {
		logger.DefaultLogger.Error("Failed to write to client connection", "addr", c.conn.RemoteAddr(), "error", err)
		if err == io.ErrClosedPipe {
			logger.DefaultLogger.Debug("Client connection closed", "addr", c.conn.RemoteAddr())
			c.closedCh <- struct{}{}
			return 0, err
		}
		return 0, err
	}
	return n, nil
}
