package listener

import (
	"errors"
	"slices"
	"sync"
	"time"

	getty "github.com/apache/dubbo-getty"
	"github.com/zjfcyefeng/rtctp/internal/config"
	"github.com/zjfcyefeng/rtctp/internal/handler"
	"github.com/zjfcyefeng/rtctp/internal/model"
)

var (
	ErrTooManySessions = errors.New("too many sessions")
)

type DefaultEventListener struct {
	locker   sync.RWMutex
	cfg      config.Config
	log      getty.Logger
	handler  handler.EventHandler
	sessions []getty.Session
}

func NewEventListener(cfg config.Config, log getty.Logger) *DefaultEventListener {

	return &DefaultEventListener{
		cfg:      cfg,
		log:      log,
		handler:  handler.NewDefaultEventHandler(log),
		sessions: make([]getty.Session, 0, cfg.MaxConns),
	}
}

func (e *DefaultEventListener) OnOpen(session getty.Session) error {
	var err error
	e.locker.RLock()
	if len(e.sessions) >= e.cfg.MaxConns {
		err = ErrTooManySessions
	}
	e.locker.RUnlock()
	if err != nil {
		return err
	}

	e.log.Debug("got new session{%s}", session.Stat())
	e.locker.Lock()
	e.sessions = append(e.sessions, session)
	e.locker.Unlock()
	return nil
}

func (e *DefaultEventListener) OnCron(session getty.Session) {
	var (
		flag   bool
		active time.Time
	)
	e.locker.RLock()
	if slices.Contains(e.sessions, session) {
		active = session.GetActive()
		if time.Since(active).Nanoseconds() > e.cfg.SessionTimeout.Nanoseconds() {
			flag = true
			e.log.Warnf("session{%s} timeout{%s}", session.Stat(), time.Since(active).String())
		}
	}
	e.locker.RUnlock()

	if flag {
		e.locker.Lock()
		e.sessions = slices.DeleteFunc(e.sessions, func(s getty.Session) bool {
			return slices.Contains(e.sessions, session)
		})
		e.locker.Unlock()
		session.Close()
	}
}

func (e *DefaultEventListener) OnMessage(session getty.Session, pkg interface{}) {
	req, ok := pkg.(*model.Request)
	if !ok {
		e.log.Warnf("illegal request{%#v}]", pkg)
		return
	}

	err := e.handler.Handle(session, req)
	if err != nil {
		e.log.Errorf("handle request{%#v} error{%v}", pkg, err)
	}
}

func (e *DefaultEventListener) OnClose(session getty.Session) {
	e.log.Debug("session{%s} is closing......", session.Stat())
	e.locker.Lock()
	defer e.locker.Unlock()
	e.sessions = slices.DeleteFunc(e.sessions, func(s getty.Session) bool {
		return slices.Contains(e.sessions, session)
	})
}

func (e *DefaultEventListener) OnError(session getty.Session, err error) {
	e.log.Errorf("session{%s} got error{%v}, will be closed", session.Stat(), err)
	e.locker.Lock()
	defer e.locker.Unlock()
	e.sessions = slices.DeleteFunc(e.sessions, func(s getty.Session) bool {
		return slices.Contains(e.sessions, session)
	})
}
