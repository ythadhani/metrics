package main

import (
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"net/http"
	"time"
)

var statusCodeMap = map[int]int{http.StatusOK: 1, http.StatusServiceUnavailable: 0}

func writeMetricsForURL(urlUpGauge, responseTimeGauge *prometheus.GaugeVec, url string) {
	go func() {
		for {
			start := time.Now()
			resp, err := http.Get(url)
			if err != nil {
				panic(err)
			}
			elapsed := time.Since(start).Milliseconds()

			var (
				present bool
				statusCode int
			)
			if statusCode, present = statusCodeMap[resp.StatusCode]; !present {
				statusCode = -1
			}
			urlUpGauge.WithLabelValues(url).Set(float64(statusCode))
			responseTimeGauge.WithLabelValues(url).Set(float64(elapsed))
			time.Sleep(2*time.Second)
		}
	}()
}

func main() {

	urlUpGauge := prometheus.NewGaugeVec(prometheus.GaugeOpts{
		Name: "sample_external_url_up",
		Help: "Is URL up",
	}, []string{"url"})

	responseTimeGauge := prometheus.NewGaugeVec(prometheus.GaugeOpts{
		Name: "sample_external_url_response_ms",
		Help: "URL Response time in milliseconds",
	}, []string{"url"})

	prometheus.MustRegister(urlUpGauge)
	prometheus.MustRegister(responseTimeGauge)

	writeMetricsForURL(urlUpGauge, responseTimeGauge,"https://httpstat.us/503")
	writeMetricsForURL(urlUpGauge, responseTimeGauge,"https://httpstat.us/200")

	http.Handle("/metrics", promhttp.Handler())
	http.ListenAndServe(":9090", nil)
}