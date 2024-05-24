package listener

import (
	"time"

	getty "github.com/apache/dubbo-getty"
	"github.com/zjfcyefeng/rtctp/internal/config"
	"github.com/zjfcyefeng/rtctp/internal/handler"
	"github.com/zjfcyefeng/rtctp/internal/model"
)

type ClientEventListener struct {
	cfg     config.Config
	log     getty.Logger
	handler handler.EventHandler
}

func NewClientEventListener(cfg config.Config, log getty.Logger) *ClientEventListener {

	return &ClientEventListener{
		cfg:     cfg,
		log:     log,
		handler: handler.NewDefaultEventHandler(log),
	}
}

func (e *ClientEventListener) OnOpen(session getty.Session) error {
	e.log.Debug("got new session{%s}", session.Stat())
	return nil
}

func (e *ClientEventListener) OnCron(session getty.Session) {
	req := model.Request{
		Version: "1.0.0",
		Xid:     "rtctp",
		ID:      uint64(time.Now().UnixMilli()),
		Event:   handler.EventHeartbeat,
		Body:    map[string]interface{}{},
	}
	_, _, _ = session.WritePkg(req, handler.WritePkgTimeout)
}

func (e *ClientEventListener) OnMessage(session getty.Session, pkg interface{}) {
	resp, ok := pkg.(*model.Response)
	if !ok {
		e.log.Warnf("illegal response{%#v}]", pkg)
		return
	}

	err := e.handler.HandleResponse(session, resp)
	if err != nil {
		e.log.Errorf("handle response{%#v} error{%v}", pkg, err)
	}
}

func (e *ClientEventListener) OnClose(session getty.Session) {
	e.log.Debug("session{%s} is closing......", session.Stat())
}

func (e *ClientEventListener) OnError(session getty.Session, err error) {
	e.log.Errorf("session{%s} got error{%v}, will be closed", session.Stat(), err)
}
