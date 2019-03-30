package herschel

import (
	"fmt"
	"strings"
	"testing"
)

func TestTable(t *testing.T) {
	table := NewTable(10, 10)
	table.PutValue(0, 0, "Hello")
	table.PutValue(1, 2, "World")

	// Get method
	if table.GetValue(0, 0) != "Hello" {
		t.Errorf("Hello expected, got: %v", table.GetValue(0, 0))
	}

	if table.GetValue(1, 2) != "World" {
		t.Errorf("World expected, got: %v", table.GetValue(1, 2))
	}

	if table.GetValue(100, 100) != nil {
		t.Errorf("Nil expected, got: %v", table.GetValue(100, 100))
	}

	// Slice output
	testCases := []struct {
		row   int
		col   int
		value string
	}{
		{0, 0, "Hello"},
		{1, 2, "World"},
	}
	values := table.Values()
	for _, tc := range testCases {
		t.Run(fmt.Sprintf("(%d, %d) %s", tc.row, tc.col, tc.value), func(t *testing.T) {
			if values[tc.row][tc.col] != tc.value {
				t.Errorf("%s expected, got: %v", tc.value, values[tc.row][tc.col])
			}
		})
	}
}

func TestRowsAndCols(t *testing.T) {
	table := NewTable(3, 5)
	if table.GetRows() != 3 {
		t.Errorf("Table should have 3 rows, got: %d", table.GetRows())
	}

	if table.GetCols() != 5 {
		t.Errorf("Table should have 5 cols, got: %d", table.GetCols())
	}
}

func TestAppendingBottom(t *testing.T) {
	firstTable := NewTable(2, 2)
	secondTable := NewTable(3, 3)
	firstTable.PutValue(0, 0, "First")
	secondTable.PutValue(1, 1, "Second")
	merged := firstTable.AppendTableAtBottom(secondTable)

	if merged.cols != 3 {
		t.Errorf("Merged table should have 3 cols.")
	}

	if merged.rows != 5 {
		t.Errorf("Merged table should have 5 rows.")
	}

	testCases := []struct {
		row   int
		col   int
		value string
	}{
		{0, 0, "First"},
		{3, 1, "Second"},
	}
	values := merged.Values()
	for _, tc := range testCases {
		t.Run(fmt.Sprintf("(%d, %d) %s", tc.row, tc.col, tc.value), func(t *testing.T) {
			if values[tc.row][tc.col] != tc.value {
				t.Errorf("%s expected, got: %v", tc.value, values[tc.row][tc.col])
			}
		})
	}
}

func TestAppendingRight(t *testing.T) {
	firstTable := NewTable(2, 2)
	secondTable := NewTable(3, 3)
	firstTable.PutValue(0, 0, "First")
	secondTable.PutValue(1, 1, "Second")
	merged := firstTable.AppendTableAtRight(secondTable)

	if merged.cols != 5 {
		t.Errorf("Merged table should have 5 cols.")
	}

	if merged.rows != 3 {
		t.Errorf("Merged table should have 3 rows.")
	}

	testCases := []struct {
		row   int
		col   int
		value string
	}{
		{0, 0, "First"},
		{1, 3, "Second"},
	}
	values := merged.Values()
	for _, tc := range testCases {
		t.Run(fmt.Sprintf("(%d, %d) %s", tc.row, tc.col, tc.value), func(t *testing.T) {
			if values[tc.row][tc.col] != tc.value {
				t.Errorf("%s expected, got: %v", tc.value, values[tc.row][tc.col])
			}
		})
	}
}

