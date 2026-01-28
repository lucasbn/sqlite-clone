package instructions

import (
	"github/com/lucasbn/sqlite-clone/app/machine/common"
	"github/com/lucasbn/sqlite-clone/app/machine/state"
	"github/com/lucasbn/sqlite-clone/app/types"
)

type OpenReadFromRegister[T any] struct {
	CursorID     uint64
	PageRegister int
}

var _ Instruction[any] = OpenReadFromRegister[any]{}

func (o OpenReadFromRegister[T]) Execute(s *state.MachineState[T], b common.BTreeEngine[T]) [][]T {
	s.CurrentAddress++

	// Get the root page number from the register
	regValue := s.Registers.Get(o.PageRegister)

	// Type assert to extract the int64 value from NumberEntry
	entry, ok := any(regValue).(types.Entry)
	if !ok {
		panic("OpenReadFromRegister: register value is not an Entry")
	}

	numEntry, ok := entry.(types.NumberEntry)
	if !ok {
		panic("OpenReadFromRegister: register value is not a NumberEntry")
	}

	rootPage := uint64(numEntry.Value)
	b.NewCursor(o.CursorID, rootPage)

	return [][]T{}
}
