package traceable

import (
	"context"

	"github.com/libp2p/go-libp2p/core/peer"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/trace"

	"github.com/blocklessnetwork/b7s/models/blockless"
	"github.com/blocklessnetwork/b7s/store"
	"github.com/blocklessnetwork/b7s/telemetry/b7ssemconv"
	"github.com/blocklessnetwork/b7s/telemetry/tracing"
)

// Store is a thin wrapper around the standard b7s store, adding a tracer to it.
type Store struct {
	store  *store.Store
	tracer *tracing.Tracer
}

func New(store *store.Store) *Store {

	s := Store{
		store:  store,
		tracer: tracing.NewTracer(tracerName),
	}

	return &s
}

func (s *Store) SavePeer(ctx context.Context, peer blockless.Peer) error {

	callback := func() error {
		return s.store.SavePeer(ctx, peer)
	}

	opts := storeSpanOptions(tracing.SpanAttributes(peerAttributes(peer))...)
	return s.tracer.WithSpanFromContext(ctx, "SavePeer", callback, opts...)
}

func (s *Store) SaveFunction(ctx context.Context, function blockless.FunctionRecord) error {

	callback := func() error {
		return s.store.SaveFunction(ctx, function)
	}

	opts := storeSpanOptions(trace.WithAttributes(b7ssemconv.FunctionCID.String(function.CID)))
	return s.tracer.WithSpanFromContext(ctx, "SaveFunction", callback, opts...)
}

func (s *Store) RetrievePeer(ctx context.Context, id peer.ID) (blockless.Peer, error) {

	var peer blockless.Peer
	var err error
	callback := func() error {
		peer, err = s.store.RetrievePeer(ctx, id)
		return err
	}

	opts := storeSpanOptions(trace.WithAttributes(b7ssemconv.PeerID.String(id.String())))
	_ = s.tracer.WithSpanFromContext(ctx, "GetPeer", callback, opts...)
	return peer, err
}

func (s *Store) RetrievePeers(ctx context.Context) ([]blockless.Peer, error) {

	var peers []blockless.Peer
	var err error
	callback := func() error {
		peers, err = s.store.RetrievePeers(ctx)
		return err
	}

	_ = s.tracer.WithSpanFromContext(ctx, "ListPeers", callback, storeSpanOptions()...)
	return peers, err
}

func (s *Store) RetrieveFunction(ctx context.Context, cid string) (blockless.FunctionRecord, error) {

	var function blockless.FunctionRecord
	var err error
	callback := func() error {
		function, err = s.store.RetrieveFunction(ctx, cid)
		return err
	}

	opts := storeSpanOptions(trace.WithAttributes(b7ssemconv.FunctionCID.String(cid)))
	_ = s.tracer.WithSpanFromContext(ctx, "GetFunction", callback, opts...)
	return function, err
}

func (s *Store) RetrieveFunctions(ctx context.Context) ([]blockless.FunctionRecord, error) {

	var functions []blockless.FunctionRecord
	var err error
	callback := func() error {
		functions, err = s.store.RetrieveFunctions(ctx)
		return err
	}

	_ = s.tracer.WithSpanFromContext(ctx, "ListFunctions", callback, storeSpanOptions()...)
	return functions, err
}

func (s *Store) RemovePeer(ctx context.Context, id peer.ID) error {

	opts := storeSpanOptions(trace.WithAttributes(b7ssemconv.PeerID.String(id.String())))
	return s.tracer.WithSpanFromContext(
		ctx,
		"RemovePeer",
		func() error { return s.store.RemovePeer(ctx, id) },
		opts...)
}

func (s *Store) RemoveFunction(ctx context.Context, cid string) error {

	opts := storeSpanOptions(trace.WithAttributes(b7ssemconv.FunctionCID.String(cid)))
	return s.tracer.WithSpanFromContext(
		ctx,
		"RemoveFunction",
		func() error { return s.store.RemoveFunction(ctx, cid) },
		opts...)
}

func peerAttributes(peer blockless.Peer) []attribute.KeyValue {
	return []attribute.KeyValue{
		b7ssemconv.PeerID.String(peer.ID.String()),
		b7ssemconv.PeerMultiaddr.String(peer.MultiAddr),
	}
}
