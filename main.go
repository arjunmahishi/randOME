package main

import (
	"math/rand"
	"time"

	"gopkg.in/alecthomas/kingpin.v2"
)

var (
	// commands
	printData   = kingpin.Command("print", "print the data to stdout")
	remoteWrite = kingpin.Command("remote-write", "Remote write to a Prometheus-compatible TSDB")
	emitMetrics = kingpin.Command("emit", "Emit metrics over HTTP on a given port (localhost:<port>/metrics)")

	// flags
	frequency = kingpin.Flag("frequency", "Frequency of data points").Short('f').Default("4s").Duration()
	conf      = kingpin.Flag("config", "path to the config file").Short('c').File()
	addr      = remoteWrite.Flag("addr", "the HTTP address of the TSDB. Should include the basic auth info if required").Required().URL()
	port      = emitMetrics.Flag("port", "HTTP port for emitting metrics").Default("9090").Int()
)

const (
	// Version to be overridden at build time
	Version = "0.0.1"
)

type metricDumper interface {
	dumpMetrics([]byte) error
}

func main() {
	kingpin.Version(Version)
	kingpin.Parse()
	rand.Seed(time.Now().UnixNano())

	config := defaultConfig
	if *conf != nil {
		var err error
		config, err = GetConfig(*conf)
		if err != nil {
			kingpin.Fatalf("error parsing config: %v", err)
		}
	}

	var dumper metricDumper
	switch kingpin.Parse() {
	case printData.FullCommand():
		dumper = &stdout{}
	case remoteWrite.FullCommand():
		dumper = newRemoteWriter()
	case emitMetrics.FullCommand():
		dumper = newHTTPEmitter()
	default:
		kingpin.Usage()
		return
	}

	for {
		dumper.dumpMetrics(generateMetrics(config))
		time.Sleep(*frequency)
	}
}
