package herschel

import (
	"fmt"
)

// AppendTableAtBottom returns new instance of table appending another table.
func (t *Table) AppendTableAtBottom(a *Table) *Table {
	maxCols := t.cols
	if a.cols > maxCols {
		maxCols = a.cols
	}
	newTable := NewTable(t.rows+a.rows, maxCols)
	newTable.FrozenRowCount = t.FrozenRowCount
	newTable.FrozenColumnCount = t.FrozenColumnCount
	newTable.copyFromTable(t)

	for row := 0; row < a.rows; row++ {
		for col := 0; col < a.cols; col++ {
			newTable.copyCellFromTable(row+t.rows, col, a, row, col)
		}
	}
	return newTable
}

// AppendTableAtRight returns new instance of table appending another table to the right.
func (t *Table) AppendTableAtRight(a *Table) *Table {
	maxRows := t.rows
	if a.rows > maxRows {
		maxRows = a.rows
	}
	newTable := NewTable(maxRows, t.cols+a.cols)
	newTable.FrozenRowCount = t.FrozenRowCount
	newTable.FrozenColumnCount = t.FrozenColumnCount
	newTable.copyFromTable(t)

	for row := 0; row < a.rows; row++ {
		for col := 0; col < a.cols; col++ {
			newTable.copyCellFromTable(row, col+t.cols, a, row, col)
		}
	}
	return newTable
}

// SubTable returns new instance of table with sliced cells copied from original table
func (t *Table) SubTable(rowStart, colStart, numRows, numCols int) (*Table, error) {
	// Validate range
	if (rowStart + numRows) > t.rows {
		return nil, fmt.Errorf("%d is larger than original table row count (%d)", (rowStart + numRows), t.rows)
	}
	if (colStart + numCols) > t.cols {
		return nil, fmt.Errorf("%d is larger than original table col count (%d)", (colStart + numCols), t.cols)
	}

	s := NewTable(numRows, numCols)
	for row := 0; row < numRows; row++ {
		for col := 0; col < numCols; col++ {
			s.copyCellFromTable(row, col, t, row+rowStart, col+colStart)
		}
	}
	return s, nil
}

// SubTableByFilteringRows returns new instance of table with filtered rows from original table.
func (t *Table) SubTableByFilteringRows(f func(values []interface{}) bool) *Table {
	s := NewTable(0, t.cols)

	for i := 0; i < t.rows; i++ {
		r := t.GetValuesAtRow(i)
		if f(r) {
			st, err := t.SubTable(i, 0, 1, t.cols)
			if err != nil {
				panic(err)
			}
			s = s.AppendTableAtBottom(st)
		}
	}

	return s
}

// InsertColAtIndex inserts new column at index
func (t *Table) InsertColAtIndex(index int) error {
	if index < 0 || index > t.cols {
		return fmt.Errorf("invalid index %d", index)
	}

	t.cols++

	for col := t.cols - 1; col > index; col-- {
		for row := 0; row < t.rows; row++ {
			t.copyCellFromTable(row, col, t, row, col-1)
		}
	}

	for row := 0; row < t.rows; row++ {
		t.clearCell(row, index)
	}

	return nil
}

// RemoveColAtIndex removes column at index
func (t *Table) RemoveColAtIndex(index int) error {
	if index < 0 || index >= t.cols {
		return fmt.Errorf("invalid index %d", index)
	}

	for col := index; col < t.cols-1; col++ {
		for row := 0; row < t.rows; row++ {
			t.copyCellFromTable(row, col, t, row, col+1)
		}
	}

	t.cols = t.cols - 1

	return nil
}

func (t *Table) copyCellFromTable(targetRow int, targetCol int, sourceTable *Table, sourceRow int, sourceCol int) {
	t.PutValue(targetRow, targetCol, sourceTable.GetValue(sourceRow, sourceCol))
	t.SetBackgroundColor(targetRow, targetCol, sourceTable.getBackgroundColor(sourceRow, sourceCol))
	t.SetNumberFormatPattern(targetRow, targetCol, sourceTable.getNumberFormatPattern(sourceRow, sourceCol))
	t.SetNumberFormatType(targetRow, targetCol, sourceTable.getNumberFormatType(sourceRow, sourceCol))
}

func (t *Table) copyFromTable(a *Table) {
	for row := 0; row < a.rows; row++ {
		for col := 0; col < a.cols; col++ {
			t.copyCellFromTable(row, col, a, row, col)
		}
	}
}

//ClearValues clears all values of table.
func (t *Table) ClearValues() error {

	for row := 0; row < t.rows; row++ {
		for col := 0; col < t.cols; col++ {
			t.PutValue(row, col, "")
		}
	}

	return nil
}

//ClearValuesInRange clears the value of a cell in the specified range.
func (t *Table) ClearValuesInRange(rowStart, colStart, numRows, numCols int) error {
	// Validate range
	if (rowStart + numRows) > t.rows {
		return fmt.Errorf("%d is larger than original table row count (%d)", (rowStart + numRows), t.rows)
	}
	if (colStart + numCols) > t.cols {
		return fmt.Errorf("%d is larger than original table col count (%d)", (colStart + numCols), t.cols)
	}

	for row := rowStart; row < (rowStart + numRows); row++ {
		for col := colStart; col < (colStart + numCols); col++ {
			t.PutValue(row, col, "")
		}
	}
	return nil
}
