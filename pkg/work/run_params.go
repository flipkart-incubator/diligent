package work

import (
	"database/sql"
	"github.com/flipkart-incubator/diligent/pkg/datagen"
	"github.com/flipkart-incubator/diligent/pkg/metrics"
)

type RunParams struct {
	DB          *sql.DB
	DataGen     *datagen.DataGen
	Metrics     *metrics.DiligentMetrics
	Table       string
	Concurrency int
	EnableTxn   bool
	BatchSize   int
	DurationSec int
}
