package main

import (
	"context"
	"time"

	"github.com/m3db/prometheus_remote_client_golang/promremote"
)

type remoteWriter struct {
	promremote.Client
}

func newRemoteWriter() (*remoteWriter, error) {
	promConfig := promremote.NewConfig(
		promremote.WriteURLOption(*addr),
		promremote.HTTPClientTimeoutOption(60*time.Second),
	)

	client, err := promremote.NewClient(promConfig)
	if err != nil {
		return nil, err
	}

	return &remoteWriter{client}, nil
}

func (r *remoteWriter) dumpMetrics(metrics *timeSeries) error {
	_, err := r.WriteTimeSeries(
		context.Background(), metrics.TSList, promremote.WriteOptions{},
	)
	return err
}
