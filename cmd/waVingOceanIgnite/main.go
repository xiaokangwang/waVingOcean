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
	wavingocean.Ignite(*conffile)
}
