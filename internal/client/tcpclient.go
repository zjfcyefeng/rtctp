package client

import (
	"fmt"
	"net"

	getty "github.com/apache/dubbo-getty"
	gxnet "github.com/dubbogo/gost/net"
	"github.com/zjfcyefeng/rtctp/internal/codec"
	"github.com/zjfcyefeng/rtctp/internal/config"
	"github.com/zjfcyefeng/rtctp/internal/listener"
)

type TCPClient struct {
	done          chan struct{}
	cfg           config.Config
	log           getty.Logger
	pkgHandler    getty.ReadWriter
	eventListener getty.EventListener
}

func NewTCPClient(cfg config.Config, log getty.Logger) *TCPClient {
	return &TCPClient{
		done:          make(chan struct{}),
		cfg:           cfg,
		log:           log,
		pkgHandler:    codec.NewJsonResponseReadWriter(),
		eventListener: listener.NewClientEventListener(cfg, log),
	}
}

func (c *TCPClient) Start() {
	fmt.Println("start tcp client...")
	clientOpts := []getty.ClientOption{getty.WithServerAddress(gxnet.HostAddress(c.cfg.Host, c.cfg.Port))}
	clientOpts = append(clientOpts, getty.WithConnectionNumber(1))
	client := getty.NewTCPClient(clientOpts...)
	defer client.Close()
	client.RunEventLoop(c.newSession)

	for range c.done {
		close(c.done)
		return
	}
}

func (c *TCPClient) Stop() {
	fmt.Println("stop tcp client...")
	c.done <- struct{}{}
}

func (c *TCPClient) newSession(session getty.Session) error {
	var (
		ok      bool
		tcpConn *net.TCPConn
	)

	if tcpConn, ok = session.Conn().(*net.TCPConn); !ok {
		panic(fmt.Sprintf("%s, session.conn{%#v} is not tcp connection\n", session.Stat(), session.Conn()))
	}

	tcpConn.SetNoDelay(c.cfg.SessionConfig.NoDelay)
	tcpConn.SetKeepAlive(c.cfg.SessionConfig.KeepAlive)
	if c.cfg.SessionConfig.KeepAlive {
		tcpConn.SetKeepAlivePeriod(c.cfg.SessionConfig.KeepAlivePeriod)
	}
	tcpConn.SetReadBuffer(c.cfg.SessionConfig.ReadBufferBytes)
	tcpConn.SetWriteBuffer(c.cfg.SessionConfig.WriteBufferBytes)

	if c.cfg.SessionConfig.Compress {
		session.SetCompressType(getty.CompressZip)
	}
	session.SetName(c.cfg.SessionConfig.Name)
	session.SetMaxMsgLen(c.cfg.MaxBytes)
	session.SetPkgHandler(c.pkgHandler)
	session.SetEventListener(c.eventListener)
	session.SetReadTimeout(c.cfg.SessionConfig.ReadTimeout)
	session.SetWriteTimeout(c.cfg.SessionConfig.WriteTimeout)
	session.SetCronPeriod((int)(c.cfg.HeartbeatPeriod.Nanoseconds()/1e6) - 100)
	session.SetWaitTime(c.cfg.SessionConfig.WaitTimeout)

	return nil
}
