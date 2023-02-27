package main

import (
	"log"

	"gopkg.in/alecthomas/kingpin.v2"
)

var (
	count     = kingpin.Flag("count", "Number of data points to generate").Short('c').Default("100").Int()
	frequency = kingpin.Flag("frequency", "Frequency of data points").Short('f').Default("10s").Duration()

	dump = kingpin.Command("dump", "Dump data to stdout")

	remoteWrite = kingpin.Command("remote-write", "Remote write to a Prometheus-compatible TSDB")
	addr        = remoteWrite.Flag("addr", "").URL()

	emitMetrics = kingpin.Command("emit", "Emit metrics over HTTP on a given port (localhost:<port>/metrics)")
	port        = emitMetrics.Flag("port", "HTTP port for emitting metrics").Default("9090").Int()
)

const (
	Version = "0.0.1"
)

func main() {
	kingpin.Version(Version)
	kingpin.Parse()

	switch kingpin.Parse() {
	case dump.FullCommand():
		log.Println("dump")
	case remoteWrite.FullCommand():
		log.Println("remoteWrite")
	case emitMetrics.FullCommand():
		log.Println("emitMetrics")
	}
}
