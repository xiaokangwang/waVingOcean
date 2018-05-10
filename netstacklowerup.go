package wavingocean

import (
	"bytes"
	"errors"
	"net"
	"strings"

	"github.com/xiaokangwang/waVingOcean/configure"
	"v2ray.com/core"

	netstack "github.com/FlowerWrong/netstack/tcpip/stack"

	"github.com/FlowerWrong/netstack/tcpip"
	"github.com/FlowerWrong/netstack/tcpip/link/fdbased"
	"github.com/FlowerWrong/netstack/tcpip/network/ipv4"
	"github.com/FlowerWrong/netstack/tcpip/network/ipv6"
	"github.com/FlowerWrong/netstack/tcpip/transport/tcp"
	"github.com/FlowerWrong/netstack/tcpip/transport/udp"
	"github.com/FlowerWrong/water"

	"github.com/FlowerWrong/netstack/waiter"
	"github.com/xiaokangwang/waVingOcean/definition"
)

type NetstackHolder struct {
	nstack *netstack.Stack
	dialer definition.SurrogateDialer
}

const (
	// NICId is global nicid for stack
	NICId = 1
	// Backlog is tcp listen backlog
	Backlog = 1024
)
const netstackHookport = 45001

func (nh *NetstackHolder) Ignite(cfg configure.WaVingOceanConfigure) {
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
}

func (nh *NetstackHolder) setupTCPHandler() error {
	var wq waiter.Queue
	ep, err := nh.nstack.NewEndpoint(tcp.ProtocolNumber, 4, &wq)
	if err != nil {
		return errors.New(err.String())
	}

	defer ep.Close()
	if err := ep.Bind(tcpip.FullAddress{NICId, "", netstackHookport}, nil); err != nil {
		return errors.New(err.String())
	}
	if err := ep.Listen(Backlog); err != nil {
		return errors.New(err.String())
	}

	// Wait for connections to appear.
	waitEntry, notifyCh := waiter.NewChannelEntry(nil)
	wq.EventRegister(&waitEntry, waiter.EventIn)
	defer wq.EventUnregister(&waitEntry)

}

func (nh *NetstackHolder) initializeStack(tunip string, ifce *water.Interface, mtu uint32) {
	tunIP, _, _ := net.ParseCIDR(tunip)

	var addr tcpip.Address
	var proto tcpip.NetworkProtocolNumber
	if tunIP.To4() != nil {
		addr = tcpip.Address(tunIP.To4())
		proto = ipv4.ProtocolNumber
	} else if tunIP.To16() != nil {
		addr = tcpip.Address(tunIP.To16())
		proto = ipv6.ProtocolNumber
	} else {
		//log.Fatalf("Unknown IP type: %v", app.Cfg.General.Network)
	}
	nh.nstack = netstack.New([]string{ipv4.ProtocolName, ipv6.ProtocolName}, []string{tcp.ProtocolName, udp.ProtocolName})

	// Parse the mac address.
	maddr, err := net.ParseMAC("aa:00:01:01:01:01")
	if err != nil {
		//log.Fatalf("Bad MAC address: aa:00:01:01:01:01")
	}

	linkID := fdbased.New(ifce, &fdbased.Options{
		FD:             ifce.Fd(),
		MTU:            mtu,
		EthernetHeader: false,
		Address:        tcpip.LinkAddress(maddr),
	})

	if err := nh.nstack.CreateNIC(NICId, linkID, true, addr, netstackHookport); err != nil {
		//log.Fatal("Create NIC failed", err)
	}
	if err := nh.nstack.AddAddress(NICId, proto, addr); err != nil {
		//log.Fatal("Add address to stack failed", err)
	}

	// Add default route.
	nh.nstack.SetRouteTable([]tcpip.Route{
		{
			Destination: tcpip.Address(strings.Repeat("\x00", len(addr))),
			Mask:        tcpip.Address(strings.Repeat("\x00", len(addr))),
			Gateway:     "",
			NIC:         NICId,
		},
	})
}
