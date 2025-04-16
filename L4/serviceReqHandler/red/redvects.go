package red

import "github.com/prometheus/client_golang/prometheus"

const (
	labelHandler = "handler"
	labelMethod  = "method"
	labelStatus  = "status"
	labelResult  = "result"
	labelService = "service"
	labelQuery   = "query"
)

var (
	duration = prometheus.NewSummaryVec(
		prometheus.SummaryOpts{
			Name:       "request_duration_seconds",
			Help:       "Request duration in seconds",
			Objectives: map[float64]float64{0.5: 0.05, 0.9: 0.01, 0.99: 0.001},
		},
		[]string{labelHandler, labelMethod, labelStatus},
	)

	errorsTotal = prometheus.NewCounterVec(
		prometheus.CounterOpts{
			Name: "errors_total",
			Help: "Total number of errors",
		},
		[]string{labelHandler, labelMethod, labelStatus},
	)

	requestsTotal = prometheus.NewCounterVec(
		prometheus.CounterOpts{
			Name: "requests_total",
			Help: "Total number of requests",
		},
		[]string{labelHandler, labelMethod},
	)
)

var (
	Duration = prometheus.NewSummaryVec(
		prometheus.SummaryOpts{
			Name:       "out duration seconds",
			Help:       "Request duration in seconds",
			Objectives: map[float64]float64{0.9: 0.01, 0.95: 0.005, 0.99: 0.001},
		},
		[]string{labelService, labelQuery, labelResult},
	)

	ErrorsTotal = prometheus.NewCounterVec(
		prometheus.CounterOpts{
			Name: "out errors_total",
			Help: "Total number of errors",
		},
		[]string{labelService, labelQuery, labelResult},
	)

	RequestsTotal = prometheus.NewCounterVec(
		prometheus.CounterOpts{
			Name: "out requests_total",
			Help: "Total number of requests",
		},
		[]string{labelService, labelQuery},
	)
)

func init() {
	prometheus.MustRegister(duration)
	prometheus.MustRegister(errorsTotal)
	prometheus.MustRegister(requestsTotal)

	prometheus.MustRegister(Duration)
	prometheus.MustRegister(ErrorsTotal)
	prometheus.MustRegister(RequestsTotal)
}
