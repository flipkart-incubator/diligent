package work

import (
	"github.com/flipkart-incubator/diligent/pkg/intgen"
	"github.com/flipkart-incubator/diligent/pkg/sqlgen"
	log "github.com/sirupsen/logrus"
	"time"
)

// UpdateRowWork is a CompositeWork that updates one row each time DoNext is called
// This class is Thread Safe. DoNext can be invoked from multiple threads
type UpdateRowWork struct {
	id        int
	runParams *RunParams
	recRange  *intgen.Range
	sqlGen    *sqlgen.SqlGen
}

func NewUpdateRowWork(id int, rp *RunParams, recRange *intgen.Range) CompositeWork {
	return &UpdateRowWork{
		id:        id,
		runParams: rp,
		recRange:  recRange,
		sqlGen:    sqlgen.NewSqlGen(rp.Table, rp.DataGen),
	}
}

// DoNext updates a single record, directly without transactions
// This DoNext method never returns false (done) as records can be updated forever
func (w *UpdateRowWork) DoNext() bool {
	// Generate SQL statement
	sqlStmt := w.sqlGen.UpdatePayloadByPkStatement(w.recRange.Rand())

	// Timed section of code starts here
	t := time.Now()
	log.Tracef("SQL: %s", sqlStmt)
	result, err := w.runParams.DB.Exec(sqlStmt)
	// Timed section ends here
	w.runParams.Metrics.ObserveStmtDuration("update", time.Since(t))

	if err != nil {
		log.Error(err)
		w.runParams.Metrics.ObserveStmtFailure("update")
		return true
	}
	if count, _ := result.RowsAffected(); count != 1 {
		log.Errorf("Expected 1 row to be affected, but instead got %d", count)
		w.runParams.Metrics.ObserveStmtRowMismatch("update")
	}

	return true
}

type UpdateTxnWork struct {
	id        int
	runParams *RunParams
	recRange  *intgen.Range
	sqlGen    *sqlgen.SqlGen
}

func NewUpdateTxnWork(id int, rp *RunParams, recRange *intgen.Range) CompositeWork {
	return &UpdateTxnWork{
		id:        id,
		runParams: rp,
		recRange:  recRange,
		sqlGen:    sqlgen.NewSqlGen(rp.Table, rp.DataGen),
	}
}

// DoNext updates a batch of records with transaction
// This DoNext method never returns false (done) as records can be updated forever
func (w *UpdateTxnWork) DoNext() bool {
	// Generate SQL statements
	sqlStatements := make([]string, w.runParams.BatchSize)
	for i := 0; i < len(sqlStatements); i++ {
		sqlStatements[i] = w.sqlGen.UpdatePayloadByPkStatement(w.recRange.Rand())
	}

	// Timed section of code starts here
	txnStartTime := time.Now()

	stmtStartTime := time.Now()
	log.Tracef("SQL: Begin")
	tx, err := w.runParams.DB.Begin()
	w.runParams.Metrics.ObserveStmtDuration("begin", time.Since(stmtStartTime))
	if err != nil {
		log.Error(err)
		w.runParams.Metrics.ObserveStmtFailure("begin")
		return true
	}

	// Delete each record
	for _, stmt := range sqlStatements {
		stmtStartTime = time.Now()
		log.Tracef("SQL: %s", stmt)
		result, err := tx.Exec(stmt)
		w.runParams.Metrics.ObserveStmtDuration("update", time.Since(stmtStartTime))
		if err != nil {
			log.Error(err)
			w.runParams.Metrics.ObserveStmtFailure("update")
			continue
		}
		if count, _ := result.RowsAffected(); count != 1 {
			log.Errorf("Expected 1 row to be affected, but instead got %d", count)
			w.runParams.Metrics.ObserveStmtRowMismatch("update")
		}
	}

	// End the transaction
	stmtStartTime = time.Now()
	log.Tracef("SQL: Commit")
	err = tx.Commit()
	w.runParams.Metrics.ObserveStmtDuration("commit", time.Since(stmtStartTime))

	// Timed section for transaction ends here
	w.runParams.Metrics.ObserveTxnDuration(time.Since(txnStartTime))

	if err != nil {
		log.Error(err)
		w.runParams.Metrics.ObserveStmtFailure("commit")
		return true
	}

	return true
}
