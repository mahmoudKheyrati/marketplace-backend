package metric

import (
	"sync"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
)

type Metrics struct {
	MethodDurationsHistogram      *prometheus.HistogramVec
	MethodErrorDurationsHistogram *prometheus.HistogramVec //Duration of all error, include user and system
	MethodMaxSaturation           *prometheus.CounterVec
	MethodDurations               *prometheus.SummaryVec
	MethodErrorDurations          *prometheus.SummaryVec //Duration of all error, include user and system
	MethodCount                   *prometheus.CounterVec //Count of all rpc
	MethodSuccessCount            *prometheus.CounterVec //Count of RpcOk and RpcError
	MethodFailCount               *prometheus.CounterVec //Count of InternalError
	MethodUserErrorCount          *prometheus.CounterVec //Count of RpcError
	LogCount                      *prometheus.CounterVec
}

//nolint:gochecknoglobals
var (
	metricsOnce sync.Once
	metrics     *Metrics
)

func GetMetrics() *Metrics {
	metricsOnce.Do(func() {
		metrics = NewMetrics()
	})
	return metrics
}

func NewMetrics() *Metrics {
	methodLabels := []string{"service_name", "method", "origin"}
	errorLabels := []string{"service_name", "method", "origin", "error"}
	logLabels := []string{"service_name", "severity"}
	buckets := []float64{0.005, 0.01, 0.025, 0.05, 0.1, 0.2, 0.3, 0.4, 0.5, 0.6, 0.7, 0.8, 0.9,
		1, 1.25, 1.5, 1.75, 2, 2.25, 2.5, 2.75, 3, 3.5, 4, 4.5, 5, 6, 7, 8, 9, 10, 15, 20, 30, 40, 50, 100, 200, 300, 1000, 3600, 7200}
	metrics := &Metrics{
		MethodDurationsHistogram: prometheus.NewHistogramVec(
			prometheus.HistogramOpts{
				Name:    "method_durations_seconds",
				Help:    "Total Rpc latency.",
				Buckets: buckets,
			}, methodLabels),
		MethodErrorDurationsHistogram: prometheus.NewHistogramVec(
			prometheus.HistogramOpts{
				Name:    "method_error_durations_seconds",
				Help:    "Total Rpc latency.",
				Buckets: buckets,
			}, errorLabels),
		MethodDurations: prometheus.NewSummaryVec(
			prometheus.SummaryOpts{
				Name:       "method_durations_nanoseconds",
				Help:       "Total Rpc latency.",
				Objectives: map[float64]float64{0.5: 0.05, 0.9: 0.01, 0.99: 0.001},
			}, methodLabels),
		MethodErrorDurations: prometheus.NewSummaryVec(
			prometheus.SummaryOpts{
				Name:       "method_error_durations_nanoseconds",
				Help:       "Total Rpc latency.",
				Objectives: map[float64]float64{0.5: 0.05, 0.9: 0.01, 0.99: 0.001},
			}, methodLabels),
		MethodCount: promauto.NewCounterVec(prometheus.CounterOpts{
			Name: "method_total",
			Help: "The total number of rpc",
		}, methodLabels),
		MethodMaxSaturation: promauto.NewCounterVec(prometheus.CounterOpts{
			Name: "method_max",
			Help: "The maximum number of RPCs that service can handle, according to capacity planed.",
		}, methodLabels),
		MethodSuccessCount: promauto.NewCounterVec(prometheus.CounterOpts{
			Name: "method_success_total",
			Help: "The total number of successful rpc",
		}, methodLabels),
		MethodFailCount: promauto.NewCounterVec(prometheus.CounterOpts{
			Name: "method_failed_total",
			Help: "The total number of failed rpc",
		}, methodLabels),
		MethodUserErrorCount: promauto.NewCounterVec(prometheus.CounterOpts{
			Name: "method_error_total",
			Help: "The total number of user error",
		}, errorLabels),
		LogCount: promauto.NewCounterVec(prometheus.CounterOpts{
			Name: "log_total",
			Help: "The total number of logs",
		}, logLabels),
	}
	prometheus.MustRegister(metrics.MethodDurations, metrics.MethodErrorDurations, metrics.MethodDurationsHistogram, metrics.MethodErrorDurationsHistogram)
	return metrics
}
