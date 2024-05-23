package config

import "time"

type Config struct {
	Name            string        `json:",default=server"`
	Host            string        `json:",default=0.0.0.0"`
	Port            int           `json:",default=10086"`
	Path            string        `json:",default=/rtctp"`
	Mode            string        `json:",default=pro,options=dev|test|rt|pre|pro"`
	LogPath         string        `json:",default=rtctp.log"`
	SslEnable       bool          `json:",optional"`
	CertFile        string        `json:",optional"`
	KeyFile         string        `json:",optional"`
	RootCertFile    string        `json:",optional"`
	HeartbeatPeriod time.Duration `json:",default=15s"`
	SessionTimeout  time.Duration `json:",default=60s"`
	FailFastTimeout time.Duration `json:",default=5s"`
	TaskPoolSize    int           `json:",default=0"`
	MaxConns        int           `json:",default=1000"`
	MaxBytes        int           `json:",default=65536"`
	SessionConfig   SessionConfig
}

type SessionConfig struct {
	Name               string        `json:",default=session"`
	NoDelay            bool          `json:",default=true"`
	KeepAlive          bool          `json:",default=true"`
	KeepAlivePeriod    time.Duration `json:",default=180s"`
	ReadBufferBytes    int           `json:",default=65536"`
	WriteBufferBytes   int           `json:",default=65536"`
	WriteQueueCapacity int           `json:",default=1024"`
	ReadTimeout        time.Duration `json:",default=1s"`
	WriteTimeout       time.Duration `json:",default=5s"`
	WaitTimeout        time.Duration `json:",default=8s"`
	Compress           bool          `json:",optional"`
}
