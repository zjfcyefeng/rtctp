package handler

import (
	"time"

	"dario.cat/mergo"
	getty "github.com/apache/dubbo-getty"
	"github.com/zjfcyefeng/rtctp/internal/model"
)

const (
	WritePkgTimeout          = 1e9
	StatusContinue           = 100
	StatusSuccess            = 200
	StatusFound              = 302
	StatusBadRequest         = 400
	StatusUnauthorized       = 401
	StatusNotFound           = 404
	StatusNotAllowed         = 405
	StatusFail               = 500
	StatusServiceUnavailable = 503
	StatusGatewayTimeout     = 504
)

const (
	EventHeartbeat uint16 = iota + 10000
	EventEcho
)

type EventHandler interface {
	HandleRequest(getty.Session, *model.Request) error
	HandleResponse(getty.Session, *model.Response) error
}

type DefaultEventHandler struct {
	log getty.Logger
}

func NewDefaultEventHandler(log getty.Logger) *DefaultEventHandler {
	return &DefaultEventHandler{
		log: log,
	}
}

func (h *DefaultEventHandler) HandleRequest(session getty.Session, req *model.Request) error {
	var resp model.Response
	resp.Xid = req.Xid
	resp.ID = req.ID
	// TODO business logic
	switch req.Event {
	case EventHeartbeat:
		// TODO heartbeat logic
		heartbeat := model.Heartbeat{
			Timestamp: time.Now().UnixMilli(),
		}
		resp.Code = StatusSuccess
		resp.Msg = "pong"
		resp.Data = heartbeat
	case EventEcho:
		echo := model.Echo{}
		mergo.Map(&echo, req.Body)
		h.log.Debugf("receive echo request: %#v", echo)
	default:
		resp.Code = StatusFail
		resp.Msg = "unknown event"
		resp.Data = ""
	}
	h.log.Debugf("send response: %#v", resp)
	_, _, err := session.WritePkg(&resp, WritePkgTimeout)
	return err
}

func (h *DefaultEventHandler) HandleResponse(session getty.Session, resp *model.Response) error {
	h.log.Debugf("handle response: %#v", resp)
	return nil
}
