package wavingocean

import (
	"context"
	"io"
	"net"

	"v2ray.com/core/common/buf"
	v2net "v2ray.com/core/common/net"
)

type V2Dialer struct {
	ser *simpleServer
}

func (vd *V2Dialer) Dial(network, address string, port uint16, ctx context.Context) (net.Conn, error) {
	var dest net.Addr
	switch network {
	case "tcp4":
		dest, _ = net.ResolveTCPAddr(network, address)
	case "udp4":
		dest, _ = net.ResolveUDPAddr(network, address)
	}
	v2dest := v2net.DestinationFromAddr(dest)
	ray, err := vd.ser.disp.Dispatch(ctx, v2dest)
	if err != nil {
		panic(err)
	}
	//Copy data
	conn1, conn2 := net.Pipe()
	go func() {
		buf, _ := ray.InboundOutput().ReadMultiBuffer()
		io.Copy(conn1, &buf)
	}()
	go func() {
		buf, _ := ray.InboundOutput().ReadMultiBuffer()
		io.Copy(&V2WriteAdapter{inter: buf}, conn1)
	}()
	return conn2, nil
}

type V2WriteAdapter struct {
	inter buf.MultiBuffer
}

func (V *V2WriteAdapter) Write(b []byte) (int, error) {
	V.inter.Write(b)
	return len(b), nil
}

func (vd *V2Dialer) NotifyMeltdown(reason error) {}
