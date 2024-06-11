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

func (s *Store) SavePeer(peer blockless.Peer) error {

	callback := func() error {
		return s.store.SavePeer(peer)
	}

	opts := storeSpanOptions(tracing.SpanAttributes(peerAttributes(peer))...)
	return s.tracer.WithSpanFromContext(context.Background(), "SavePeer", callback, opts...)
}

func (s *Store) SaveFunction(function blockless.FunctionRecord) error {

	callback := func() error {
		return s.store.SaveFunction(function)
	}

	opts := storeSpanOptions(trace.WithAttributes(b7ssemconv.FunctionCID.String(function.CID)))
	return s.tracer.WithSpanFromContext(context.Background(), "SaveFunction", callback, opts...)
}

func (s *Store) RetrievePeer(id peer.ID) (blockless.Peer, error) {

	var peer blockless.Peer
	var err error
	callback := func() error {
		peer, err = s.store.RetrievePeer(id)
		return err
	}

	opts := storeSpanOptions(trace.WithAttributes(b7ssemconv.PeerID.String(id.String())))
	_ = s.tracer.WithSpanFromContext(context.Background(), "GetPeer", callback, opts...)
	return peer, err
}

func (s *Store) RetrievePeers() ([]blockless.Peer, error) {

	var peers []blockless.Peer
	var err error
	callback := func() error {
		peers, err = s.store.RetrievePeers()
		return err
	}

	_ = s.tracer.WithSpanFromContext(context.Background(), "ListPeers", callback, storeSpanOptions()...)
	return peers, err
}

func (s *Store) RetrieveFunction(cid string) (blockless.FunctionRecord, error) {

	var function blockless.FunctionRecord
	var err error
	callback := func() error {
		function, err = s.store.RetrieveFunction(cid)
		return err
	}

	opts := storeSpanOptions(trace.WithAttributes(b7ssemconv.FunctionCID.String(cid)))
	_ = s.tracer.WithSpanFromContext(context.Background(), "GetFunction", callback, opts...)
	return function, err
}

func (s *Store) RetrieveFunctions() ([]blockless.FunctionRecord, error) {

	var functions []blockless.FunctionRecord
	var err error
	callback := func() error {
		functions, err = s.store.RetrieveFunctions()
		return err
	}

	_ = s.tracer.WithSpanFromContext(context.Background(), "ListFunctions", callback, storeSpanOptions()...)
	return functions, err
}

func (s *Store) RemovePeer(id peer.ID) error {

	opts := storeSpanOptions(trace.WithAttributes(b7ssemconv.PeerID.String(id.String())))
	return s.tracer.WithSpanFromContext(
		context.Background(),
		"RemovePeer",
		func() error { return s.store.RemovePeer(id) },
		opts...)
}

func (s *Store) RemoveFunction(cid string) error {

	opts := storeSpanOptions(trace.WithAttributes(b7ssemconv.FunctionCID.String(cid)))
	return s.tracer.WithSpanFromContext(
		context.Background(),
		"RemoveFunction",
		func() error { return s.store.RemoveFunction(cid) },
		opts...)
}

func peerAttributes(peer blockless.Peer) []attribute.KeyValue {
	return []attribute.KeyValue{
		b7ssemconv.PeerID.String(peer.ID.String()),
		b7ssemconv.PeerMultiaddr.String(peer.MultiAddr),
	}
}
