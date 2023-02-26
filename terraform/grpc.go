package terraform

import (
	"context"

	"github.com/hashicorp/go-plugin"
	"github.com/rvolykh/intellij-hcl-schema/proto"
	"google.golang.org/grpc"
)

const (
	// Maximum response size for schema request, 256MB
	maxRecvSize = 256 << 20
)

// GRPCPlugin is a wrapper around plugin.Plugin to be able to configure grpc client
type GRPCPlugin struct {
	plugin.Plugin
}

func (p *GRPCPlugin) GRPCClient(ctx context.Context, broker *plugin.GRPCBroker, c *grpc.ClientConn) (interface{}, error) {
	return proto.NewProviderClient(c), nil
}

func (p *GRPCPlugin) GRPCServer(broker *plugin.GRPCBroker, s *grpc.Server) error {
	return nil
}
