package work

import (
	"github.com/flipkart-incubator/diligent/pkg/intgen"
	"github.com/flipkart-incubator/diligent/pkg/sqlgen"
	log "github.com/sirupsen/logrus"
	"time"
)

// SelectByUkRowWork is a CompositeWork that selects one row by unique key each time DoNext is called
// This class is Thread Safe. DoNext can be invoked from multiple threads
type SelectByUkRowWork struct {
	id        int
	runParams *RunParams
	recRange  *intgen.Range
	sqlGen    *sqlgen.SqlGen
}

func NewSelectByUkRowWork(id int, rp *RunParams, recRange *intgen.Range) CompositeWork {
	return &SelectByUkRowWork{
		id:        id,
		runParams: rp,
		recRange:  recRange,
		sqlGen:    sqlgen.NewSqlGen(rp.Table, rp.DataGen),
	}
}

// DoNext inserts a single record, directly without transactions
// This DoNext method never returns false (done) as records can be selected forever
func (w *SelectByUkRowWork) DoNext() bool {
	// Generate SQL statement
	sqlStmt := w.sqlGen.SelectByUkStatement(w.recRange.Rand())

	// Timed section of code starts here
	t := time.Now()

	log.Tracef("SQL: %s", sqlStmt)
	row := w.runParams.DB.QueryRow(sqlStmt)
	var pk, uniq, smallGrp, largeGrp, fixedVal, ts, payload string
	var seqNum int
	err := row.Scan(&pk, &uniq, &smallGrp, &largeGrp, &fixedVal, &seqNum, &ts, &payload)

	// Timed section ends here
	w.runParams.Metrics.ObserveStmtDuration("select", time.Since(t))

	if err != nil {
		log.Error(err)
		w.runParams.Metrics.ObserveStmtFailure("select")
		return true
	}

	return true
}

type SelectByUkTxnWork struct {
	id            int
	runParams     *RunParams
	recRange      *intgen.Range
	sqlGen        *sqlgen.SqlGen
}

func NewSelectByUkTxnWork(id int, rp *RunParams, recRange *intgen.Range) CompositeWork {
	return &SelectByUkTxnWork{
		id:        id,
		runParams: rp,
		recRange:  recRange,
		sqlGen:    sqlgen.NewSqlGen(rp.Table, rp.DataGen),
	}
}

// DoNext selects a batch of records with transaction
func (w *SelectByUkTxnWork) DoNext() bool {
	// Generate SQL statements
	sqlStatements := make([]string, w.runParams.BatchSize)
	for i := 0; i < len(sqlStatements); i++ {
		sqlStatements[i] = w.sqlGen.SelectByUkStatement(w.recRange.Rand())
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

	// Select each record
	for _, stmt := range sqlStatements {
		stmtStartTime = time.Now()
		log.Tracef("SQL: %s", stmt)
		row := tx.QueryRow(stmt)
		var pk, uniq, smallGrp, largeGrp, fixedVal, ts, payload string
		var seqNum int
		err = row.Scan(&pk, &uniq, &smallGrp, &largeGrp, &fixedVal, &seqNum, &ts, &payload)
		w.runParams.Metrics.ObserveStmtDuration("select", time.Since(stmtStartTime))
		if err != nil {
			log.Error(err)
			w.runParams.Metrics.ObserveStmtFailure("select")
			continue
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
