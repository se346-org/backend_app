package observability

import (
	"time"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
)

type Metrics struct {
	// HTTP metrics
	HTTPRequestsTotal    *prometheus.CounterVec
	HTTPRequestDuration  *prometheus.HistogramVec
	HTTPRequestsInFlight prometheus.Gauge

	// Database metrics
	DBConnectionsActive prometheus.Gauge
	DBQueriesTotal      *prometheus.CounterVec
	DBQueryDuration     *prometheus.HistogramVec

	// WebSocket metrics
	WSConnectionsActive prometheus.Gauge
	WSMessagesTotal     *prometheus.CounterVec

	// Business metrics
	UsersRegistered     prometheus.Counter
	MessagesTotal       *prometheus.CounterVec
	ConversationsTotal  prometheus.Counter
}

func NewMetrics() *Metrics {
	return &Metrics{
		HTTPRequestsTotal: promauto.NewCounterVec(
			prometheus.CounterOpts{
				Name: "http_requests_total",
				Help: "Total number of HTTP requests",
			},
			[]string{"method", "endpoint", "status"},
		),
		HTTPRequestDuration: promauto.NewHistogramVec(
			prometheus.HistogramOpts{
				Name:    "http_request_duration_seconds",
				Help:    "HTTP request duration in seconds",
				Buckets: prometheus.DefBuckets,
			},
			[]string{"method", "endpoint"},
		),
		HTTPRequestsInFlight: promauto.NewGauge(
			prometheus.GaugeOpts{
				Name: "http_requests_in_flight",
				Help: "Number of HTTP requests currently being processed",
			},
		),
		DBConnectionsActive: promauto.NewGauge(
			prometheus.GaugeOpts{
				Name: "db_connections_active",
				Help: "Number of active database connections",
			},
		),
		DBQueriesTotal: promauto.NewCounterVec(
			prometheus.CounterOpts{
				Name: "db_queries_total",
				Help: "Total number of database queries",
			},
			[]string{"operation", "table", "status"},
		),
		DBQueryDuration: promauto.NewHistogramVec(
			prometheus.HistogramOpts{
				Name:    "db_query_duration_seconds",
				Help:    "Database query duration in seconds",
				Buckets: prometheus.DefBuckets,
			},
			[]string{"operation", "table"},
		),
		WSConnectionsActive: promauto.NewGauge(
			prometheus.GaugeOpts{
				Name: "ws_connections_active",
				Help: "Number of active WebSocket connections",
			},
		),
		WSMessagesTotal: promauto.NewCounterVec(
			prometheus.CounterOpts{
				Name: "ws_messages_total",
				Help: "Total number of WebSocket messages",
			},
			[]string{"type", "status"},
		),
		UsersRegistered: promauto.NewCounter(
			prometheus.CounterOpts{
				Name: "users_registered_total",
				Help: "Total number of registered users",
			},
		),
		MessagesTotal: promauto.NewCounterVec(
			prometheus.CounterOpts{
				Name: "messages_total",
				Help: "Total number of messages sent",
			},
			[]string{"conversation_type"},
		),
		ConversationsTotal: promauto.NewCounter(
			prometheus.CounterOpts{
				Name: "conversations_total",
				Help: "Total number of conversations created",
			},
		),
	}
}

func (m *Metrics) RecordHTTPRequest(method, endpoint, status string, duration time.Duration) {
	m.HTTPRequestsTotal.WithLabelValues(method, endpoint, status).Inc()
	m.HTTPRequestDuration.WithLabelValues(method, endpoint).Observe(duration.Seconds())
}

func (m *Metrics) RecordDBQuery(operation, table, status string, duration time.Duration) {
	m.DBQueriesTotal.WithLabelValues(operation, table, status).Inc()
	m.DBQueryDuration.WithLabelValues(operation, table).Observe(duration.Seconds())
} 