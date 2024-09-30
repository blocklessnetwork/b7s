package helpers

import (
	"errors"
	"testing"

	"github.com/prometheus/client_golang/prometheus"
	dto "github.com/prometheus/client_model/go"
	"github.com/stretchr/testify/require"
)

func MetricMap(t *testing.T, g prometheus.Gatherer) map[string]*dto.MetricFamily {
	t.Helper()

	metrics, err := g.Gather()
	require.NoError(t, err)

	// Create a map of gathered metricMap.
	metricMap := make(map[string]*dto.MetricFamily)
	for _, m := range metrics {
		metricMap[*m.Name] = m
	}

	return metricMap
}

func GetMetric(m map[string]*dto.MetricFamily, name string, params ...string) ([]*dto.Metric, error) {

	family, ok := m[name]
	if !ok {
		return nil, errors.New("not found")
	}

	metric := family.GetMetric()
	if len(params) == 0 {
		return metric, nil
	}

	if len(params)%2 != 0 {
		return nil, errors.New("key-value pairs required for label values")
	}

	// Passed in param list.
	labels := make(map[string]string)
	for i := 0; i < len(params); i += 2 {
		labels[params[i]] = params[i+1]
	}

	var out []*dto.Metric
	for _, m := range metric {

		// Labels part of this metric.
		ml := make(map[string]string)
		for _, label := range m.GetLabel() {
			ml[label.GetName()] = label.GetValue()
		}

		match := true
		for k, v := range labels {

			mv, ok := ml[k]
			if !ok {
				match = false
			}

			if mv != v {
				match = false
			}
		}

		if match {
			out = append(out, m)
		}
	}

	return out, nil
}

func CounterCmp(t *testing.T, metrics map[string]*dto.MetricFamily, value float64, name string, params ...string) {
	t.Helper()

	metric, err := GetMetric(metrics, name, params...)
	require.NoError(t, err)
	require.Len(t, metric, 1)
	require.Equal(t, value, metric[0].GetCounter().GetValue())
}
