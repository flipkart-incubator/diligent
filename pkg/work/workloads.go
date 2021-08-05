package work

import (
	"fmt"
	"github.com/flipkart-incubator/diligent/pkg/intgen"
)

func NewInsertRowWorkload(assignedRange *intgen.Range, rp *RunParams) *Workload {
	return makeWorkload("insert-row-rand", assignedRange, Partitioned, rp, NewInsertRowWork)
}

func NewInsertTxnWorkload(assignedRange *intgen.Range, rp *RunParams) *Workload {
	return makeWorkload("insert-txn-rand", assignedRange, Partitioned, rp, NewInsertTxnWork)
}

func NewSelectByPkRowWorkload(assignedRange *intgen.Range, rp *RunParams) *Workload {
	return makeWorkload("select-pk-row-rand", assignedRange, Shared, rp, NewSelectByPkRowWork)
}

func NewSelectByPkTxnWorkload(assignedRange *intgen.Range, rp *RunParams) *Workload {
	return makeWorkload("select-pk-txn-rand", assignedRange, Shared, rp, NewSelectByPkTxnWork)
}

func NewSelectByUkRowWorkload(assignedRange *intgen.Range, rp *RunParams) *Workload {
	return makeWorkload("select-uk-row-rand", assignedRange, Shared, rp, NewSelectByPkRowWork)
}

func NewSelectByUkTxnWorkload(assignedRange *intgen.Range, rp *RunParams) *Workload {
	return makeWorkload("select-uk-txn-rand", assignedRange, Shared, rp, NewSelectByPkTxnWork)
}

func NewUpdateRowWorkload(assignedRange *intgen.Range, rp *RunParams) *Workload {
	return makeWorkload("update-row-rand", assignedRange, Shared, rp, NewUpdateRowWork)
}

func NewUpdateTxnWorkload(assignedRange *intgen.Range, rp *RunParams) *Workload {
	return makeWorkload("update-txn-rand", assignedRange, Shared, rp, NewUpdateTxnWork)
}

func NewDeleteRowWorkload(assignedRange *intgen.Range, rp *RunParams) *Workload {
	return makeWorkload("delete-row-rand", assignedRange, Partitioned, rp, NewDeleteRowWork)
}

func NewDeleteTxnWorkload(assignedRange *intgen.Range, rp *RunParams) *Workload {
	return makeWorkload("delete-txn-rand", assignedRange, Partitioned, rp, NewDeleteTxnWork)
}


type WorkAllocationStrategy int
const (
	Partitioned WorkAllocationStrategy = iota
	Shared
)

func makeWorkload(
	name string,
	assignedRange *intgen.Range,
	workAllocStrategy WorkAllocationStrategy,
	rp *RunParams,
	workFactory CompositeWorkFactory) *Workload {

	// Create new workload with given name
	workload := NewWorkload(name, rp)

	// Partition the range of all records (number of partitions = concurrency)
	var workerRanges []*intgen.Range
	switch workAllocStrategy {
	case Partitioned:
		workerRanges = assignedRange.Partition(rp.Concurrency)
	case Shared:
		workerRanges = assignedRange.Duplicate(rp.Concurrency)
	default:
		panic(fmt.Errorf("invalid work workAllocStrategy: %d", workAllocStrategy))
	}

	// For each partition, create a work definition using the workFactory
	// Then create worker in workload that will do this work
	for i, workRange := range workerRanges {
		work := workFactory(i, rp, workRange)
		workload.CreateWorkerWithWork(work)
	}

	return workload
}