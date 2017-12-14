package definition

import (
	"context"
	"net"
)

/*
SurrogateDialer represent a Dial interface
for connections that needs to be transmited
with an alternative tunnel
*/
type SurrogateDialer interface {
	Dial(network, address string, port uint16, ctx context.Context) (net.Conn, error)
	NotifyMeltdown(reason error)
}
