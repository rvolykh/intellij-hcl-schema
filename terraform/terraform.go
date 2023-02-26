package terraform

import (
	"context"
	"fmt"
	"os/exec"

	"github.com/hashicorp/go-hclog"
	"github.com/hashicorp/go-plugin"
	"github.com/rvolykh/intellij-hcl-schema/proto"
	"google.golang.org/grpc"
)

func GetProviderSchema(ctx context.Context, pluginPath string) (*proto.GetProviderSchema_Response, error) {
	handshake := plugin.HandshakeConfig{
		ProtocolVersion:  5,
		MagicCookieKey:   "TF_PLUGIN_MAGIC_COOKIE",
		MagicCookieValue: "d602bf8f470bc67ca7faa0386276bbdd4330efaf76d1a219cb4d6991ca9872b2",
	}

	pluginClient := plugin.NewClient(
		&plugin.ClientConfig{
			Cmd:              exec.Command(pluginPath),
			HandshakeConfig:  handshake,
			AllowedProtocols: []plugin.Protocol{plugin.ProtocolGRPC},
			Managed:          true,
			Logger:           hclog.NewNullLogger(),
			Plugins: map[string]plugin.Plugin{
				"provider": &GRPCPlugin{},
			},
		},
	)
	defer pluginClient.Kill()

	protocolClient, err := pluginClient.Client()
	if err != nil {
		return nil, fmt.Errorf("failed to start provider: %w", err)
	}

	rawProviderClient, err := protocolClient.Dispense("provider")
	if err != nil {
		return nil, fmt.Errorf("failed to setup client for provider: %w", err)
	}

	grpcProviderClient := rawProviderClient.(proto.ProviderClient)

	request := new(proto.GetProviderSchema_Request)
	options := []grpc.CallOption{
		grpc.MaxRecvMsgSizeCallOption{MaxRecvMsgSize: maxRecvSize},
	}

	resp, err := grpcProviderClient.GetSchema(ctx, request, options...)
	if err != nil {
		return nil, fmt.Errorf("failed to get provider schema: %w", err)
	}

	return resp, nil
}
