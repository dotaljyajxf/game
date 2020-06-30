package gateway

import "net"

type UserManager struct {
	conn *net.TCPConn
	uid  uint64
}
