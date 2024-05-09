package bytex

import (
	"io"
	gxbytes "github.com/dubbogo/gost/bytes"
)

type ByteBuffer struct {
	buf *gxbytes.Buffer
}

func NewByteBuffer(bytes []byte) *ByteBuffer {
	return &ByteBuffer{
		buf: gxbytes.NewBuffer(bytes),
	}
}

func (b *ByteBuffer) Bytes() []byte {
	return b.buf.Bytes()
}

func (b *ByteBuffer) Read(p []byte) (n int, err error) {
	return b.buf.Read(p)
}

func (b *ByteBuffer) ReadByte() (byte, error) {
	data := make([]byte, 1)
	n, err := b.Read(data)
	if err != nil {
		return 0, err
	}
	if n < 1 {
		return 0, io.ErrShortBuffer
	}
	return data[0], nil
}

func (b *ByteBuffer) Write(p []byte) (n int, err error) {
	return b.buf.Write(p)
}