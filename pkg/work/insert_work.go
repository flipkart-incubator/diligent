package work

import (
	"fmt"
	"github.com/flipkart-incubator/diligent/pkg/intgen"
	"github.com/flipkart-incubator/diligent/pkg/sqlgen"
	log "github.com/sirupsen/logrus"
	"go.uber.org/atomic"
	"math"
	"time"
)

// InsertRowWork is a CompositeWork that inserts one row each time DoNext is called
//
// Since the same record cannot be inserted twice, the InsertRowWork takes a slice of row indexes
// which has the indexes of the rows (as per the DataGen) that this instance should insert.
// The creator of InsertRowWork must ensure that it is giving mutually exclusive sets of row indexes
// to the different InsertRowWork instances.
//
// This class is Thread Safe. DoNext can be invoked from multiple threads
type InsertRowWork struct {
	id         int
	runParams  *RunParams
	recIndexes []int
	sqlGen     *sqlgen.SqlGen
	currentPos *atomic.Int32
}

func NewInsertRowWork(id int, rp *RunParams, recRange *intgen.Range) CompositeWork {
	recIndexes := recRange.Ints()
	intgen.Shuffle(recIndexes)
	return &InsertRowWork{
		id:         id,
		runParams:  rp,
		recIndexes: recIndexes,
		sqlGen:     sqlgen.NewSqlGen(rp.Table, rp.DataGen),
		currentPos: atomic.NewInt32(0),
	}
}

// DoNext inserts a single record, directly without transactions
// Returns true as long as there are more rows to be inserted
// Returns false (indicating done) once insert has been attempted for all the rows that this InsertRowWork
// has been assigned
func (w *InsertRowWork) DoNext() (bool, error) {
	pos := int(w.currentPos.Load())
	w.currentPos.Inc()

	// Are we done?
	if pos >= len(w.recIndexes) {
		return false, nil
	}

	sqlStmt := w.sqlGen.InsertStatement(w.recIndexes[pos])

	// Timed section of code starts here
	t := time.Now()
	log.Tracef("SQL: %s", sqlStmt)
	result, err := w.runParams.DB.Exec(sqlStmt)
	// Timed section ends here
	w.runParams.Metrics.ObserveStmtDuration("insert", time.Since(t))

	if err != nil {
		log.Error(err)
		w.runParams.Metrics.ObserveStmtFailure("insert")
		return true, err
	}
	if count, _ := result.RowsAffected(); count != 1 {
		e := fmt.Errorf("expected 1 row to be affected, but instead got %d", count)
		log.Error(e)
		w.runParams.Metrics.ObserveStmtRowMismatch("insert")
		return true, e
	}

	return true, nil
}

// InsertTxnWork is a CompositeWork that inserts one botch of rows using a transaction each time DoNext is called
//
// Since the same record cannot be inserted twice, the InsertRowWork takes a slice of row indexes
// which has the indexes of the rows (as per the DataGen) that this instance should insert.
// The creator of InsertRowWork must ensure that it is giving mutually exclusive sets of row indexes
// to the different InsertRowWork instances.
//
// This class is Thread Safe. DoNext can be invoked from multiple threads
type InsertTxnWork struct {
	id         int
	runParams  *RunParams
	recIndexes []int
	sqlGen     *sqlgen.SqlGen
	currentPos *atomic.Int32
}

func NewInsertTxnWork(id int, rp *RunParams, recRange *intgen.Range) CompositeWork {
	recIndexes := recRange.Ints()
	intgen.Shuffle(recIndexes)
	return &InsertTxnWork{
		id:         id,
		runParams:  rp,
		recIndexes: recIndexes,
		sqlGen:     sqlgen.NewSqlGen(rp.Table, rp.DataGen),
		currentPos: atomic.NewInt32(0),
	}
}

// DoNext inserts a batch of records with transaction
// Returns true as long as there are more rows to be inserted
// Returns false (indicating done) once insert has been attempted for all the rows that this InsertRowWork
// has been assigned
func (w *InsertTxnWork) DoNext() (bool, error) {
	batchStart := int(w.currentPos.Load())
	w.currentPos.Add(int32(w.runParams.BatchSize))

	// Are we done?
	if batchStart >= len(w.recIndexes) {
		return false, nil
	}

	batchLimit := int(math.Min(float64(batchStart+w.runParams.BatchSize), float64(len(w.recIndexes))))
	batch := w.recIndexes[batchStart:batchLimit]
	sqlStatements := make([]string, len(batch))
	for i, rec := range batch {
		sqlStatements[i] = w.sqlGen.InsertStatement(rec)
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
		return true, err
	}

	// Insert each record
	statementErrors := 0
	for _, stmt := range sqlStatements {
		stmtStartTime = time.Now()
		log.Tracef("SQL: %s", stmt)
		result, err := tx.Exec(stmt)
		w.runParams.Metrics.ObserveStmtDuration("insert", time.Since(stmtStartTime))
		if err != nil {
			log.Error(err)
			w.runParams.Metrics.ObserveStmtFailure("insert")
			statementErrors++
			continue
		}
		if count, _ := result.RowsAffected(); count != 1 {
			log.Errorf("Expected 1 row to be affected, but instead got %d", count)
			w.runParams.Metrics.ObserveStmtRowMismatch("insert")
			statementErrors++
		}
	}

	// End the transaction
	stmtStartTime = time.Now()
	log.Trace("SQL: Commit")
	err = tx.Commit()
	w.runParams.Metrics.ObserveStmtDuration("commit", time.Since(stmtStartTime))

	// Timed section for transaction ends here
	w.runParams.Metrics.ObserveTxnDuration(time.Since(txnStartTime))

	if err != nil {
		log.Error(err)
		w.runParams.Metrics.ObserveStmtFailure("commit")
		return true, err
	}

	if statementErrors > 0 {
		e := fmt.Errorf("encountered errors with statements within the transaction")
		log.Error(e)
		return true, e
	}

	return true, nil
}
