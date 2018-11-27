package herschel

import (
	"testing"
)

func TestReading(t *testing.T) {
	spreadsheetID := createNewSpreadsheet(t)
	c := newTestClient(t)

	// Create test sheet
	sheetTitle := t.Name()
	if err := c.RecreateSheet(spreadsheetID, sheetTitle); err != nil {
		t.Fatal(err)
	}

	table := NewTable(2, 1)
	table.PutValue(0, 0, "Hello")
	table.PutValue(1, 0, "World")

	if err := c.WriteTable(spreadsheetID, sheetTitle, table); err != nil {
		t.Fatal(err)
	}

	t.Run("Read", func(t *testing.T) {
		_, err := c.Read(spreadsheetID, sheetTitle)
		if err != nil {
			t.Fatal(err)
		}
	})

	t.Run("ReadTable", func(t *testing.T) {
		_, err := c.ReadTable(spreadsheetID, sheetTitle)
		if err != nil {
			t.Fatal(err)
		}
	})

}
