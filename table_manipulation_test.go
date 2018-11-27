package herschel

import "testing"

func TestSubTable(t *testing.T) {
	orig := NewTable(3, 3)
	orig.PutValuesAtRow(0, "a", "b", "c")
	orig.PutValuesAtRow(1, "d", "e", "f")
	orig.PutValuesAtRow(2, "g", "h", "i")

	t.Run("SuccessfulSlicing", func(t *testing.T) {
		s, err := orig.SubTable(1, 1, 1, 2)
		if err != nil {
			t.Fatal(err)
		}

		t.Run("SizeOfNewTable", func(t *testing.T) {
			if s.rows != 1 || s.cols != 2 {
				t.Errorf("Sub table should have dimension %dx%d, got: %dx%d", 1, 2, s.rows, s.cols)
			}
		})

		t.Run("ValuesOfNewTable", func(t *testing.T) {
			if s.GetValue(0, 0) != "e" {
				t.Errorf("Value at %d, %d should be %s, got: %s", 0, 0, "e", s.GetValue(0, 0))
			}
			if s.GetValue(0, 1) != "f" {
				t.Errorf("Value at %d, %d should be %s, got: %s", 0, 1, "f", s.GetValue(0, 1))
			}
		})
	})

}

func TestSubTableByFilteringRows(t *testing.T) {
	orig := NewTable(5, 3)
	orig.PutValuesAtRow(0, "a", "b", "c")
	orig.PutValuesAtRow(1, "d", "e", "f")
	orig.PutValuesAtRow(2, "g", "h", "i")
	orig.PutValuesAtRow(3, "j", "a", "b")
	orig.PutValuesAtRow(4, "x", "x", "a")

	st := orig.SubTableByFilteringRows(func(row []interface{}) bool {
		for _, r := range row {
			if r == "a" {
				return true
			}
		}
		return false
	})

	if st.GetRows() != 3 {
		t.Errorf("SubTable should have 3 rows. Got %d rows.", st.GetRows())
	}
}

func TestRemovingColumn(t *testing.T) {
	orig := NewTable(5, 3)
	orig.PutValuesAtRow(0, "a", "b", "c")
	orig.PutValuesAtRow(1, "d", "e", "f")
	orig.PutValuesAtRow(2, "g", "h", "i")
	orig.PutValuesAtRow(3, "j", "a", "b")
	orig.PutValuesAtRow(4, "x", "x", "a")

	if orig.cols != 3 {
		t.Fatalf("Unexpected number of columns. 3 expected, got: %d", orig.cols)
	}

	if err := orig.RemoveColAtIndex(1); err != nil {
		t.Fatal(err)
	}

	expectedValues := []struct {
		row   int
		col   int
		value string
	}{
		{0, 0, "a"},
		{0, 1, "c"},
		{1, 0, "d"},
		{1, 1, "f"},
	}

	for _, expectedValue := range expectedValues {
		actualValue := orig.GetValue(expectedValue.row, expectedValue.col)
		if actualValue != expectedValue.value {
			t.Errorf("Unexpected value at (%d,%d). %s expected, got: %s", expectedValue.row, expectedValue.col, expectedValue.value, actualValue)
		}
	}
}

func TestClearValues(t *testing.T) {
	orig := NewTable(3, 3)
	orig.PutValuesAtRow(0, "a", "b", "c")
	orig.PutValuesAtRow(1, "d", "e", "f")
	orig.PutValuesAtRow(2, "g", "h", "i")

	if err := orig.ClearValues(); err != nil {
		t.Fatal(err)
	}

	for row := 0; row < orig.rows; row++ {
		for col := 0; col < orig.cols; col++ {
			if value := orig.GetValue(row, col); value != "" {
				t.Errorf("The value of cell at (%d,%d) is not empty.", row, col)
			}
		}
	}

}

