package main

import (
	"io/ioutil"
	"log"

	"gopkg.in/yaml.v2"
)

/*
---
num_points: 1000
frequency: 1
metrics:
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
	Name     string              `yaml:"name"`
	Type     string              `yaml:"type"`
	ValueMin int                 `yaml:"value_min"`
	ValueMax int                 `yaml:"value_max"`
	Labels   map[string][]string `yaml:"labels"`
}

func (c *Config) GetConfig(filename string) *Config {
	yamlFile, err := ioutil.ReadFile(filename)
	if err != nil {
		log.Printf("yamlFile.Get err   #%v ", err)
	}

	if err = yaml.Unmarshal(yamlFile, c); err != nil {
		log.Fatalf("Unmarshal: %v", err)
	}

	return c
}
