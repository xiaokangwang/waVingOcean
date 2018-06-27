package wavingocean

import (
	"bytes"

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
	nh.InitializeStack(cfg.GetTun().Address, wi, 1500)
}