func TestGetValuesAtRow(t *testing.T) {
	table := NewTable(3, 3)
	for i := 0; i < 3; i++ {
		for j := 0; j < 3; j++ {
			table.PutValue(i, j, fmt.Sprintf("%d,%d", i, j))
		}
	}

	testCases := []struct {
		row    int
		values []string
	}{
		{0, []string{"0,0", "0,1", "0,2"}},
		{1, []string{"1,0", "1,1", "1,2"}},
		{2, []string{"2,0", "2,1", "2,2"}},
	}

	for _, tc := range testCases {
		t.Run(fmt.Sprintf("GetValuesAtRow %d", tc.row), func(t *testing.T) {
			values := table.GetValuesAtRow(tc.row)
			if len(values) != len(tc.values) {
				t.Fatalf("Length of values should be %d, got: %d", len(tc.values), len(values))
			}

			for i := 0; i < len(values); i++ {
				if values[i] != tc.values[i] {
					t.Errorf("Value of %d,%d should be %v, got: %v", tc.row, i, tc.values[i], values[i])
				}
			}
		})
	}
}

func TestPutValuesAtRow(t *testing.T) {
	table := NewTable(3, 3)

	table.PutValuesAtRow(0, "a", "b", "c")
	table.PutValuesAtRow(1, "d", "e", "f")
	table.PutValuesAtRow(2, "g", "h", "i")

	testCases := []struct {
		row    int
		values []string
	}{
		{0, []string{"a", "b", "c"}},
		{1, []string{"d", "e", "f"}},
		{2, []string{"g", "h", "i"}},
	}

	for _, tc := range testCases {
		t.Run(fmt.Sprintf("GetValuesAtRow %d", tc.row), func(t *testing.T) {
			values := table.GetValuesAtRow(tc.row)
			if len(values) != len(tc.values) {
				t.Fatalf("Length of values should be %d, got: %d", len(tc.values), len(values))
			}

			for i := 0; i < len(values); i++ {
				if values[i] != tc.values[i] {
					t.Errorf("Value of %d,%d should be %v, got: %v", tc.row, i, tc.values[i], values[i])
				}
			}
		})
	}
}

func TestIndexOfRow(t *testing.T) {
	table := NewTable(3, 3)
	table.PutValuesAtRow(0, "a", "b", "c")
	table.PutValuesAtRow(1, "d", "e", "f")
	table.PutValuesAtRow(2, "g", "h", "i")

	testCases := []struct {
		values []interface{}
		row    int
	}{
		{[]interface{}{"a", "b", "c"}, 0},
		{[]interface{}{"d", "e", "f"}, 1},
		{[]interface{}{"g", "h", "i"}, 2},
		{[]interface{}{"d", "e"}, 1},
		{[]interface{}{"d"}, 1},
		{[]interface{}{"e", "f"}, -1},
		{[]interface{}{[]string{"a"}}, -1},
	}

	for _, tc := range testCases {
		t.Run(fmt.Sprintf("IndexOfRow %v", tc.values), func(t *testing.T) {
			index := table.IndexOfRowWithPrefix(tc.values...)
			if index != tc.row {
				t.Errorf("Index of row starts with %+v should be %d, got: %d", tc.values, tc.row, index)
			}
		})
	}

}

func describeTableCells(table *Table) string {
	rows := []string{}
	for row := 0; row < table.rows; row++ {
		values := []string{}
		for col := 0; col < table.cols; col++ {
			val := table.GetValue(col, row)
			if val != nil {
				values = append(values, val.(string))
			}
			values = append(values, " ")
		}
		rows = append(rows, "["+strings.Join(values, ", ")+"]")
	}
	return strings.Join(rows, "\n")
}

func TestTableToMapConversion(t *testing.T) {
	tb := NewTable(3, 2)
	tb.PutValuesAtRow(0, "foo", "value")
	tb.PutValuesAtRow(1, "bar", 123)
	tb.PutValuesAtRow(2, "foo", "overwritten")

	m := tb.ToMap()
	if m["foo"] != "overwritten" {
		t.Errorf("m[foo] = %v, want 'overwritten'", m["foo"])
	}
	if m["bar"] != 123 {
		t.Errorf("m[bar] = %v, want 123", m["bar"])
	}
}
