package herschel

import (
	"bytes"
	"encoding/csv"
	"reflect"
	"strings"
	"testing"
)

func TestCSVExport(t *testing.T) {
	table := NewTable(2, 2)
	table.PutValuesAtRow(0, "a", "b")
	table.PutValuesAtRow(1, 1, int64(2))

	buf := bytes.NewBufferString("")
	if err := table.ToCSV(buf); err != nil {
		t.Fatal(err)
	}

	r := csv.NewReader(strings.NewReader(buf.String()))

	for i, expect := range [][]string{{"a", "b"}, {"1", "2"}} {
		row, err := r.Read()
		if err != nil {
			t.Fatal(err)
		}
		if !reflect.DeepEqual(row, expect) {
			t.Errorf("Line %d: Expect a, b got %v\nCSV: %s", i, row, buf.String())
		}
	}
}
