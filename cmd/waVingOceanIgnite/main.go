package main

import (
	"bytes"
	"io"
	"os"
)
import (
	"github.com/golang/protobuf/proto"
	"github.com/xiaokangwang/waVingOcean"
	"github.com/xiaokangwang/waVingOcean/configure"
)

func main() {
	var inbuf bytes.Buffer
	io.Copy(&inbuf, os.Stdin)
	conffile := new(configure.WaVingOceanConfigure)
	err := proto.Unmarshal(inbuf.Bytes(), conffile)
	if err != nil {
		panic(err)
	}
	if len(os.Args) == 1 {
		if os.Args[0] == "waVingOceanIgnite" {
			wavingocean.Ignite(*conffile)
		} else {
			wavingocean.IgniteNH(*conffile)
		}
	} else {
		wavingocean.Ignite(*conffile)
	}

}
