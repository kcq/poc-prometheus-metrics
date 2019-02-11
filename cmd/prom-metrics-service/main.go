package main

import (
	"fmt"
	"github.com/go-chi/chi"
	"github.com/go-chi/render"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"math/rand"
	"net/http"
	"os"
	"time"
)

// General Prometheus metrics flow:
//
// * Add the Prom client package imports
// * Define your custom metrics objects
// * Optionally define a custom metrics registry (or use the default registry)
// * Register your metrics with the Prom registry
// * Optionally define a custom HTTP handler for your Prometheus metrics (or use the default HTTP handler)
// * Add the Prometheus HTTP handler to your HTTP router/mux

var (
	instrumentMetrics = false
	pid               = os.Getpid()
	registry          = prometheus.NewRegistry()
)

var basicHandler = promhttp.HandlerFor(
	registry,
	promhttp.HandlerOpts{},
)

var (
	callCounter = prometheus.NewCounter(prometheus.CounterOpts{
		Name: "ac_call_count",
		Help: "Number of successful API calls.",
	})
	callGauge = prometheus.NewGauge(prometheus.GaugeOpts{
		Name: "ac_call_active",
		Help: "Number of in-flight calls",
	})
	callHistogram = prometheus.NewHistogramVec(prometheus.HistogramOpts{
		Name:    "ac_call_time_ms",
		Help:    "API call duration distributions.",
		Buckets: []float64{5, 10, 15, 20, 25, 50, 100, 125, 150, 175, 200, 225, 250, 275, 300},
		//Buckets: prometheus.LinearBuckets(5,5,8),
		//Buckets: prometheus.ExponentialBuckets(5,2,8)
	},
		[]string{"code"},
	)
	callSummary = prometheus.NewSummaryVec(prometheus.SummaryOpts{
		Name:       "ac_call_time_distribution_ms",
		Help:       "API call duration distribution.",
		Objectives: map[float64]float64{0.5: 0.05, 0.9: 0.01, 0.95: 0.001, 0.99: 0.001},
	},
		[]string{"code"},
	)
)

// By default, the Prometheus client is configured to collect
// Go metrics (all platforms) and Process metrics (on Linux only).
// The custom registry in the PoC allows you to choose,
// which metrics collectors to use.

var collectors = []prometheus.Collector{
	//Process metrics
	prometheus.NewProcessCollector(
		//uses current process' PID from os.Getpid()
		prometheus.ProcessCollectorOpts{}),
	//Go runtime metrics
	prometheus.NewGoCollector(),
	//Custom metrics:
	callCounter,
	callGauge,
	callHistogram,
	callSummary,
}

// If you don't want to use a custom metrics registry you can register
// your custom metrics with the default registry in Prometheus:
//   prometheus.MustRegister(yourCustomMetrics)
// and then add the default http metrics handler:
//   r.Handle("/metrics",promhttp.Handler()) <- (chi http mux)
//   http.Handle("/metrics", promhttp.Handler()) <- (default http mux)

func init() {
	//Register the metrics collectors you want to use
	registry.MustRegister(collectors...)
}

func main() {
	fmt.Println("Prom PoC pid =", pid)

	r := chi.NewRouter()

	r.Get("/", func(w http.ResponseWriter, r *http.Request) {
		fmt.Println("GET /")
		stime := time.Now()

		callCounter.Inc()
		callGauge.Inc()
		defer callGauge.Dec()

		hcode := 200

		result := map[string]interface{}{
			"pid":  pid,
			"data": "something",
		}

		rv := rand.Intn(290)
		fmt.Println("doing work (random sleep to 290 ms) -", rv)
		time.Sleep(time.Millisecond * time.Duration(rv))

		dtime := time.Since(stime)
		ms := uint64(dtime / time.Millisecond)
		fmt.Printf("worked for %v milliseconds\n", ms)
		result["wtime"] = ms
		render.JSON(w, r, &result)

		callHistogram.WithLabelValues(fmt.Sprintf("%d", hcode)).Observe(float64(ms))
		callSummary.WithLabelValues(fmt.Sprintf("%d", hcode)).Observe(float64(ms))
	})

	var metricsHandler http.Handler
	if instrumentMetrics {
		metricsHandler = promhttp.InstrumentMetricHandler(registry, basicHandler)
	} else {
		metricsHandler = basicHandler
	}

	r.Handle("/metrics", metricsHandler)

	server := http.Server{Addr: ":7000", Handler: r}
	if err := server.ListenAndServe(); err != nil {
		fmt.Println("ListenAndServe error:", err)
	}
}
