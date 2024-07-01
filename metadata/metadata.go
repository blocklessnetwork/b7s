package metadata

import (
	"github.com/blocklessnetwork/b7s/models/execute"
)

type Provider interface {
	Metadata(execute.Request, execute.RuntimeOutput) (any, error)
}

type noopProvider struct{}

func (p noopProvider) Metadata(execute.Request, execute.RuntimeOutput) (any, error) {
	return map[string]any{}, nil
}

func NewNoopProvider() Provider {
	return noopProvider{}
}
