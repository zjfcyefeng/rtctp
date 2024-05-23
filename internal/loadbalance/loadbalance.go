package loadbalance

import (
	"strings"
	"sync"

	getty "github.com/apache/dubbo-getty"
)

func Select(sessions *sync.Map, xid string) getty.Session {
	var session getty.Session

	// ip:port:transactionId
	splits := strings.Split(xid, ":")
	if len(splits) == 3 {
		ip := splits[0]
		port := splits[1]
		host := ip+":"+port
		sessions.Range(func(key, value interface{}) bool {
			tmpSession := key.(getty.Session)
			if tmpSession.IsClosed() {
				sessions.Delete(tmpSession)
				return true
			}
			remoteAddr := tmpSession.RemoteAddr()
			if host == remoteAddr {
				session = tmpSession
				return false
			}
			return true
		})
	}
	return session
}