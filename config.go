package main

import (
	"io"
	"io/ioutil"

	"gopkg.in/yaml.v2"
)

/*
---
num_points: 1000
frequency: 1
metrics:
	- name: cpu_usage
		type: gauge
		value_min: 0
		value_max: 100
		labels:
			instance: [localhost:8080]
			cluster: [dev, prod, staging, test, qa]

  - name: requests_total
    type: counter
    labels:
      method: [GET, POST, PUT, DELETE]
      status: [200, 400, 404, 500]

  - name: response_time_seconds
    type: histogram
    labels:
      method: [GET, POST, PUT, DELETE]
*/

type Config struct {
	NumPoints int      `yaml:"num_points"`
	Frequency int      `yaml:"frequency"`
	Metrics   []Metric `yaml:"metrics"`
}

type Metric struct {
	Name           string              `yaml:"name"`
	Type           string              `yaml:"type"`
	ValueMin       int                 `yaml:"value_min"`
	ValueMax       int                 `yaml:"value_max"`
	Labels         map[string][]string `yaml:"labels"`
	MaxCardinality int                 `yaml:"max_cardinality"`
}

var (
	defaultConfig = &Config{
		Frequency: 1,
		Metrics: []Metric{
			{
				Name:     "cpu_usage",
				Type:     "gauge",
				ValueMin: 0,
				ValueMax: 100,
				Labels: map[string][]string{
					"instance": {"localhost:8080"},
					"cluster":  {"dev", "prod", "staging", "test", "qa"},
					"region":   {"us-east-1", "us-west-1", "us-west-2", "eu-west-1", "eu-central-1"},
					"service":  {"api", "web", "db", "cache", "queue", "worker"},
				},
				MaxCardinality: 10,
			},
		},
	}
)

func GetConfig(confReader io.Reader) (*Config, error) {
	yamlFile, err := ioutil.ReadAll(confReader)
	if err != nil {
		return nil, err
	}

	var conf Config
	if err = yaml.Unmarshal(yamlFile, &conf); err != nil {
		return nil, err
	}

	return &conf, nil
}
