package main

import (
	"flag"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

var (
	apikey         string
	locationIds    string
	listenAddress  string
	cacheTTLString string
)

func init() {
	flag.StringVar(&apikey, "apikey", "", "API key for OpenWeatherMap (required)")
	flag.StringVar(&locationIds, "location-ids", "", "Comma separated list of location IDs to fetch weather for (required)")
	flag.StringVar(&listenAddress, "listen-address", "localhost:9755", "Hostname and port to listen on")
	flag.StringVar(&cacheTTLString, "cache-ttl", "10m", "TTL for API response cache")
}

func main() {
	flag.Parse()

	if apikey == "" || locationIds == "" {
		flag.PrintDefaults()
		os.Exit(1)
	}

	cacheTTL, err := time.ParseDuration(cacheTTLString)
	if err != nil {
		log.Println("Cache TTL could not be parsed:", err)
		os.Exit(1)
	}

	exporter := newExporter(apikey, locationIds, cacheTTL)
	if len(exporter.getData().Locations) == 0 {
		os.Exit(1)
	}

	prometheus.MustRegister(exporter)

	http.Handle("/metrics", promhttp.Handler())
	log.Fatal(http.ListenAndServe(listenAddress, nil))
}
