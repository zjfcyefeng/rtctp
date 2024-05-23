package collection

import "strings"

const (
	KvSplit   = "="
	PairSplit = "&"
)

var (
	kvSplitBytes   = []byte(KvSplit)
	pairSplitBytes = []byte(PairSplit)
)

func EncodeMap(dataMap map[string]string) []byte {
	if dataMap == nil {
		return nil
	}

	bytes := make([]byte, 0)
	if len(dataMap) == 0 {
		return bytes
	}

	for k, v := range dataMap {
		bytes = append(bytes, []byte(k)...)
		bytes = append(bytes, kvSplitBytes...)
		bytes = append(bytes, []byte(v)...)
		bytes = append(bytes, pairSplitBytes...)
	}

	return bytes[:len(bytes)-1]
}

func DecodeMap(data []byte) map[string]string {
	if data == nil {
		return nil
	}

	ctxMap := make(map[string]string, 0)

	dataStr := string(data)
	if dataStr == "" {
		return ctxMap
	}

	kvPairs := strings.Split(dataStr, PairSplit)
	if len(kvPairs) == 0 {
		return ctxMap
	}

	for _, kvPair := range kvPairs {
		if kvPair == "" {
			continue
		}

		kvs := strings.Split(kvPair, KvSplit)
		if len(kvs) != 2 {
			continue
		}

		ctxMap[kvs[0]] = kvs[1]
	}

	return ctxMap
}