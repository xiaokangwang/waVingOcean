package wavingocean

import (
	"context"
	"io"

	"github.com/xiaokangwang/waVingOcean/configure"
	"github.com/xiaokangwang/waVingOcean/definition"
	"github.com/yinghuocho/gotun2socks"
)

type LowerUp struct {
	tuns *gotun2socks.Tun2Socks
}

func NewLowerUp(cfg configure.WaVingOceanConfigure, f io.ReadWriteCloser, dialer definition.SurrogateDialer, ctx context.Context) *LowerUp {
	tunc := gotun2socks.New(f, dialer, cfg.DNSServers, cfg.PublicOnly, cfg.EnableDnsCache, ctx)
	return &LowerUp{tuns: tunc}
}

func (l *LowerUp) Up() {
	l.tuns.Run()
}

func (l *LowerUp) Down() {
	l.tuns.Stop()
}
