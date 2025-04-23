package main

import (
	"sort"
	"testing"

	"github.com/m3db/prometheus_remote_client_golang/promremote"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestLabelCombos(t *testing.T) {
	tests := []struct {
		name   string
		labels map[string][]string
		want   int // number of expected combinations
	}{
		{
			name: "Simple case with one label",
			labels: map[string][]string{
				"instance": {"localhost:8080"},
			},
			want: 1,
		},
		{
			name: "Two labels with multiple values",
			labels: map[string][]string{
				"instance": {"localhost:8080"},
				"region":   {"us-east-1", "us-west-1"},
			},
			want: 2,
		},
		{
			name: "Three labels with multiple values",
			labels: map[string][]string{
				"instance": {"localhost:8080"},
				"region":   {"us-east-1", "us-west-1"},
				"cluster":  {"dev", "prod"},
			},
			want: 4,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := labelCombos(tt.labels)
			assert.Equal(t, tt.want, len(result), "Wrong number of combinations generated")

			// Check for duplicate combinations
			seen := make(map[string]bool)
			for _, combo := range result {
				key := mapToString(combo)
				assert.False(t, seen[key], "Duplicate combination found: %s", key)
				seen[key] = true
			}
		})
	}
}

func TestNoDuplicateTimeSeries(t *testing.T) {
	config := &Config{
		Metrics: []Metric{
			{
				Name:     "test_metric",
				ValueMin: 0,
				ValueMax: 100,
				Labels: map[string][]string{
					"instance": {"localhost:8080"},
					"region":   {"us-east-1", "us-west-1"},
					"cluster":  {"dev", "prod"},
				},
			},
		},
	}

	ts := generateMetrics(config)

	// Check for duplicate label sets
	seen := make(map[string]bool)
	for _, series := range ts.TSList {
		key := timeSeriesLabelsToString(series.Labels)
		assert.False(t, seen[key], "Duplicate time series found with labels: %s", key)
		seen[key] = true
	}
}

func TestGenerateMetricsRespectsMaxCardinality(t *testing.T) {
	tests := []struct {
		name            string
		metric          Metric
		wantCardinality int
	}{
		{
			name: "No max cardinality specified",
			metric: Metric{
				Name:     "test_metric",
				ValueMin: 0,
				ValueMax: 100,
				Labels: map[string][]string{
					"instance": {"localhost:8080"},
					"region":   {"us-east-1", "us-west-1"},
					"cluster":  {"dev", "prod"},
				},
				// MaxCardinality not set (0)
			},
			wantCardinality: 4, // 1 instance × 2 regions × 2 clusters
		},
		{
			name: "Max cardinality lower than possible combinations",
			metric: Metric{
				Name:     "test_metric",
				ValueMin: 0,
				ValueMax: 100,
				Labels: map[string][]string{
					"instance": {"localhost:8080"},
					"region":   {"us-east-1", "us-west-1"},
					"cluster":  {"dev", "prod"},
				},
				MaxCardinality: 2,
			},
			wantCardinality: 2,
		},
		{
			name: "Max cardinality higher than possible combinations",
			metric: Metric{
				Name:     "test_metric",
				ValueMin: 0,
				ValueMax: 100,
				Labels: map[string][]string{
					"instance": {"localhost:8080"},
					"region":   {"us-east-1", "us-west-1"},
				},
				MaxCardinality: 10,
			},
			wantCardinality: 2, // 1 instance × 2 regions
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			config := &Config{
				Metrics: []Metric{tt.metric},
			}

			ts := generateMetrics(config)

			assert.Equal(t, tt.wantCardinality, len(ts.TSList), 
				"Wrong number of time series generated")

			// Check for duplicate label sets
			seen := make(map[string]bool)
			for _, series := range ts.TSList {
				key := timeSeriesLabelsToString(series.Labels)
				assert.False(t, seen[key], "Duplicate time series found with labels: %s", key)
				seen[key] = true
			}
		})
	}
}

func TestFixedValueRespected(t *testing.T) {
	fixedValue := 42.0
	config := &Config{
		Metrics: []Metric{
			{
				Name:     "test_metric",
				ValueMin: 0,
				ValueMax: 100,
				Value:    &fixedValue,
				Labels: map[string][]string{
					"instance": {"localhost:8080"},
				},
			},
		},
	}

	ts := generateMetrics(config)
	require.NotEmpty(t, ts.TSList, "No time series generated")

	// Check that all time series have the fixed value
	for _, series := range ts.TSList {
		assert.Equal(t, fixedValue, series.Datapoint.Value, 
			"Time series does not have the expected fixed value")
	}
}

// Helper function to convert a map to a consistent string representation
func mapToString(m map[string]string) string {
	// Sort keys for consistent output
	keys := make([]string, 0, len(m))
	for k := range m {
		keys = append(keys, k)
	}
	sort.Strings(keys)

	// Build string representation
	result := ""
	for _, k := range keys {
		result += k + "=" + m[k] + ";"
	}
	return result
}

// Helper function to convert time series labels to a consistent string representation
func timeSeriesLabelsToString(labels []promremote.Label) string {
	// Convert to map for easier sorting
	labelsMap := make(map[string]string)
	metricName := ""
	
	for _, label := range labels {
		if label.Name == "__name__" {
			metricName = label.Value
			continue
		}
		labelsMap[label.Name] = label.Value
	}
	
	return metricName + "{" + mapToString(labelsMap) + "}"
}