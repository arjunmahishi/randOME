package main

type remoteWriter struct{}

func newRemoteWriter() *remoteWriter {
	return &remoteWriter{}
}

func (r *remoteWriter) dumpMetrics(metrics []byte) error {
	return nil
}
