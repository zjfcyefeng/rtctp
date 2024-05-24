package codec

import (
	"encoding/json"

	getty "github.com/apache/dubbo-getty"
	gxbytes "github.com/dubbogo/gost/bytes"
	"github.com/zjfcyefeng/rtctp/internal/model"
)

type JsonResponseReadWriter struct {
}

func NewJsonResponseReadWriter() *JsonResponseReadWriter {
	return &JsonResponseReadWriter{}
}

func (c *JsonResponseReadWriter) Read(session getty.Session, data []byte) (interface{}, int, error) {
	// Read Parse tcp/udp/websocket pkg from buffer and if possible return a complete pkg.
	// When receiving a tcp network streaming segment, there are 4 cases as following:
	// case 1: a error found in the streaming segment;
	// case 2: can not unmarshal a pkg header from the streaming segment;
	// case 3: unmarshal a pkg header but can not unmarshal a pkg from the streaming segment;
	// case 4: just unmarshal a pkg from the streaming segment;
	// case 5: unmarshal more than one pkg from the streaming segment;
	//
	// The return value is (nil, 0, error) as case 1.
	// The return value is (nil, 0, nil) as case 2.
	// The return value is (nil, pkgLen, nil) as case 3.
	// The return value is (pkg, pkgLen, nil) as case 4.
	// The handleTcpPackage may invoke func Read many times as case 5.
	var req model.Response
	reqBuf := gxbytes.NewBuffer(data)
	reqLen := reqBuf.Len()
	err := json.Unmarshal(reqBuf.Bytes(), &req)
	if err != nil {
		return nil, 0, err
	}

	return &req, reqLen, nil
}

func (c *JsonResponseReadWriter) Write(session getty.Session, pkg interface{}) ([]byte, error) {
	// Write if @Session is udpGettySession, the second parameter is UDPContext.
	return json.Marshal(&pkg)
}
