package server

import (
	"crypto/tls"
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

type WSServer struct {
	done          chan struct{}
	cfg           config.Config
	taskPool      gxsync.GenericTaskPool
	pkgHandler    getty.ReadWriter
	eventListener getty.EventListener
}

func NewWSServer(cfg config.Config, log getty.Logger, taskPool gxsync.GenericTaskPool) *WSServer {
	return &WSServer{
		done:          make(chan struct{}),
		cfg:           cfg,
		taskPool:      taskPool,
		pkgHandler:    codec.NewJsonRequestReadWriter(),
		eventListener: listener.NewServerEventListener(cfg, log),
	}
}

func (s *WSServer) Start() {
	addr := gxnet.HostAddress(s.cfg.Host, s.cfg.Port+2)
	fmt.Printf("start ws server at ws://%s%s...\n", addr, s.cfg.Path)
	serverOpts := []getty.ServerOption{getty.WithLocalAddress(addr)}
	serverOpts = append(serverOpts, getty.WithServerTaskPool(s.taskPool))
	serverOpts = append(serverOpts, getty.WithWebsocketServerPath(s.cfg.Path))
	server := getty.NewWSServer(serverOpts...)
	defer server.Close()
	server.RunEventLoop(s.newSession)

	ticker := time.NewTicker(s.cfg.HeartbeatPeriod)
	defer ticker.Stop()
	for {
		select {
		case <-ticker.C:
			// TODO
			fmt.Println("ws server ticker time...")
		case <-s.done:
			close(s.done)
			return
		}
	}
}

func (s *WSServer) Stop() {
	fmt.Println("stop ws server...")
	s.done <- struct{}{}
}

func (s *WSServer) newSession(session getty.Session) error {
	var (
		flagTLS, flagTCP bool
		tcpConn          *net.TCPConn
	)

	_, flagTLS = session.Conn().(*tls.Conn)
	tcpConn, flagTCP = session.Conn().(*net.TCPConn)
	if !flagTLS && !flagTCP {
		panic(fmt.Sprintf("%s, session.conn{%#v} is not tcp/tls connection\n", session.Stat(), session.Conn()))
	}

	if flagTCP {
		tcpConn.SetNoDelay(s.cfg.SessionConfig.NoDelay)
		tcpConn.SetKeepAlive(s.cfg.SessionConfig.KeepAlive)
		if s.cfg.SessionConfig.KeepAlive {
			tcpConn.SetKeepAlivePeriod(s.cfg.SessionConfig.KeepAlivePeriod)
		}
		tcpConn.SetReadBuffer(s.cfg.SessionConfig.ReadBufferBytes)
		tcpConn.SetWriteBuffer(s.cfg.SessionConfig.WriteBufferBytes)
	}

	if s.cfg.SessionConfig.Compress {
		session.SetCompressType(getty.CompressZip)
	}
	session.SetName(s.cfg.SessionConfig.Name)
	session.SetMaxMsgLen(s.cfg.MaxBytes)
	session.SetPkgHandler(s.pkgHandler)
	session.SetEventListener(s.eventListener)
	session.SetReadTimeout(s.cfg.SessionConfig.ReadTimeout)
	session.SetWriteTimeout(s.cfg.SessionConfig.WriteTimeout)
	session.SetCronPeriod((int)(s.cfg.HeartbeatPeriod.Nanoseconds() / 1e6))
	session.SetWaitTime(s.cfg.SessionConfig.WaitTimeout)

	return nil
}
