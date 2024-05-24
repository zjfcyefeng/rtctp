package main

import (
	"flag"

	getty "github.com/apache/dubbo-getty"
	gxsync "github.com/dubbogo/gost/sync"
	"github.com/zeromicro/go-zero/core/conf"
	"github.com/zeromicro/go-zero/core/service"
	"github.com/zjfcyefeng/rtctp/internal/client"
	"github.com/zjfcyefeng/rtctp/internal/config"
	"github.com/zjfcyefeng/rtctp/internal/server"
	"go.uber.org/zap"
)

var configFile = flag.String("f", "etc/rtctp.yaml", "the config file")

func main() {
	flag.Parse()

	var cfg config.Config
	conf.MustLoad(*configFile, &cfg)

	serviceGroup := service.NewServiceGroup()
	defer serviceGroup.Stop()

	log := getty.GetLogger()
	if cfg.Mode == "pro" || cfg.Mode == "pre" {
		zapLoggerConfig := zap.NewProductionConfig()
		zapLoggerConfig.OutputPaths = []string{"stdout", cfg.LogPath}
		zapLoggerConfig.EncoderConfig = zap.NewProductionEncoderConfig()
		zapLogger, err := zapLoggerConfig.Build(zap.AddCallerSkip(1))
		if err == nil {
			log = zapLogger.Sugar()
			getty.SetLogger(log)
			defer zapLogger.Sync()
		} else {
			log.Error(err)
		}
	}

	taskPool := gxsync.NewTaskPoolSimple(0)
	defer taskPool.Close()

	tcpServer := server.NewTCPServer(cfg, log, taskPool)
	serviceGroup.Add(tcpServer)

	wsServer := server.NewWSServer(cfg, log, taskPool)
	serviceGroup.Add(wsServer)

	tcpClient := client.NewTCPClient(cfg, log)
	serviceGroup.Add(tcpClient)

	wsClient := client.NewWSClient(cfg, log)
	serviceGroup.Add(wsClient)

	serviceGroup.Start()
}
