package config

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestFlattenMap(t *testing.T) {

	in := map[string]any{
		"k1": "v1",
		"k2": "v2",
		"k3": map[string]any{
			"k3-1": "v3-1",
			"k3-2": "v3-2",
			"k3-3": map[string]any{
				"k3-3-1": "v3-3-1",
				"k3-3-2": "v3-3-2",
			},
		},
		"k4": map[string]any{
			"k4-1": map[string]any{
				"k4-1-1": "v4-1-1",
			},
		},
	}

	flat := make(map[string]any)
	flattenMap("", in, flat)

	require.Len(t, flat, 7)
	require.Equal(t, flat["k1"], "v1")
	require.Equal(t, flat["k2"], "v2")
	require.Equal(t, flat["k3.k3-1"], "v3-1")
	require.Equal(t, flat["k3.k3-2"], "v3-2")
	require.Equal(t, flat["k3.k3-3.k3-3-1"], "v3-3-1")
	require.Equal(t, flat["k3.k3-3.k3-3-2"], "v3-3-2")
	require.Equal(t, flat["k4.k4-1.k4-1-1"], "v4-1-1")
}
