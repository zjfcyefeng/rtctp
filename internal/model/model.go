package model

type Header struct {
	Version string `json:"version"`
	Xid     string `json:"xid"`
}

type Request struct {
	Header
	ID    uint64                 `json:"id"`
	Event uint16                 `json:"event"`
	Body  map[string]interface{} `json:"body"`
}

type Response struct {
	ID   uint64      `json:"id"`
	Code int         `json:"code"`
	Msg  string      `json:"msg"`
	Data interface{} `json:"data"`
}

type Heartbeat struct {
	Timestamp int64 `json:"timestamp"`
}
