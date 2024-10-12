package main

import (
	"slices"

	mp "github.com/armon/go-metrics/prometheus"

	"github.com/blocklessnetwork/b7s/consensus/pbft"
	"github.com/blocklessnetwork/b7s/consensus/raft"
	"github.com/blocklessnetwork/b7s/execution/executor"
	"github.com/blocklessnetwork/b7s/fstore"
	"github.com/blocklessnetwork/b7s/host"
	"github.com/blocklessnetwork/b7s/node"
)

func metricCounters() []mp.CounterDefinition {

	counters := slices.Concat(
		node.Counters,
		host.Counters,
		fstore.Counters,
		executor.Counters,
	)

	return counters
}

func metricSummaries() []mp.SummaryDefinition {

	summaries := slices.Concat(
		executor.Summaries,
		fstore.Summaries,
		pbft.Summaries,
		raft.Summaries,
	)

	return summaries
}

func metricGauges() []mp.GaugeDefinition {

	// Right now we have a single gauge - node info.
	gauges := node.Gauges
	return gauges
}
