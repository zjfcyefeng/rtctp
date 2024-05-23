package server

import (
	"fmt"
	"net"
	"time"

	getty "github.com/apache/dubbo-getty"
	gxnet "github.com/dubbogo/gost/net"
	gxsync "github.com/dubbogo/gost/sync"
	"github.com/zjfcyefeng/rtctp/internal/codec"
	"github.com/zjfcyefeng/rtctp/internal/config"
	"github.com/zjfcyefeng/rtctp/internal/listener"
)

type TCPServer struct {
	done          chan struct{}
	cfg           config.Config
	taskPool      gxsync.GenericTaskPool
	pkgHandler    getty.ReadWriter
	eventListener getty.EventListener
}

func NewTCPServer(cfg config.Config, log getty.Logger, taskPool gxsync.GenericTaskPool) *TCPServer {
	return &TCPServer{
		done:          make(chan struct{}),
		cfg:           cfg,
		taskPool:      taskPool,
		pkgHandler:    codec.NewJsonReadWriter(),
		eventListener: listener.NewEventListener(cfg, log),
	}
}

func (s *TCPServer) Start() {
	fmt.Println("start tcp server...")
	addr := gxnet.HostAddress(s.cfg.Host, s.cfg.Port)
	serverOpts := []getty.ServerOption{getty.WithLocalAddress(addr)}
	serverOpts = append(serverOpts, getty.WithServerTaskPool(s.taskPool))
	server := getty.NewTCPServer(serverOpts...)
	server.RunEventLoop(s.newSession)

	ticker := time.NewTicker(s.cfg.FailFastTimeout)
	defer ticker.Stop()
	for {
		select {
		case <-ticker.C:
			// TODO
			fmt.Println("ticker time...")
		case <-s.done:
			server.Close()
			return
		}
	}
}

func (s *TCPServer) Stop() {
	fmt.Println("stop tcp server...")
	s.done <- struct{}{}
	close(s.done)
}

func (s *TCPServer) newSession(session getty.Session) error {
	var (
		ok      bool
		tcpConn *net.TCPConn
	)

	if s.cfg.SessionConfig.Compress {
		session.SetCompressType(getty.CompressZip)
	}

	if tcpConn, ok = session.Conn().(*net.TCPConn); !ok {
		panic(fmt.Sprintf("%s, session.conn{%#v} is not tcp connection\n", session.Stat(), session.Conn()))
	}

	tcpConn.SetNoDelay(s.cfg.SessionConfig.NoDelay)
	tcpConn.SetKeepAlive(s.cfg.SessionConfig.KeepAlive)
	if s.cfg.SessionConfig.KeepAlive {
		tcpConn.SetKeepAlivePeriod(s.cfg.SessionConfig.KeepAlivePeriod)
	}
	tcpConn.SetReadBuffer(s.cfg.SessionConfig.ReadBufferBytes)
	tcpConn.SetWriteBuffer(s.cfg.SessionConfig.WriteBufferBytes)

	session.SetName(s.cfg.SessionConfig.Name)
	session.SetMaxMsgLen(s.cfg.MaxBytes)
	session.SetPkgHandler(s.pkgHandler)
	session.SetEventListener(s.eventListener)
	session.SetReadTimeout(s.cfg.SessionConfig.ReadTimeout)
	session.SetWriteTimeout(s.cfg.SessionConfig.WriteTimeout)
	session.SetCronPeriod((int)(s.cfg.SessionTimeout.Nanoseconds() / 1e6))
	session.SetWaitTime(s.cfg.SessionConfig.WaitTimeout)

	return nil
}