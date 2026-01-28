package instructions

import (
	"github/com/lucasbn/sqlite-clone/app/machine/common"
	"github/com/lucasbn/sqlite-clone/app/machine/state"
)

type Goto[T any] struct {
	JumpAddress uint64
}

var _ Instruction[any] = Goto[any]{}

func (g Goto[T]) Execute(s *state.MachineState[T], b common.BTreeEngine[T]) [][]T {
	s.CurrentAddress = int(g.JumpAddress)
	return [][]T{}
}
