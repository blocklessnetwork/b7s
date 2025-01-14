package telemetry

import (
	"context"
	"errors"
	"fmt"

	"go.opentelemetry.io/otel/propagation"
	"go.opentelemetry.io/otel/sdk/resource"
	semconv "go.opentelemetry.io/otel/semconv/v1.4.0"

	"github.com/blessnetwork/b7s/models/bls"
	"github.com/blessnetwork/b7s/telemetry/b7ssemconv"
)

func CreateResource(ctx context.Context, id string, role bls.NodeRole) (*resource.Resource, error) {

	if id == "" {
		return nil, errors.New("instance ID is required")
	}

	if !role.Valid() {
		return nil, errors.New("invalid node role")
	}

	opts := append(
		defaultResourceOpts,
		resource.WithAttributes(semconv.ServiceInstanceIDKey.String(id)),
		resource.WithAttributes(b7ssemconv.ServiceRole.String(role.String())),
	)

	resource, err := resource.New(ctx, opts...)
	if err != nil {
		return nil, fmt.Errorf("could not initialize resource: %w", err)
	}

	return resource, nil
}

func CreatePropagator() propagation.TextMapPropagator {
	return propagation.NewCompositeTextMapPropagator(
		propagation.TraceContext{},
		propagation.Baggage{},
	)
}
