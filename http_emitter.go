package main

import (
	"fmt"
	"log"
	"net/http"
	"sync"
)

type httpEmitter struct {
	body []byte
	sync.RWMutex
}

func newHTTPEmitter() *httpEmitter {
	h := &httpEmitter{}
	go h.startServer()

	return h
}

func (h *httpEmitter) startServer() {
	http.HandleFunc("/metrics", func(w http.ResponseWriter, r *http.Request) {
		h.RLock()
		defer h.RUnlock()

		w.Write(h.body)
		return
	})

	fmt.Printf("serving metrics on: http://localhost:%d/metrics\n", *port)
	log.Fatal(http.ListenAndServe(fmt.Sprintf(":%d", *port), nil))
}

func (h *httpEmitter) dumpMetrics(metrics []byte) error {
	h.Lock()
	defer h.Unlock()

	h.body = metrics
	return nil
}
