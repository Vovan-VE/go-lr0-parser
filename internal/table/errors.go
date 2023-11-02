package table

import (
	"github.com/pkg/errors"
	"github.com/vovan-ve/go-lr0-parser/internal/symbol"
)

var (
	ErrState = errors.Wrap(symbol.ErrDefine, "bad state for table")

	// ErrConflictReduceReduce means that there are a number of rules which
	// applicable to reduce in the current state.
	ErrConflictReduceReduce = errors.Wrap(ErrState, "reduce-reduce conflict")
	// ErrConflictShiftReduce means that Shift and Reduce are both applicable
	// in the current state.
	ErrConflictShiftReduce = errors.Wrap(ErrState, "shift-reduce conflict")
)
