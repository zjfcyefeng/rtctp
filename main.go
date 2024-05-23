package main

import (
	"flag"

	getty "github.com/apache/dubbo-getty"
	gxsync "github.com/dubbogo/gost/sync"
	"github.com/zeromicro/go-zero/core/conf"
	"github.com/zeromicro/go-zero/core/service"
	"github.com/zjfcyefeng/rtctp/internal/config"
	"github.com/zjfcyefeng/rtctp/internal/server"
)

var configFile = flag.String("f", "etc/rtctp.yaml", "the config file")

func main() {
	flag.Parse()

	var cfg config.Config
	conf.MustLoad(*configFile, &cfg)

	serviceGroup := service.NewServiceGroup()
	defer serviceGroup.Stop()

	log := getty.GetLogger()
	taskPool := gxsync.NewTaskPoolSimple(0)
	defer taskPool.Close()

	tcpServer := server.NewTCPServer(cfg, log, taskPool)
	serviceGroup.Add(tcpServer)

	wsServer := server.NewWSServer(cfg, log, taskPool)
	serviceGroup.Add(wsServer)

	serviceGroup.Start()
}
