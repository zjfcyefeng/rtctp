package bytex
func ReadBytes(n int, buf *ByteBuffer) []byte {
	bytes := make([]byte, n)
	buf.Read(bytes)
	return bytes
}

func ReadByte(buf *ByteBuffer) byte {
	value, _ := buf.ReadByte()
	return value
}

func ReadUint8(buf *ByteBuffer) uint8 {
	value, _ := buf.ReadByte()
	return value
}

func ReadUInt16(buf *ByteBuffer) uint16 {
	value, _ := buf.ReadUint16()
	return value
}

func ReadUInt32(buf *ByteBuffer) uint32 {
	value, _ := buf.ReadUint32()
	return value
}

func ReadUInt64(buf *ByteBuffer) uint64 {
	value, _ := buf.ReadUint64()
	return value
}

func ReadString8(buf *ByteBuffer) string {
	bytes := make([]byte, 1)
	buf.Read(bytes)
	return string(bytes)
}

func Read1String16(buf *ByteBuffer) string {
	bytes := make([]byte, 2)
	buf.Read(bytes)
	return string(bytes)
}

func ReadString32(buf *ByteBuffer) string {
	bytes := make([]byte, 4)
	buf.Read(bytes)
	return string(bytes)
}

func ReadString64(buf *ByteBuffer) string {
	bytes := make([]byte, 8)
	buf.Read(bytes)
	return string(bytes)
}

func ReadString8Length(buf *ByteBuffer) string {
	length, _ := buf.ReadByte()
	if length > 0 {
		p := make([]byte, length)
		buf.Read(p)
		return string(p)
	}
	return ""
}

func ReadString16Length(buf *ByteBuffer) string {
	length, _ := buf.ReadUint16()
	if length > 0 {
		p := make([]byte, length)
		buf.Read(p)
		return string(p)
	}
	return ""
}

func ReadString32Length(buf *ByteBuffer) string {
	length, _ := buf.ReadUint32()
	if length > 0 {
		p := make([]byte, length)
		buf.Read(p)
		return string(p)
	}
	return ""
}

func ReadString64Length(buf *ByteBuffer) string {
	length, _ := buf.ReadUint64()
	if length > 0 {
		p := make([]byte, length)
		buf.Read(p)
		return string(p)
	}
	return ""
}

func WriteString8Length(value string, buf *ByteBuffer) {
	if value != "" {
		buf.WriteByte(byte(len(value)))
		buf.WriteString(value)
	} else {
		buf.WriteByte(byte(0))
	}
}

func WriteString16Length(value string, buf *ByteBuffer) {
	if value != "" {
		buf.WriteUint16(uint16(len(value)))
		buf.WriteString(value)
	} else {
		buf.WriteUint16(uint16(0))
	}
}

func WriteString32Length(value string, buf *ByteBuffer) {
	if value != "" {
		buf.WriteUint32(uint32(len(value)))
		buf.WriteString(value)
	} else {
		buf.WriteUint32(uint32(0))
	}
}

func WriteString64Length(value string, buf *ByteBuffer) {
	if value != "" {
		buf.WriteUint64(uint64(len(value)))
		buf.WriteString(value)
	} else {
		buf.WriteUint64(uint64(0))
	}
}