package herschel

import (
	"fmt"
	"testing"
)

func TestTableWriteRead(t *testing.T) {
	spreadsheetID := createNewSpreadsheet(t)
	sheetTitle := t.Name()

	table := NewTable(2, 2)
	table.PutValue(0, 0, "0,0")
	table.PutValue(1, 0, "1,0")
	table.PutValue(0, 1, "0,1")
	table.PutValue(1, 1, "1,1")

	client := newTestClient(t)
	client.RecreateSheet(spreadsheetID, sheetTitle)
	if err := client.WriteTable(spreadsheetID, sheetTitle, table); err != nil {
		t.Fatal(err)
	}

	readTable, err := client.ReadTable(spreadsheetID, sheetTitle)
	if err != nil {
		t.Fatal(err)
	}

	if readTable.rows < 2 || readTable.cols < 2 {
		t.Fatalf("Table read from spreadsheet should have at last 2 x 2 dimension. Got: %d x %d", readTable.rows, readTable.cols)
	}

	for c := 0; c < 2; c++ {
		for r := 0; r < 2; r++ {
			if readTable.GetValue(r, c) != fmt.Sprintf("%d,%d", r, c) {
				t.Errorf("Unexpected value at %d, %d. %+v", r, c, readTable.GetValue(r, c))
			}
		}
	}

}
