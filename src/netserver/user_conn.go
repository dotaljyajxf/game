package netserver

import (
	"net"
)

type UserConn struct {
	conn *net.TCPConn
}
