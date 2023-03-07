package main

import (
	"fmt"
)

type stdout struct{}

func (m *stdout) dumpMetrics(metrics *timeSeries) error {
	fmt.Println(metrics)
	return nil
}
