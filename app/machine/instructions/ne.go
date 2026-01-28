package instructions

import (
	"github/com/lucasbn/sqlite-clone/app/machine/common"
	"github/com/lucasbn/sqlite-clone/app/machine/state"
	"github/com/lucasbn/sqlite-clone/app/types"
)

type Ne[T any] struct {
	Register1   int
	Register2   int
	JumpAddress uint64
}

var _ Instruction[any] = Ne[any]{}

func (ne Ne[T]) Execute(s *state.MachineState[T], b common.BTreeEngine[T]) [][]T {
	val1 := s.Registers.Get(ne.Register1)
	val2 := s.Registers.Get(ne.Register2)

	equal := compareValues(val1, val2)

	if !equal {
		s.CurrentAddress = int(ne.JumpAddress)
	} else {
		s.CurrentAddress++
	}

	return [][]T{}
}

// compareValues compares two values for equality. It handles the specific
// Entry types used by this SQLite clone.
func compareValues[T any](a, b T) bool {
	// Type assert to Entry interface and compare based on type
	aEntry, aOk := any(a).(types.Entry)
	bEntry, bOk := any(b).(types.Entry)

	if !aOk || !bOk {
		return false
	}

	switch aVal := aEntry.(type) {
	case types.TextEntry:
		if bVal, ok := bEntry.(types.TextEntry); ok {
			return aVal.Value == bVal.Value
		}
	case types.NumberEntry:
		if bVal, ok := bEntry.(types.NumberEntry); ok {
			return aVal.Value == bVal.Value
		}
	case types.NullEntry:
		_, ok := bEntry.(types.NullEntry)
		return ok
	}

	return false
}
