package instructions

import (
	"github/com/lucasbn/sqlite-clone/app/machine/common"
	"github/com/lucasbn/sqlite-clone/app/machine/state"
)

type String[T any] struct {
	Register int
	Value    T
}

var _ Instruction[int] = String[int]{}

func (str String[T]) Execute(s *state.MachineState[T], b common.BTreeEngine[T]) [][]T {
	s.CurrentAddress++

	s.Registers = s.Registers.Set(str.Register, str.Value)
	return [][]T{}
}
