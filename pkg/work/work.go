package work

import "github.com/flipkart-incubator/diligent/pkg/intgen"

// CompositeWork represents a unit work that has multiple sub parts
type CompositeWork interface {
	// DoNext executes the next sub part of the work. It returns false when all sub parts are done
	DoNext() (hasMore bool)
}

type CompositeWorkFactory func(id int, rp *RunParams, recRange *intgen.Range) CompositeWork
