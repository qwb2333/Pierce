package main

import (
	"github.com/qwb2333/Pierce/config"
	"flag"
	"github.com/golang/glog"
	"github.com/qwb2333/Pierce/outer"
	"github.com/qwb2333/Pierce/common"
)

func main() {
	defer glog.Flush()
	flag.Set("logtostderr", "true")
	flag.Parse()

	cf := config.NewConfig(config.ReadArgs(1, "config/outer.properties"))
	msgLists := config.ReadDotConf(config.ReadArgs(2, "config/outer_dot.conf"))

	outerIp, ok := cf.ReadString("outerIp")
	if !ok {
		glog.Fatal("outerIp is not valid.\n")
	}

	outerPort, ok := cf.ReadInt("outerPort")
	if !ok {
		glog.Fatal("outerPort is not valid.\n")
	}

	service := outer.NewService(common.OuterMsg {
		OuterIp: outerIp,
		OuterPort: outerPort,
	})

	service.Run(msgLists)
}