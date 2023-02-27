package main

import "fmt"

type stdout struct{}

func (m *stdout) dumpMetrics(metrics []byte) error {
	fmt.Println(string(metrics))
	return nil
}
