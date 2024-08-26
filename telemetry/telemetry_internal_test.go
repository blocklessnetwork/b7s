package telemetry

import (
	"context"
	"testing"

	"github.com/stretchr/testify/require"
	"go.opentelemetry.io/otel/attribute"

	"github.com/blocklessnetwork/b7s/models/blockless"
)

func TestTelemetry_Resource(t *testing.T) {

	var (
		id   = "resource-id"
		role = blockless.WorkerNode
	)

	resource, err := CreateResource(context.Background(), id, role)
	require.NoError(t, err)

	// Convert attributes to a map.
	attrs := make(map[string]attribute.Value)
	for _, attr := range resource.Attributes() {
		attrs[string(attr.Key)] = attr.Value
	}

	// These are values returned by Otel, so we will not verify their correctness as it would be an overkill.

	// Verify existence of OS attributes.
	require.NotEmpty(t, attrs["os.type"].AsString())
	require.NotEmpty(t, attrs["os.description"].AsString())

	// Verify existence of process attributes.
	require.NotEmpty(t, attrs["process.pid"].AsInt64())
	require.NotEmpty(t, attrs["process.executable.name"].AsString())
	require.NotEmpty(t, attrs["process.executable.path"].AsString())

	// Verify telemetry attributes.
	require.Equal(t, "opentelemetry", attrs["telemetry.sdk.name"].AsString())

	// Verify existence of service attributes.
	require.Equal(t, id, attrs["service.instance.id"].AsString())
	require.Equal(t, "b7s", attrs["service.name"].AsString())
	require.Equal(t, role.String(), attrs["service.role"].AsString())
}
