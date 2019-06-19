package netserver

import (
	"net"
	"netserver/log"
	"time"
)

type Server struct {
	logger        *log.Logger
	listener      *net.TCPListener
}

func NewServer(aAddr string) *Server {
	log := log.NewLogger()
	addr, err := net.ResolveTCPAddr("tcp4", aAddr)
	if err != nil {
		log.Fatal("resolve tcp addr:%v failed. error:%v", aAddr, err.Error())
		return nil
	}

	listener, err := net.ListenTCP("tcp4", addr)
	if err != nil {
		log.Fatal("resolve tcp addr:%v failed. error:%v", aAddr, err.Error())
		return nil
	}

	return &Server{
		log,
		listener,
	}
}

func (this *Server) Start() {
	this.logger.Info("server started. listen on:%s", this.listener.Addr().String())
	for {
		conn, err := this.listener.AcceptTCP()
		if err != nil {
			this.logger.Warn("accept failed:%s", err.Error())
			time.Sleep(time.Second)
			continue
		}

		this.logger.Info("client:%s connected", conn.RemoteAddr().String())
		userConn := NewUserConn(conn)
		userConn.Start()
	}
}
