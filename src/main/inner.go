package main

import (
	"config"
	"inner"
	"flag"
	"github.com/golang/glog"
	"common"
)

func main() {
	defer glog.Flush()
	flag.Set("logtostderr", "true")
	flag.Parse()


	cf := config.NewConfig(config.ReadArgs(1, "config/inner.properties"))

	outerIp, ok := cf.ReadString("outerIp")
	if !ok {
		glog.Fatal("outerIp is not valid.\n")
	}

	outerPort, ok := cf.ReadInt("outerPort")
	if !ok {
		glog.Fatal("outerPort is not valid.\n")
	}

	service := inner.NewService(common.OuterMsg {
		OuterIp: outerIp,
		OuterPort: outerPort,
	})
	service.Run()
}