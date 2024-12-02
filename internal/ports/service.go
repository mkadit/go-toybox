package ports

import "net"

type TCPHandler interface {
	HandleConnection(conn net.Conn)
}
