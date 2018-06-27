package wavingocean

import (
	"bytes"
	"log"
	"runtime"
	"time"

	"github.com/xiaokangwang/waVingOcean/configure"
	"v2ray.com/core"
	//"github.com/xiaokangwang/waVingOcean/definition"
	"github.com/xiaokangwang/waVingOcean/netstackadoptor"
	_ "v2ray.com/core/main/distro/all"
)

func IgniteNH(cfg configure.WaVingOceanConfigure) {
	configure, err := core.LoadConfig("protobuf", "", bytes.NewBuffer(cfg.V2RayConfigure))
	if err != nil {
		panic(err)
	}
	ns, err := core.New(configure)
	if err != nil {
		panic(err)
	}
	err = ns.Start()
	if err != nil {
		panic(err)
	}
	var nh *netstackadoptor.NetstackHolder
	nh = &netstackadoptor.NetstackHolder{}
	wi, err := netstackadoptor.OpenDefaultWaterInterface()
	if err != nil {
		panic(err)
	}

	nh.SetDialer(&V2Dialer{ser: ns})
	log.Print("Dialer Set")
	nh.InitializeStack(cfg.GetTun().GetGateway(), wi, 1500)
	log.Print("Stack init")
	for {
		runtime.Gosched()
		time.Sleep(time.Minute)
	}
}
