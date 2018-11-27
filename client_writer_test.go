package herschel

import (
	"fmt"
	"image/color"
	"os"
	"testing"
	"time"

	"github.com/yokoe/herschel/option"
)

func TestWriting(t *testing.T) {
	spreadsheetID := createNewSpreadsheet(t)
	c := newTestClient(t)

	t.Run("Write values", func(t *testing.T) {
		table := NewTable(4, 3)
		table.PutValue(0, 0, "Updated: "+time.Now().Format("2006-01-02 03:04"))
		table.PutValue(1, 0, 1234567890)
		table.PutValue(2, 0, "12345000")
		table.PutValue(0, 1, "fuga")
		table.PutValue(1, 1, 0.2530)
		table.PutValue(2, 1, 49501)
		table.PutValue(3, 2, "Hello world")

		table.SetBackgroundColor(0, 0, color.Black)
		table.SetBackgroundColor(1, 1, color.RGBA{128, 0, 0, 0})
		table.SetNumberFormatPattern(1, 0, "#,###")
		table.SetNumberFormatPattern(2, 0, "#,##0")
		table.SetNumberFormatType(2, 0, "CURRENCY")
		table.SetNumberFormatPattern(1, 1, "#.00%")
		table.SetNumberFormatPattern(2, 1, "yyyy/MM")
		table.SetNumberFormatType(2, 1, "DATE")

		table.FrozenRowCount = 1
		table.FrozenColumnCount = 2

		c.RecreateSheet(spreadsheetID, t.Name())
		if err := c.WriteTable(spreadsheetID, t.Name(), table); err != nil {
			t.Fatal(err)
		}
	})

	t.Run("Set background colors beyond auto sized range", func(t *testing.T) {
		table := NewTable(3, 50)
		table.PutValue(0, 0, "Updated: "+time.Now().Format("2006-01-02 03:04"))
		table.SetBackgroundColor(2, 49, color.Black)

		if err := c.RecreateSheet(spreadsheetID, t.Name()); err != nil {
			t.Fatal(err)
		}

		if err := c.UpdateSheetGridLimits(spreadsheetID, t.Name(), table.GetRows(), table.GetCols()); err != nil {
			t.Fatal(err)
		}

		if err := c.WriteTable(spreadsheetID, t.Name(), table); err != nil {
			t.Fatal(err)
		}
	})

	t.Run("Clearing sheet values", func(t *testing.T) {
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

		t1, err := c.ReadTable(spreadsheetID, sheetTitle)
		if err != nil {
			t.Fatal(err)
		}
		if t1.GetRows() != 2 || t1.GetCols() != 1 {
			t.Fatalf("Unexpected dimensions of table. 2 x 1 expected, got %d x %d", t1.GetRows(), t1.GetCols())
		}

		if t1.GetValue(0, 0) == nil || t1.GetValue(1, 0) == nil {
			t.Fatal("Values at (0,0) and (1,0) should not be nil.")
		}

		if err := c.ClearSheetValues(spreadsheetID, sheetTitle); err != nil {
			t.Fatal(err)
		}

		t2, err := c.ReadTable(spreadsheetID, sheetTitle)
		if err != nil {
			t.Fatal(err)
		}
		if t2.GetRows() != 0 || t2.GetCols() != 0 {
			t.Fatalf("Unexpected dimensions of table. 0 x 0 after clearing sheet values expected, got %d x %d", t2.GetRows(), t2.GetCols())
		}
	})

}

/*
 * Helper functions
 */
func newTestClient(t *testing.T) *Client {
	credentialFilePath := os.Getenv("SPREADSHEET_CREDENTIAL_FILE")
	if len(credentialFilePath) == 0 {
		t.Skip("SPREADSHEET_CREDENTIAL_FILE is empty.")
	}

	client, err := NewClient(option.WithServiceAccountCredentials(credentialFilePath))
	if err != nil {
		t.Fatal(err)
	}
	return client
}

func createNewSpreadsheet(t *testing.T) string {
	ssID, err := newTestClient(t).CreateNewSpreadsheet(fmt.Sprintf("HerschelTest: %s", t.Name()))
	if err != nil {
		t.Fatalf("Failed to create new spreadsheet: %s", err)
	}
	return ssID
}
