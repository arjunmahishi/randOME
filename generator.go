package main

import (
	"fmt"
	"math/rand"
	"strings"
)

// generateMetrics generates random metrics data in open metrics format
func generateMetrics(conf *Config) []byte {
	var metrics []string
	for _, m := range conf.Metrics {
		labels := labelCombos(m.Labels)

		maxCardinality := len(labels)
		if m.MaxCardinality > 0 && m.MaxCardinality < maxCardinality {
			maxCardinality = m.MaxCardinality
		}

		for i := 0; i < maxCardinality; i++ {
			labelVals := []string{}
			for key, value := range labels[i] {
				labelVals = append(labelVals, fmt.Sprintf("%s='%s'", key, value))
			}

			metrics = append(metrics, fmt.Sprintf(
				"%s{%s} %d", m.Name, strings.Join(labelVals, ","),
				rand.Intn(m.ValueMax-m.ValueMin)+m.ValueMin,
			))
		}
	}

	return []byte(strings.Join(metrics, "\n"))
}

func labelCombos(labels map[string][]string) []map[string]string {
	var res []map[string]string
	labelComboHelper(labels, []string{}, map[string]string{}, &res)
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

	// Iterate over the label keys in sorted order
	var sortedKeys []string
	for k := range labels {
		sortedKeys = append(sortedKeys, k)
	}

	for _, k := range sortedKeys {
		if _, ok := cur[k]; !ok {
			for _, v := range labels[k] {
				cur[k] = v
				labelComboHelper(labels, append(keys, k), cur, res)
			}
			// Remove the current label from the current combination
			delete(cur, k)
		}
	}
}
