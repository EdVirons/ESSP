package metrics

import (
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"net/http"
)

var (
	HTTPRequestsTotal = prometheus.NewCounterVec(
		prometheus.CounterOpts{
			Name: "http_requests_total",
			Help: "Total HTTP requests",
		},
		[]string{"method", "path", "status"},
	)

	HTTPRequestDuration = prometheus.NewHistogramVec(
		prometheus.HistogramOpts{
			Name:    "http_request_duration_seconds",
			Help:    "HTTP request duration",
			Buckets: prometheus.DefBuckets,
		},
		[]string{"method", "path"},
	)

	DBConnectionsActive = prometheus.NewGauge(prometheus.GaugeOpts{
		Name: "db_connections_active",
		Help: "Active database connections",
	})

	IncidentsCreated = prometheus.NewCounter(prometheus.CounterOpts{
		Name: "incidents_created_total",
		Help: "Total incidents created",
	})

	WorkOrdersCreated = prometheus.NewCounter(prometheus.CounterOpts{
		Name: "work_orders_created_total",
		Help: "Total work orders created",
	})
)

func init() {
	prometheus.MustRegister(HTTPRequestsTotal)
	prometheus.MustRegister(HTTPRequestDuration)
	prometheus.MustRegister(DBConnectionsActive)
	prometheus.MustRegister(IncidentsCreated)
	prometheus.MustRegister(WorkOrdersCreated)
}

func Handler() http.Handler {
	return promhttp.Handler()
}
