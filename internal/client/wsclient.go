package client

import (
	"crypto/tls"
	"fmt"
	"net"

	getty "github.com/apache/dubbo-getty"
	gxnet "github.com/dubbogo/gost/net"
	"github.com/zjfcyefeng/rtctp/internal/codec"
	"github.com/zjfcyefeng/rtctp/internal/config"
	"github.com/zjfcyefeng/rtctp/internal/listener"
)

type WSClient struct {
	done          chan struct{}
	cfg           config.Config
	log           getty.Logger
	pkgHandler    getty.ReadWriter
	eventListener getty.EventListener
}

func NewWSClient(cfg config.Config, log getty.Logger) *WSClient {
	return &WSClient{
		done:          make(chan struct{}),
		cfg:           cfg,
		log:           log,
		pkgHandler:    codec.NewJsonResponseReadWriter(),
		eventListener: listener.NewClientEventListener(cfg, log),
	}
}

func (c *WSClient) Start() {
	fmt.Println("start ws client...")
	clientOpts := []getty.ClientOption{getty.WithServerAddress(gxnet.WSHostAddress(c.cfg.Host, c.cfg.Port+2, c.cfg.Path))}
	clientOpts = append(clientOpts, getty.WithConnectionNumber(1))
	client := getty.NewWSClient(clientOpts...)
	defer client.Close()
	client.RunEventLoop(c.newSession)

	for range c.done {
		close(c.done)
		return
	}
}

func (c *WSClient) Stop() {
	fmt.Println("stop ws client...")
	c.done <- struct{}{}
}

func (c *WSClient) newSession(session getty.Session) error {
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
		tcpConn.SetNoDelay(c.cfg.SessionConfig.NoDelay)
		tcpConn.SetKeepAlive(c.cfg.SessionConfig.KeepAlive)
		if c.cfg.SessionConfig.KeepAlive {
			tcpConn.SetKeepAlivePeriod(c.cfg.SessionConfig.KeepAlivePeriod)
		}
		tcpConn.SetReadBuffer(c.cfg.SessionConfig.ReadBufferBytes)
		tcpConn.SetWriteBuffer(c.cfg.SessionConfig.WriteBufferBytes)
	}

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
