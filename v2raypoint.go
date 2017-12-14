package wavingocean

import (
	"context"

	"v2ray.com/core"
	"v2ray.com/core/app"
	"v2ray.com/core/app/dispatcher"
	"v2ray.com/core/app/dns"
	"v2ray.com/core/app/log"
	"v2ray.com/core/app/policy"
	"v2ray.com/core/app/proxyman"
	"v2ray.com/core/common"
	"v2ray.com/core/common/net"
)

// simpleServer shell of V2Ray.
type simpleServer struct {
	space app.Space
	disp  dispatcher.Interface
}

// newSimpleServer returns a new Point server based on given configuration.
// The server is not started at this point.
func newSimpleServer(config *core.Config) (*simpleServer, error) {
	var server = new(simpleServer)

	if err := config.Transport.Apply(); err != nil {
		return nil, err
	}

	space := app.NewSpace()
	ctx := app.ContextWithSpace(context.Background(), space)

	server.space = space

	for _, appSettings := range config.App {
		settings, err := appSettings.GetInstance()
		if err != nil {
			return nil, err
		}
		application, err := app.CreateAppFromConfig(ctx, settings)
		if err != nil {
			return nil, err
		}
		if err := space.AddApplication(application); err != nil {
			return nil, err
		}
	}

	if log.FromSpace(space) == nil {
		l, err := app.CreateAppFromConfig(ctx, &log.Config{
			ErrorLogType:  log.LogType_Console,
			ErrorLogLevel: log.LogLevel_Warning,
			AccessLogType: log.LogType_None,
		})
		if err != nil {
			return nil, err //newError("failed apply default log settings").Base(err)
		}
		common.Must(space.AddApplication(l))
	}

	outboundHandlerManager := proxyman.OutboundHandlerManagerFromSpace(space)
	if outboundHandlerManager == nil {
		o, err := app.CreateAppFromConfig(ctx, new(proxyman.OutboundConfig))
		if err != nil {
			return nil, err
		}
		if err := space.AddApplication(o); err != nil {
			return nil, err //newError("failed to add default outbound handler manager").Base(err)
		}
		outboundHandlerManager = o.(proxyman.OutboundHandlerManager)
	}

	if dns.FromSpace(space) == nil {
		dnsConfig := &dns.Config{
			NameServers: []*net.Endpoint{{
				Address: net.NewIPOrDomain(net.LocalHostDomain),
			}},
		}
		d, err := app.CreateAppFromConfig(ctx, dnsConfig)
		if err != nil {
			return nil, err
		}
		common.Must(space.AddApplication(d))
	}

	if disp := dispatcher.FromSpace(space); disp == nil {
		d, err := app.CreateAppFromConfig(ctx, new(dispatcher.Config))
		if err != nil {
			return nil, err
		}
		common.Must(space.AddApplication(d))
		server.disp = disp
	}

	if p := policy.FromSpace(space); p == nil {
		p, err := app.CreateAppFromConfig(ctx, &policy.Config{
			Level: map[uint32]*policy.Policy{
				1: {
					Timeout: &policy.Policy_Timeout{
						ConnectionIdle: &policy.Second{
							Value: 600,
						},
					},
				},
			},
		})
		if err != nil {
			return nil, err
		}
		common.Must(space.AddApplication(p))
	}

	for _, outbound := range config.Outbound {
		if err := outboundHandlerManager.AddHandler(ctx, outbound); err != nil {
			return nil, err
		}
	}

	if err := server.space.Initialize(); err != nil {
		return nil, err
	}

	return server, nil
}

func (s *simpleServer) Close() {
	s.space.Close()
}

func (s *simpleServer) Start() error {
	if err := s.space.Start(); err != nil {
		return err
	}
	//log.Trace(newError("V2Ray started").AtWarning())

	return nil
}
