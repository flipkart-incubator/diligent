package metrics

import (
	"database/sql"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	log "github.com/sirupsen/logrus"
	"net/http"
	"time"
)

type DiligentMetrics struct {
	metricsAddr               string
	configTxnsEnabledGauge    prometheus.Gauge
	configOpsPerTxnGauge      prometheus.Gauge
	configConcurrencyGauge    prometheus.Gauge
	concurrencyGauge          prometheus.Gauge
	stmtDurationHistVec       *prometheus.HistogramVec
	txnDurationHist           prometheus.Histogram
	stmtFailCounterVec        *prometheus.CounterVec
	stmtRowMismatchCounterVec *prometheus.CounterVec
	dbConnectionsGaugeVec     *prometheus.GaugeVec
}

func NewDiligentMetrics(metricsAddr string) *DiligentMetrics {
	dm := &DiligentMetrics{
		metricsAddr: metricsAddr,
	}

	dm.configTxnsEnabledGauge = prometheus.NewGauge(
		prometheus.GaugeOpts{
			Namespace: "diligent_sql",
			Subsystem: "basic",
			Name:      "config_transactions_enabled",
			Help:      "Set to 0 or 1 for transactions disabled, enabled respectively",
		})

	dm.configOpsPerTxnGauge = prometheus.NewGauge(
		prometheus.GaugeOpts{
			Namespace: "diligent_sql",
			Subsystem: "basic",
			Name:      "config_operations_per_transaction",
			Help:      "Configured number of operations per transaction",
		})

	dm.configConcurrencyGauge = prometheus.NewGauge(
		prometheus.GaugeOpts{
			Namespace: "diligent_sql",
			Subsystem: "basic",
			Name:      "config_concurrency",
			Help:      "Configured concurrency",
		})

	dm.concurrencyGauge = prometheus.NewGauge(
		prometheus.GaugeOpts{
			Namespace: "diligent_sql",
			Subsystem: "basic",
			Name:      "concurrency",
			Help:      "Actual number of concurrent workers",
		})

	dm.stmtDurationHistVec = prometheus.NewHistogramVec(
		prometheus.HistogramOpts{
			Namespace: "diligent_sql",
			Subsystem: "basic",
			Name:      "statement_duration_seconds",
			Help:      "Histogram of the duration taken to execute SQL statements",
			Buckets:   prometheus.ExponentialBuckets(0.000001, 2, 30),
		}, []string{"statement"})

	dm.txnDurationHist = prometheus.NewHistogram(
		prometheus.HistogramOpts{
			Namespace: "diligent_sql",
			Subsystem: "basic",
			Name:      "transaction_duration_seconds",
			Help:      "Histogram of the duration of SQL transactions",
			Buckets:   prometheus.ExponentialBuckets(0.000001, 2, 30),
		})

	dm.stmtFailCounterVec = prometheus.NewCounterVec(
		prometheus.CounterOpts{
			Namespace: "diligent_sql",
			Subsystem: "basic",
			Name:      "statement_failure",
			Help:      "Counter for failed SQL statements",
		}, []string{"statement"})

	dm.stmtRowMismatchCounterVec = prometheus.NewCounterVec(
		prometheus.CounterOpts{
			Namespace: "diligent_sql",
			Subsystem: "basic",
			Name:      "statement_affected_rows_mismatch",
			Help:      "Counter for SQL statements with mismatch in affected number of rows",
		}, []string{"statement"})

	dm.dbConnectionsGaugeVec = prometheus.NewGaugeVec(
		prometheus.GaugeOpts{
			Namespace: "diligent_sql",
			Subsystem: "basic",
			Name:      "db_connections",
			Help:      "Number of db connections - max open, open, in use and idle",
		}, []string{"label"})

	return dm
}

func (dm *DiligentMetrics) Register() {
	prometheus.MustRegister(dm.configTxnsEnabledGauge)
	prometheus.MustRegister(dm.configOpsPerTxnGauge)
	prometheus.MustRegister(dm.configConcurrencyGauge)
	prometheus.MustRegister(dm.concurrencyGauge)
	prometheus.MustRegister(dm.stmtDurationHistVec)
	prometheus.MustRegister(dm.txnDurationHist)
	prometheus.MustRegister(dm.stmtFailCounterVec)
	prometheus.MustRegister(dm.stmtRowMismatchCounterVec)
	prometheus.MustRegister(dm.dbConnectionsGaugeVec)

	http.Handle("/metrics", promhttp.Handler())
	go func() {
		err := http.ListenAndServe(dm.metricsAddr, nil)
		if err != nil {
			log.Fatal(err)
		}
	}()
}

func (dm *DiligentMetrics) ObserveStmtDuration(statement string, d time.Duration) {
	dm.stmtDurationHistVec.WithLabelValues(statement).Observe(d.Seconds())
}

func (dm *DiligentMetrics) ObserveTxnDuration(d time.Duration) {
	dm.txnDurationHist.Observe(d.Seconds())
}

func (dm *DiligentMetrics) SetConfigMetricsForWorkload(txnsEnabled bool, opsPerTxn int, concurrency int) {
	if txnsEnabled {
		dm.configTxnsEnabledGauge.Set(1.0)
	} else {
		dm.configTxnsEnabledGauge.Set(0.0)
	}

	dm.configOpsPerTxnGauge.Set(float64(opsPerTxn))

	dm.configConcurrencyGauge.Set(float64(concurrency))
}

func (dm *DiligentMetrics) UnsetConfigMetricsForWorkload() {
	dm.configTxnsEnabledGauge.Set(0)
	dm.configOpsPerTxnGauge.Set(0)
	dm.configConcurrencyGauge.Set(0)
}

func (dm *DiligentMetrics) IncConcurrencyForWorkload() {
	dm.concurrencyGauge.Inc()
}

func (dm *DiligentMetrics) DecConcurrencyForWorkload() {
	dm.concurrencyGauge.Dec()
}

func (dm *DiligentMetrics) ObserveStmtFailure(statement string) {
	dm.stmtFailCounterVec.WithLabelValues(statement).Inc()
}

func (dm *DiligentMetrics) ObserveStmtRowMismatch(statement string) {
	dm.stmtRowMismatchCounterVec.WithLabelValues(statement).Inc()
}

func (dm *DiligentMetrics) ObserveDbConn(db *sql.DB) {
	stats := db.Stats()
	dm.dbConnectionsGaugeVec.WithLabelValues("maxopen").Set(float64(stats.MaxOpenConnections))
	dm.dbConnectionsGaugeVec.WithLabelValues("open").Set(float64(stats.OpenConnections))
	dm.dbConnectionsGaugeVec.WithLabelValues("inuse").Set(float64(stats.InUse))
	dm.dbConnectionsGaugeVec.WithLabelValues("idle").Set(float64(stats.Idle))
}
