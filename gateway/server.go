package gateway

import (
	log "github.com/sirupsen/logrus"
	"net"
	"time"
)

type Server struct {
	listener *net.TCPListener
}

func NewServer(aAddr string) *Server {
	addr, err := net.ResolveTCPAddr("tcp4", aAddr)
	if err != nil {
		log.Fatalf("resolve tcp addr:%v failed. error:%v", aAddr, err.Error())
		return nil
	}

	listener, err := net.ListenTCP("tcp4", addr)
	if err != nil {
		log.Fatalf("resolve tcp addr:%v failed. error:%v", aAddr, err.Error())
		return nil
	}

	return &Server{
		listener,
	}
}

func (this *Server) Start() {
	log.Infof("server started. listen on:%s", this.listener.Addr().String())
	for {
		conn, err := this.listener.AcceptTCP()
		if err != nil {
			log.Infof("accept failed:%s", err.Error())
			time.Sleep(time.Second)
			continue
		}

		log.Infof("client:%s connected", conn.RemoteAddr().String())
		userConn := NewUserConn(conn)
		userConn.Start()
	}
}
