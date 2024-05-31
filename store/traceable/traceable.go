package traceable

import (
	"context"

	"github.com/libp2p/go-libp2p/core/peer"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/codes"
	"go.opentelemetry.io/otel/trace"

	"github.com/blocklessnetwork/b7s/models/blockless"
	"github.com/blocklessnetwork/b7s/store"
	"github.com/blocklessnetwork/b7s/telemetry/b7ssemconv"
)

// Store is a thin wrapper around the standard b7s store, adding tracer to it.
type Store struct {
	store  *store.Store
	tracer trace.Tracer
}

func New(store *store.Store) *Store {

	s := Store{
		store:  store,
		tracer: otel.Tracer(tracerName),
	}

	return &s
}

func (s *Store) SavePeer(peer blockless.Peer) error {

	opts := traceOptions
	opts = append(opts,
		trace.WithAttributes(
			b7ssemconv.PeerID.String(peer.ID.String()),
			b7ssemconv.PeerMultiaddr.String(peer.MultiAddr),
		),
	)
	_, span := s.tracer.Start(context.Background(), "save peer", opts...)
	defer span.End()

	err := s.store.SavePeer(peer)
	setSpanStatus(span, err)
	return err
}

func (s *Store) SaveFunction(function blockless.FunctionRecord) error {

	opts := traceOptions
	opts = append(opts,
		trace.WithAttributes(
			b7ssemconv.FunctionCID.String(function.CID),
			// TODO: Perhaps more details for function?
		),
	)
	_, span := s.tracer.Start(context.Background(), "save function", opts...)
	defer span.End()

	err := s.store.SaveFunction(function)
	setSpanStatus(span, err)
	return err
}

func (s *Store) RetrievePeer(id peer.ID) (blockless.Peer, error) {

	opts := traceOptions
	opts = append(opts, trace.WithAttributes(
		b7ssemconv.PeerID.String(id.String()),
	))
	_, span := s.tracer.Start(context.Background(), "get peer", opts...)
	defer span.End()

	peer, err := s.store.RetrievePeer(id)
	setSpanStatus(span, err)
	return peer, err
}

func (s *Store) RetrievePeers() ([]blockless.Peer, error) {

	_, span := s.tracer.Start(context.Background(), "list peers", traceOptions...)
	defer span.End()

	peers, err := s.store.RetrievePeers()
	setSpanStatus(span, err)
	return peers, err
}

func (s *Store) RetrieveFunction(cid string) (blockless.FunctionRecord, error) {

	opts := traceOptions
	opts = append(opts, trace.WithAttributes(
		b7ssemconv.FunctionCID.String(cid),
	))
	_, span := s.tracer.Start(context.Background(), "get function", opts...)
	defer span.End()

	function, err := s.store.RetrieveFunction(cid)
	setSpanStatus(span, err)
	return function, err
}

func (s *Store) RetrieveFunctions() ([]blockless.FunctionRecord, error) {

	_, span := s.tracer.Start(context.Background(), "get function", traceOptions...)
	defer span.End()

	functions, err := s.store.RetrieveFunctions()
	setSpanStatus(span, err)
	return functions, err
}

func (s *Store) RemovePeer(id peer.ID) error {

	opts := traceOptions
	opts = append(opts, trace.WithAttributes(
		b7ssemconv.PeerID.String(id.String()),
	))
	_, span := s.tracer.Start(context.Background(), "remove peer", opts...)
	defer span.End()

	err := s.store.RemovePeer(id)
	setSpanStatus(span, err)
	return err
}

func (s *Store) RemoveFunction(cid string) error {

	opts := traceOptions
	opts = append(opts, trace.WithAttributes(
		b7ssemconv.FunctionCID.String(cid),
	))
	_, span := s.tracer.Start(context.Background(), "remove function", opts...)
	defer span.End()

	err := s.store.RemoveFunction(cid)
	setSpanStatus(span, err)
	return err
}

func setSpanStatus(span trace.Span, err error) {
	if err != nil {
		span.SetStatus(codes.Error, err.Error())
		return
	}

	span.SetStatus(codes.Ok, "")
}
