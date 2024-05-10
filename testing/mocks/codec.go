package mocks

import (
	"testing"
)

type Codec struct {
	MarshalFunc   func(any) ([]byte, error)
	UnmarshalFunc func([]byte, any) error
}

func BaselineCodec(t *testing.T) *Codec {
	t.Helper()

	codec := Codec{
		MarshalFunc: func(any) ([]byte, error) {
			return []byte{}, nil
		},
		UnmarshalFunc: func([]byte, any) error {
			return nil
		},
	}

	return &codec
}

func (c *Codec) Marshal(obj any) ([]byte, error) {
	return c.MarshalFunc(obj)
}

func (c *Codec) Unmarshal(data []byte, obj any) error {
	return c.UnmarshalFunc(data, obj)
}
