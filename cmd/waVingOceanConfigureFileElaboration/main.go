package main

import (
	"bytes"
	"io"
	"os"

	"github.com/golang/protobuf/proto"
	"github.com/ld9999999999/go-interfacetools"
	"github.com/nahanni/go-ucl"
	"github.com/xiaokangwang/waVingOcean/configure"
)

func main() {
	switch os.Args[1] {
	case "T":
		par := ucl.NewParser(os.Stdin)
		out, err := par.Ucl()
		if err != nil {
			panic(err)
		}
		conffile := new(configure.WaVingOceanConfigure)
		err = interfacetools.CopyOut(out, conffile)
		if err != nil {
			panic(err)
		}
		pbout, err := proto.Marshal(conffile)
		if err != nil {
			panic(err)
		}
		io.Copy(os.Stdout, bytes.NewBuffer(pbout))
	case "V":
		var inv2buf bytes.Buffer
		io.Copy(&inv2buf, os.Stdin)
		cfgfd, err := os.Open(os.Args[2])
		if err != nil {
			panic(err)
		}
		var incfbuf bytes.Buffer
		io.Copy(&incfbuf, cfgfd)
		cfgfd.Close()
		conffile := new(configure.WaVingOceanConfigure)
		err = proto.Unmarshal(incfbuf.Bytes(), conffile)
		if err != nil {
			panic(err)
		}
		conffile.V2RayConfigure = inv2buf.Bytes()
		pbout, err := proto.Marshal(conffile)
		if err != nil {
			panic(err)
		}
		io.Copy(os.Stdout, bytes.NewBuffer(pbout))
	}
}
