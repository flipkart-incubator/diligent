package work

import (
	"github.com/flipkart-incubator/diligent/pkg/intgen"
	"github.com/flipkart-incubator/diligent/pkg/sqlgen"
	log "github.com/sirupsen/logrus"
	"go.uber.org/atomic"
	"math"
	"time"
)

// DeleteRowWork is a CompositeWork that deletes one row each time DoNext is called
//
// Since the same record cannot be deleted twice, the DeleteRowWork takes a slice of row indexes
// which has the indexes of the rows (as per the DataGen) that this instance should delete.
// The creator of DeleteRowWork must ensure that it is giving mutually exclusive sets of row indexes
// to the different DeleteRowWork instances.
//
// This class is Thread Safe. DoNext can be invoked from multiple threads
type DeleteRowWork struct {
	id         int
	runParams  *RunParams
	recIndexes []int
	sqlGen     *sqlgen.SqlGen
	currentPos *atomic.Int32
}

func NewDeleteRowWork(id int, rp *RunParams, recRange *intgen.Range) CompositeWork {
	recIndexes := recRange.Ints()
	intgen.Shuffle(recIndexes)
	return &DeleteRowWork{
		id:         id,
		runParams:  rp,
		recIndexes: recIndexes,
		sqlGen:     sqlgen.NewSqlGen(rp.Table, rp.DataGen),
		currentPos: atomic.NewInt32(0),
	}
}

// DoNext deletes a single record, directly without transactions
// Returns true as long as there are more rows to be deleteed
// Returns false (indicating done) once delete has been attempted for all the rows that this DeleteRowWork
// has been assigned
func (w *DeleteRowWork) DoNext() bool {
	pos := int(w.currentPos.Load())
	w.currentPos.Inc()

	// Are we done?
	if pos >= len(w.recIndexes) {
		return false
	}

	sqlStmt := w.sqlGen.DeleteByPkStatement(w.recIndexes[pos])

	// Timed section of code starts here
	t := time.Now()
	log.Tracef("SQL: %s", sqlStmt)
	result, err := w.runParams.DB.Exec(sqlStmt)
	// Timed section ends here
	w.runParams.Metrics.ObserveStmtDuration("delete", time.Since(t))

	if err != nil {
		log.Error(err)
		w.runParams.Metrics.ObserveStmtFailure("delete")
		return true
	}
	if count, _ := result.RowsAffected(); count != 1 {
		log.Errorf("Expected 1 row to be affected, but instead got %d", count)
		w.runParams.Metrics.ObserveStmtRowMismatch("delete")
	}

	return true
}

// DeleteTxnWork is a CompositeWork that deletes one botch of rows using a transaction each time DoNext is called
//
// Since the same record cannot be deleteed twice, the DeleteRowWork takes a slice of row indexes
// which has the indexes of the rows (as per the DataGen) that this instance should delete.
// The creator of DeleteRowWork must ensure that it is giving mutually exclusive sets of row indexes
// to the different DeleteRowWork instances.
//
// This class is Thread Safe. DoNext can be invoked from multiple threads
type DeleteTxnWork struct {
	id         int
	runParams  *RunParams
	recIndexes []int
	sqlGen     *sqlgen.SqlGen
	currentPos *atomic.Int32
}

func NewDeleteTxnWork(id int, rp *RunParams, recRange *intgen.Range) CompositeWork {
	recIndexes := recRange.Ints()
	intgen.Shuffle(recIndexes)
	return &DeleteTxnWork{
		id:         id,
		runParams:  rp,
		recIndexes: recIndexes,
		sqlGen:     sqlgen.NewSqlGen(rp.Table, rp.DataGen),
		currentPos: atomic.NewInt32(0),
	}
}

// DoNext deletes a batch of records with transaction
// Returns true as long as there are more rows to be deleted
// Returns false (indicating done) once delete has been attempted for all the rows that this DeleteRowWork
// has been assigned
func (w *DeleteTxnWork) DoNext() bool {
	batchStart := int(w.currentPos.Load())
	w.currentPos.Add(int32(w.runParams.BatchSize))

	// Are we done?
	if batchStart >= len(w.recIndexes) {
		return false
	}

	batchLimit := int(math.Min(float64(batchStart+w.runParams.BatchSize), float64(len(w.recIndexes))))
	batch := w.recIndexes[batchStart:batchLimit]
	sqlStatements := make([]string, len(batch))
	for i, rec := range batch {
		sqlStatements[i] = w.sqlGen.DeleteByPkStatement(rec)
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
		w.runParams.Metrics.ObserveStmtDuration("delete", time.Since(stmtStartTime))
		if err != nil {
			log.Error(err)
			w.runParams.Metrics.ObserveStmtFailure("delete")
			continue
		}
		if count, _ := result.RowsAffected(); count != 1 {
			log.Errorf("Expected 1 row to be affected, but instead got %d", count)
			w.runParams.Metrics.ObserveStmtRowMismatch("delete")
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
