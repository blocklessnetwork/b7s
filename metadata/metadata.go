package metadata

import (
	"github.com/blocklessnetwork/b7s/models/execute"
)

type Metadata map[string]any

type Provider interface {
	Metadata(execute.Request, execute.RuntimeOutput) (Metadata, error)
}

type noopProvider struct{}

func (p noopProvider) Metadata(execute.Request, execute.RuntimeOutput) (Metadata, error) {
	return map[string]any{}, nil
}

func NewNoopProvider() Provider {
	return noopProvider{}
}
