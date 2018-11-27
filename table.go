package herschel

import (
	"fmt"
	"image/color"
	"log"
	"reflect"
)

// Table represents 2 dimension cells.
type Table struct {
	cols              int
	rows              int
	values            map[int]map[int]interface{}
	backgroundColors  map[int]map[int]color.Color
	numberFormats     map[int]map[int]string
	numberFormatTypes map[int]map[int]string
	FrozenRowCount    int64
	FrozenColumnCount int64
}

func (t Table) String() string {
	return fmt.Sprintf("{Table %d rows x %d cols}", t.rows, t.cols)
}

// NewTable returns instance of Table.
func NewTable(rows int, cols int) *Table {
	instance := &Table{cols: cols, rows: rows}

	values := map[int]map[int]interface{}{}
	backgroundColors := map[int]map[int]color.Color{}
	numberFormats := map[int]map[int]string{}
	numberFormatTypes := map[int]map[int]string{}

	for i := 0; i < rows; i++ {
		values[i] = map[int]interface{}{}
		backgroundColors[i] = map[int]color.Color{}
		numberFormats[i] = map[int]string{}
		numberFormatTypes[i] = map[int]string{}
	}

	instance.values = values
	instance.backgroundColors = backgroundColors
	instance.numberFormats = numberFormats
	instance.numberFormatTypes = numberFormatTypes
	instance.FrozenRowCount = 0
	instance.FrozenColumnCount = 0
	return instance
}

// GetRows returns number of rows
func (t *Table) GetRows() int {
	return t.rows
}

// GetCols returns numfer of cols
func (t *Table) GetCols() int {
	return t.cols
}

// PutValue updates value of cell.
func (t *Table) PutValue(row int, col int, value interface{}) {
	if row >= t.rows {
		panic(fmt.Sprintf("row out of order: %d in %d", row, t.rows))
	}
	if col >= t.cols {
		panic(fmt.Sprintf("col out of order: %d in %d", col, t.cols))
	}
	t.values[row][col] = value
}

// GetValue returns value of cell.
func (t *Table) GetValue(row int, col int) interface{} {
	return t.values[row][col]
}

// GetValuesAtRow returns a slice containing value of cells at row
func (t *Table) GetValuesAtRow(row int) []interface{} {
	cells := []interface{}{}
	for i := 0; i < t.cols; i++ {
		cells = append(cells, t.GetValue(row, i))
	}
	return cells
}

// PutValuesAtRow sets values of cells at row
func (t *Table) PutValuesAtRow(row int, values ...interface{}) {
	for i := 0; i < len(values); i++ {
		t.PutValue(row, i, values[i])
	}
}

// IndexOfRowWithPrefix returns index of row which contains values as prefix. Returns -1 when there is no row matches.
func (t *Table) IndexOfRowWithPrefix(values ...interface{}) int {
	if len(values) == 0 {
		return -1
	}

	for i, v := range values {
		if reflect.TypeOf(v).Kind() == reflect.Slice {
			log.Printf("Unsupported type (slice) in args[%d] (%+v)\n", i, v)
		}
	}

L:
	for i := 0; i < t.rows; i++ {
		valuesAtRow := t.GetValuesAtRow(i)
		if len(valuesAtRow) < len(values) {
			continue L
		}

		for j := 0; j < len(values); j++ {
			if values[j] != valuesAtRow[j] {
				continue L
			}
		}

		return i
	}
	return -1
}

// Values returns slice of cell values
func (t *Table) Values() [][]interface{} {
	values := [][]interface{}{}
	for row := 0; row < t.rows; row++ {
		rowValues := []interface{}{}
		for col := 0; col < t.cols; col++ {
			rowValues = append(rowValues, t.GetValue(row, col))
		}
		values = append(values, rowValues)
	}
	return values
}

// SetBackgroundColor sets background color of cell at (row, col)
func (t *Table) SetBackgroundColor(row int, col int, c color.Color) {
	t.backgroundColors[row][col] = c
}

func (t *Table) getBackgroundColor(row int, col int) color.Color {
	if c, exists := t.backgroundColors[row][col]; exists {
		return c
	}
	return color.Transparent
}

// SetNumberFormatPattern sets number format pettern of cell at (row, col)
func (t *Table) SetNumberFormatPattern(row int, col int, pattern string) {
	t.numberFormats[row][col] = pattern
}

func (t *Table) getNumberFormatPattern(row int, col int) string {
	if f, exists := t.numberFormats[row][col]; exists {
		return f
	}
	return ""
}

// SetNumberFormatType sets number format type of cell at (row, col)
func (t *Table) SetNumberFormatType(row int, col int, formatType string) {
	t.numberFormatTypes[row][col] = formatType
}

func (t *Table) getNumberFormatType(row int, col int) string {
	if f, exists := t.numberFormatTypes[row][col]; exists {
		return f
	}
	return ""
}

// PutCommaSeparatedInt64 set value of cell at (row, col) as comma separated integer.
func (t *Table) PutCommaSeparatedInt64(row int, col int, value int64) {
	t.PutValue(row, col, value)
	t.SetNumberFormatPattern(row, col, "#,##0")
}

func (t *Table) clearCell(row int, col int) {
	delete(t.values[row], col)
	delete(t.backgroundColors[row], col)
	delete(t.numberFormats[row], col)
	delete(t.numberFormatTypes[row], col)
}
