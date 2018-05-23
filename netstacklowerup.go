package wavingocean

import (
	"bytes"

	"github.com/xiaokangwang/waVingOcean/configure"
	"v2ray.com/core"
	//"github.com/xiaokangwang/waVingOcean/definition"
	"github.com/xiaokangwang/waVingOcean/netstackadoptor"
)

func IgniteNH(cfg configure.WaVingOceanConfigure, nh *netstackadoptor.NetstackHolder) {
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
	wi, err := netstackadoptor.OpenDefaultWaterInterface()
	if err != nil {
		panic(err)
	}
	nh.SetDialer(&V2Dialer{ser: ns})
	nh.InitializeStack(cfg.GetTun().Address, wi, 1500)
}
