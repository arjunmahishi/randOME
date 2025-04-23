package main

import (
	"fmt"
	"math/rand"
	"sort"
	"strings"
	"sync"
	"time"

	"github.com/m3db/prometheus_remote_client_golang/promremote"
)

type timeSeries struct {
	promremote.TSList
}

func (ts *timeSeries) String() string {
	lines := []string{}
	for _, row := range ts.TSList {
		var (
			metricName = ""
			labelVals  = []string{}
		)

		for _, label := range row.Labels {
			if label.Name == "__name__" {
				metricName = label.Value
				continue
			}

			labelVals = append(labelVals, fmt.Sprintf("%s=\"%s\"", label.Name, label.Value))
		}

		lines = append(lines, fmt.Sprintf(
			"%s{%s} %d",
			metricName, strings.Join(labelVals, ","), int64(row.Datapoint.Value),
		))
	}

	return strings.Join(lines, "\n")
}

// generateMetrics generates random metrics data in open metrics format
func generateMetrics(conf *Config) *timeSeries {
	var (
		ts = promremote.TSList{}
		wg = sync.WaitGroup{}
	)

	wg.Add(len(conf.Metrics))
	for _, m := range conf.Metrics {
		go func(mt Metric) {
			defer wg.Done()

			labels := labelCombos(mt.Labels)

			maxCardinality := len(labels)
			if mt.MaxCardinality > 0 && mt.MaxCardinality < maxCardinality {
				maxCardinality = mt.MaxCardinality
			}

			for i := 0; i < maxCardinality; i++ {
				currLabels := []promremote.Label{
					{
						Name:  "__name__",
						Value: mt.Name,
					},
				}

				for key, value := range labels[i] {
					currLabels = append(currLabels, promremote.Label{
						Name:  key,
						Value: value,
					})
				}

				randIn := mt.ValueMax - mt.ValueMin
				if randIn <= 0 {
					randIn = 10
				}

				val := float64(rand.Intn(randIn) + mt.ValueMin)
				if mt.Value != nil {
					val = *mt.Value
				}

				ts = append(ts, promremote.TimeSeries{
					Labels: currLabels,
					Datapoint: promremote.Datapoint{
						Timestamp: time.Now(),
						Value:     val,
					},
				})
			}
		}(m)
	}

	wg.Wait()
	return &timeSeries{ts}
}

func labelCombos(labels map[string][]string) []map[string]string {
	var res []map[string]string
	
	// Create a sorted list of keys for consistent ordering
	var sortedKeys []string
	for k := range labels {
		sortedKeys = append(sortedKeys, k)
	}
	
	// Sort keys (optional but helps with consistency)
	sort.Strings(sortedKeys)
	
	labelComboHelper(labels, sortedKeys, map[string]string{}, &res)
	return res
}

func labelComboHelper(labels map[string][]string, keys []string, cur map[string]string, res *[]map[string]string) {
	if len(cur) == len(labels) {
		// Make a copy of the current combination before appending it to the result
		copy := make(map[string]string)
		for k, v := range cur {
			copy[k] = v
		}
		*res = append(*res, copy)
		return
	}

	// Use the keys parameter that tracks the order of processing
	// If no keys are provided yet, initialize with all label keys
	if len(keys) == 0 {
		for k := range labels {
			if _, ok := cur[k]; !ok {
				keys = append(keys, k)
			}
		}
	}

	if len(keys) > 0 {
		k := keys[0]
		remainingKeys := keys[1:]
		
		for _, v := range labels[k] {
			cur[k] = v
			labelComboHelper(labels, remainingKeys, cur, res)
		}
		// Remove the current label from the current combination
		delete(cur, k)
	}
}