func TestClearValuesInRange(t *testing.T) {
	orig := NewTable(5, 4)
	orig.PutValuesAtRow(0, "a", "b", "c", "d")
	orig.PutValuesAtRow(1, "e", "f", "g", "h")
	orig.PutValuesAtRow(2, "i", "j", "k", "l")
	orig.PutValuesAtRow(3, "m", "n", "o", "p")
	orig.PutValuesAtRow(4, "q", "r", "s", "t")

	if err := orig.ClearValuesInRange(1, 1, 3, 2); err != nil {
		t.Fatal(err)
	}

	exp := NewTable(5, 4)
	exp.PutValuesAtRow(0, "a", "b", "c", "d")
	exp.PutValuesAtRow(1, "e", "", "", "h")
	exp.PutValuesAtRow(2, "i", "", "", "l")
	exp.PutValuesAtRow(3, "m", "", "", "p")
	exp.PutValuesAtRow(4, "q", "r", "s", "t")

	for row := 0; row < exp.rows; row++ {
		for col := 0; col < exp.cols; col++ {
			expectedValue := exp.GetValue(row, col)
			actualValue := orig.GetValue(row, col)
			if expectedValue != actualValue {
				t.Errorf("The cell value of (% d,% d) is Unexpected.", row, col)
			}
		}
	}
}

func TestInsertCol(t *testing.T) {
	t.Run("First", func(t *testing.T) {
		orig := NewTable(3, 2)
		orig.PutValuesAtRow(0, "a", "b")
		orig.PutValuesAtRow(1, "c", "d")
		orig.PutValuesAtRow(2, "e", "f")

		if err := orig.InsertColAtIndex(0); err != nil {
			t.Fatal(err)
		}

		if orig.GetRows() != 3 || orig.GetCols() != 3 {
			t.Fatalf("Table should have dimension of 3 x 3, got: %d x %d", orig.GetRows(), orig.GetCols())
		}

		exp := NewTable(3, 3)
		exp.PutValuesAtRow(0, nil, "a", "b")
		exp.PutValuesAtRow(1, nil, "c", "d")
		exp.PutValuesAtRow(2, nil, "e", "f")

		for row := 0; row < exp.rows; row++ {
			for col := 0; col < exp.cols; col++ {
				expectedValue := exp.GetValue(row, col)
				actualValue := orig.GetValue(row, col)
				if expectedValue != actualValue {
					t.Errorf("The cell value of (% d,% d) is Unexpected.", row, col)
				}
			}
		}
	})
	t.Run("Middle", func(t *testing.T) {
		orig := NewTable(3, 2)
		orig.PutValuesAtRow(0, "a", "b")
		orig.PutValuesAtRow(1, "c", "d")
		orig.PutValuesAtRow(2, "e", "f")

		if err := orig.InsertColAtIndex(1); err != nil {
			t.Fatal(err)
		}

		if orig.GetRows() != 3 || orig.GetCols() != 3 {
			t.Fatalf("Table should have dimension of 3 x 3, got: %d x %d", orig.GetRows(), orig.GetCols())
		}

		exp := NewTable(3, 3)
		exp.PutValuesAtRow(0, "a", nil, "b")
		exp.PutValuesAtRow(1, "c", nil, "d")
		exp.PutValuesAtRow(2, "e", nil, "f")

		for row := 0; row < exp.rows; row++ {
			for col := 0; col < exp.cols; col++ {
				expectedValue := exp.GetValue(row, col)
				actualValue := orig.GetValue(row, col)
				if expectedValue != actualValue {
					t.Errorf("The cell value of (% d,% d) is Unexpected.", row, col)
				}
			}
		}
	})
	t.Run("Last", func(t *testing.T) {
		orig := NewTable(3, 2)
		orig.PutValuesAtRow(0, "a", "b")
		orig.PutValuesAtRow(1, "c", "d")
		orig.PutValuesAtRow(2, "e", "f")

		if err := orig.InsertColAtIndex(2); err != nil {
			t.Fatal(err)
		}

		if orig.GetRows() != 3 || orig.GetCols() != 3 {
			t.Fatalf("Table should have dimension of 3 x 3, got: %d x %d", orig.GetRows(), orig.GetCols())
		}

		exp := NewTable(3, 3)
		exp.PutValuesAtRow(0, "a", "b", nil)
		exp.PutValuesAtRow(1, "c", "d", nil)
		exp.PutValuesAtRow(2, "e", "f", nil)

		for row := 0; row < exp.rows; row++ {
			for col := 0; col < exp.cols; col++ {
				expectedValue := exp.GetValue(row, col)
				actualValue := orig.GetValue(row, col)
				if expectedValue != actualValue {
					t.Errorf("The cell value of (% d,% d) is Unexpected.", row, col)
				}
			}
		}
	})
}
