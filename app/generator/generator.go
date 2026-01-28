package generator

import (
	"github/com/lucasbn/sqlite-clone/app/machine/instructions"

	"github.com/xwb1989/sqlparser"
)

// resultConstructor allows the generator to construct T values for use in
// instructions like String. This follows the same pattern as BTreeEngine.
type resultConstructor[T any] interface {
	Text(string) T
	Number(int64) T
	Null() T
}

func Generate[T any](stmt sqlparser.Statement, rc resultConstructor[T]) []instructions.Instruction[T] {
	selectStmt := stmt.(*sqlparser.Select)
	tableExpr := selectStmt.From[0].(*sqlparser.AliasedTableExpr)
	tableName := sqlparser.GetTableName(tableExpr.Expr).String()

	return []instructions.Instruction[T]{
		// Schema lookup (addresses 0-10)
		// Open cursor 0 on the schema table (always page 1)
		instructions.OpenRead[T]{RootPage: 1, CursorID: 0},
		instructions.Rewind[T]{Cursor: 0},

		// Loop: check if this row is a "table" type entry
		instructions.Column[T]{Cursor: 0, Column: 0, Register: 1}, // Read 'type' into reg1
		instructions.String[T]{Register: 2, Value: rc.Text("table")},
		instructions.Ne[T]{Register1: 1, Register2: 2, JumpAddress: 9}, // If not "table", go to Next

		// Check if the name matches our target table
		instructions.Column[T]{Cursor: 0, Column: 1, Register: 3}, // Read 'name' into reg3
		instructions.String[T]{Register: 4, Value: rc.Text(tableName)},
		instructions.Ne[T]{Register1: 3, Register2: 4, JumpAddress: 9}, // If not target, go to Next

		// Found the table! Jump past the loop
		instructions.Goto[T]{JumpAddress: 11},

		// Next: advance cursor and loop back, or fall through to halt if exhausted
		instructions.Next[T]{Cursor: 0, FromAddress: 2},
		instructions.Halt[T]{}, // Table not found (error case)

		// Extract root page and open data cursor (addresses 11-13)
		instructions.Column[T]{Cursor: 0, Column: 3, Register: 5}, // Read rootpage into reg5
		instructions.OpenReadFromRegister[T]{CursorID: 1, PageRegister: 5},
		instructions.Rewind[T]{Cursor: 1},

		// Data reading (addresses 14+) - using cursor 1
		instructions.Column[T]{Cursor: 1, Column: 0, Register: 6},
		instructions.Column[T]{Cursor: 1, Column: 1, Register: 7},
		instructions.Column[T]{Cursor: 1, Column: 2, Register: 8},
		instructions.ResultRow[T]{FromRegister: 6, ToRegister: 8},
		instructions.Next[T]{Cursor: 1, FromAddress: 14},
		instructions.Halt[T]{},
	}
}
