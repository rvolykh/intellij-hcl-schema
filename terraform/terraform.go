package terraform

import (
	"context"
	"errors"

	"github.com/rvolykh/intellij-hcl-schema/model"
)

const (
	// Maximum response size for schema request, 256MB
	maxRecvSize = 256 << 20
	// Number of protocol versions supported, must match number of goroutines
	protocolsCount = 2
)

type (
	result struct {
		err    error
		schema *model.ProviderSchema
	}

	getProviderSchemaFunc func(ctx context.Context, pluginPath string) (*model.ProviderSchema, error)
)

func GetProviderSchema(ctx context.Context, pluginPath string) (*model.ProviderSchema, error) {
	var resultCh = make(chan result, protocolsCount)
	defer close(resultCh)

	run := func(fn getProviderSchemaFunc) {
		schema, err := fn(ctx, pluginPath)
		resultCh <- result{err: err, schema: schema}
	}

	// must match protocolsCount
	go run(getProviderSchemaV5)
	go run(getProviderSchemaV6)

	var err = errors.New("multiple errors occurred")
	for received := 0; received < protocolsCount; received++ {
		res := <-resultCh
		if res.err != nil {
			err = errors.Join(err, res.err)
			continue
		}

		return res.schema, nil
	}

	return nil, err
}
