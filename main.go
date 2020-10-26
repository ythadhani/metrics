package main

import (
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"net/http"
	"time"
)

var statusCodeMap = map[int]int{http.StatusOK: 1, http.StatusServiceUnavailable: 0}

func writeMetricsForURL(urlUpGauge, responseTimeGauge *prometheus.GaugeVec, url string) {
	/*
	 The main goroutine would block while exposing the scrape endpoint (listening on 9090).
	 Hence, a separate goroutine for gathering/generating metrics. Could've gathered metrics by probing
	 external URLs when Prometheus hits our /metrics endpoint but blocking a Prometheus request
	 isn't good design for production if it takes our app longer to gather metrics.
	*/
	go func() {
		for {
			/*
			 A simple approach to measure the time taken for a query by recording
			 current time and measuring time elapsed after triggering a http Get.
			 Can also calculate the time taken only for content transfer using:
			 https://github.com/davecheney/httpstat
			*/
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

			// Writing metric values to the respective gauges
			urlUpGauge.WithLabelValues(url).Set(float64(statusCode))
			responseTimeGauge.WithLabelValues(url).Set(float64(elapsed))

			time.Sleep(2*time.Second)
		}
	}()
}

func main() {

	// urlUpGauge records the metrics for the http status codes
	urlUpGauge := prometheus.NewGaugeVec(prometheus.GaugeOpts{
		Name: "sample_external_url_up",
		Help: "Is URL up",
	}, []string{"url"})

	// responseTimeGauge records the metrics for the response time in milliseconds
	responseTimeGauge := prometheus.NewGaugeVec(prometheus.GaugeOpts{
		Name: "sample_external_url_response_ms",
		Help: "URL Response time in milliseconds",
	}, []string{"url"})

	// Registering our gauges so that their data is accessible when Prometheus scrapes
	// our /metrics endpoint
	prometheus.MustRegister(urlUpGauge)
	prometheus.MustRegister(responseTimeGauge)

	writeMetricsForURL(urlUpGauge, responseTimeGauge,"https://httpstat.us/503")
	writeMetricsForURL(urlUpGauge, responseTimeGauge,"https://httpstat.us/200")

	http.Handle("/metrics", promhttp.Handler())
	// Awaiting Prometheus scrapes on port 9090
	http.ListenAndServe(":9090", nil)
}