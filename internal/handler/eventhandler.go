package handler

import (
	"dario.cat/mergo"
	getty "github.com/apache/dubbo-getty"
	"github.com/zjfcyefeng/rtctp/internal/model"
)

const (
	WritePkgTimeout  = 1e9
	StatusSuccess    = 200
	StatusNotFound   = 404
	StatusNotAllowed = 405
	StatusFail       = 500
)

const (
	EventHeartbeat uint16 = iota + 1
)

type EventHandler interface {
	Handle(getty.Session, *model.Request) error
}

type DefaultEventHandler struct {
	log getty.Logger
}

func NewDefaultEventHandler(log getty.Logger) *DefaultEventHandler {
	return &DefaultEventHandler{
		log: log,
	}
}

func (h *DefaultEventHandler) Handle(session getty.Session, req *model.Request) error {
	var resp model.Response
	resp.ID = req.ID
	// TODO business logic
	switch req.Event {
	case EventHeartbeat:
		// TODO heartbeat logic
		heartbeat := model.Heartbeat{}
		mergo.Merge(&heartbeat, req.Body)
		resp.Code = StatusSuccess
		resp.Msg = "pong"
		resp.Data = "pong"
	default:
		resp.Code = StatusFail
		resp.Msg = "unknown event"
		resp.Data = ""
	}
	_, _, err := session.WritePkg(&resp, WritePkgTimeout)
	return err
}
