package intellij

import (
	"bufio"
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"

	"github.com/rvolykh/intellij-hcl-schema/model"
)

const (
	providersHclPath = ".terraform.d/metadata-repo/terraform/model/providers"
)

func LoadToIDE(provider *model.ProviderSchema) (err error) {
	homePath, err := os.UserHomeDir()
	if err != nil {
		return fmt.Errorf("find HOME path: %w", err)
	}

	hclPath := filepath.Join(homePath, providersHclPath)
	if err := os.MkdirAll(hclPath, 0755); err != nil {
		return fmt.Errorf("ensure path %s: %w", hclPath, err)
	}

	filename := filepath.Join(hclPath, fmt.Sprintf("%s.json", provider.Name))
	file, err := os.OpenFile(filename, os.O_CREATE|os.O_TRUNC|os.O_WRONLY, 0644)
	if err != nil {
		return fmt.Errorf("create %s: %w", filename, err)
	}
	defer func() {
		if errClose := file.Close(); errClose != nil && err == nil {
			err = fmt.Errorf("close file %s: %w", filename, errClose)
		}
	}()

	bufferedWriter := bufio.NewWriter(file)
	defer func() {
		if errFlush := bufferedWriter.Flush(); errFlush != nil && err == nil {
			err = fmt.Errorf("flush buffer for %s: %w", filename, errFlush)
		}
	}()

	encoder := json.NewEncoder(bufferedWriter)
	encoder.SetIndent("", "  ")
	if err := encoder.Encode(provider); err != nil {
		return fmt.Errorf("json encode schema: %w", err)
	}

	return nil
}
