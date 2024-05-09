package bytex

import "encoding/binary"

func Byte2Int64(b []byte) int64 {
    return int64((int64(b[0])&0xff)<<56 | 
		(int64(b[1])&0xff)<<48 | 
		(int64(b[2])&0xff)<<40 | 
		(int64(b[3])&0xff)<<32 |
        (int64(b[4])&0xff)<<24 | 
		(int64(b[5])&0xff)<<16 | 
		(int64(b[6])&0xff)<<8 | 
		(int64(b[7])&0xff))
}

func Byte2UInt64(b []byte) uint64 {
	return binary.BigEndian.Uint64(b)
}