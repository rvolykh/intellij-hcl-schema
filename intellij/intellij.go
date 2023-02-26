package intellij

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"

	"github.com/rvolykh/intellij-hcl-schema/proto"
)

const (
	providersHclPath = ".terraform.d/metadata-repo/terraform/model/providers"
)

func Export(provider *proto.GetProviderSchema_Response) *ProviderSchema {
	result := &ProviderSchema{
		Type:        "provider",
		Provider:    SchemaInfo{},
		Resources:   map[string]SchemaInfo{},
		DataSources: map[string]SchemaInfo{},
	}

	for _, v := range provider.Provider.Block.Attributes {
		result.Provider[v.Name] = SchemaDefinition{
			Type:        typeValToStr(v.Type),
			Optional:    v.Optional,
			Required:    v.Required,
			Description: v.Description,
			Computed:    v.Computed,
			Deprecated:  deprecatedValToStr(v.Deprecated),
			Sensitive:   v.Sensitive,
		}
	}

	for k, v := range provider.DataSourceSchemas {
		result.DataSources[k] = SchemaInfo{}
		for _, av := range v.Block.Attributes {
			result.DataSources[k][av.Name] = SchemaDefinition{
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

	for k, v := range provider.ResourceSchemas {
		result.Resources[k] = SchemaInfo{}
		for _, av := range v.Block.Attributes {
			result.Resources[k][av.Name] = SchemaDefinition{
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

	return result
}

func LoadToIDE(provider *ProviderSchema) (err error) {
	homePath, err := os.UserHomeDir()
	if err != nil {
		return fmt.Errorf("find HOME path: %w", err)
	}

	hclPath := filepath.Join(homePath, providersHclPath)
	if err := os.MkdirAll(hclPath, os.ModeDir); err != nil {
		return fmt.Errorf("ensure path %s: %w", hclPath, err)
	}

	filename := filepath.Join(hclPath, fmt.Sprintf("%s.json", provider.Name))
	file, err := os.OpenFile(filename, os.O_CREATE|os.O_TRUNC|os.O_WRONLY, os.ModePerm)
	if err != nil {
		return fmt.Errorf("create %s: %w", filename, err)
	}
	defer func() {
		if errClose := file.Close(); errClose != nil {
			err = fmt.Errorf("close file %s: %w", filename, errClose)
		}
	}()

	encoder := json.NewEncoder(file)
	encoder.SetIndent("", "  ")
	if err := encoder.Encode(provider); err != nil {
		return fmt.Errorf("json encode schema: %w", err)
	}

	return nil
}

// typeValToStr converts GRPC response Type value from bytes to string with removing open/close quotes
func typeValToStr(v []byte) string {
	return string(v[1 : len(v)-1])
}

// deprecatedValToStr converts GRPC response Deprecated value from bool to string
func deprecatedValToStr(v bool) string {
	if !v {
		return ""
	}
	return "Deprecated! Please refer to documentation for more details"
}
