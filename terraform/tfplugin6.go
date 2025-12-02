package terraform

import (
	"context"
	"fmt"
	"os/exec"

	"github.com/hashicorp/go-hclog"
	"github.com/hashicorp/go-plugin"
	"github.com/rvolykh/intellij-hcl-schema/model"
	"github.com/rvolykh/intellij-hcl-schema/proto/tfplugin6"
	"google.golang.org/grpc"
)

func getProviderSchemaV6(ctx context.Context, pluginPath string) (*model.ProviderSchema, error) {
	handshake := plugin.HandshakeConfig{
		ProtocolVersion:  6,
		MagicCookieKey:   "TF_PLUGIN_MAGIC_COOKIE",
		MagicCookieValue: "d602bf8f470bc67ca7faa0386276bbdd4330efaf76d1a219cb4d6991ca9872b2",
	}

	pluginClient := plugin.NewClient(
		&plugin.ClientConfig{
			Cmd:              exec.CommandContext(ctx, pluginPath),
			HandshakeConfig:  handshake,
			AllowedProtocols: []plugin.Protocol{plugin.ProtocolGRPC},
			Managed:          true,
			Logger:           hclog.NewNullLogger(),
			Plugins: map[string]plugin.Plugin{
				"provider": &grpcPluginV6{},
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

	grpcProviderClient := rawProviderClient.(tfplugin6.ProviderClient)

	request := new(tfplugin6.GetProviderSchema_Request)
	options := []grpc.CallOption{
		grpc.MaxRecvMsgSizeCallOption{MaxRecvMsgSize: maxRecvSize},
	}

	resp, err := grpcProviderClient.GetProviderSchema(ctx, request, options...)
	if err != nil {
		return nil, fmt.Errorf("failed to get provider schema: %w", err)
	}

	return parseV6(resp)
}

type grpcPluginV6 struct {
	plugin.Plugin
}

func (p *grpcPluginV6) GRPCClient(ctx context.Context, broker *plugin.GRPCBroker, c *grpc.ClientConn) (interface{}, error) {
	return tfplugin6.NewProviderClient(c), nil
}
func (p *grpcPluginV6) GRPCServer(broker *plugin.GRPCBroker, s *grpc.Server) error {
	return nil
}

func parseV6(schema *tfplugin6.GetProviderSchema_Response) (*model.ProviderSchema, error) {
	result := &model.ProviderSchema{
		Type:        "provider",
		Provider:    model.SchemaInfo{},
		Resources:   map[string]model.SchemaInfo{},
		DataSources: map[string]model.SchemaInfo{},
	}

	for _, v := range schema.Provider.Block.Attributes {
		result.Provider[v.Name] = model.SchemaDefinition{
			Type:        typeValToStr(v.Type),
			Optional:    v.Optional,
			Required:    v.Required,
			Description: v.Description,
			Computed:    v.Computed,
			Deprecated:  deprecatedValToStr(v.Deprecated),
			Sensitive:   v.Sensitive,
		}
	}

	for k, v := range schema.DataSourceSchemas {
		result.DataSources[k] = model.SchemaInfo{}
		for _, av := range v.Block.Attributes {
			result.DataSources[k][av.Name] = model.SchemaDefinition{
				Type:        typeValToStr(av.Type),
				Description: av.Description,
				Optional:    av.Optional,
				Required:    av.Required,
				Computed:    av.Computed,
				Deprecated:  deprecatedValToStr(av.Deprecated),
				Sensitive:   av.Sensitive,
			}
		}
	}

	for k, v := range schema.ResourceSchemas {
		result.Resources[k] = model.SchemaInfo{}
		for _, av := range v.Block.Attributes {
			result.Resources[k][av.Name] = model.SchemaDefinition{
				Type:        typeValToStr(av.Type),
				Description: av.Description,
				Deprecated:  deprecatedValToStr(av.Deprecated),
				Optional:    av.Optional,
				Required:    av.Required,
				Computed:    av.Computed,
				Sensitive:   av.Sensitive,
			}
		}
	}

	return result, nil
}
