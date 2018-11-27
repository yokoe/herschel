package herschel

import (
	"testing"
)

func TestCreateNewSpreadsheet(t *testing.T) {
	id, err := newTestClient(t).CreateNewSpreadsheet("HerschelTest")
	if err != nil {
		t.Fatal(err)
	}
	if len(id) == 0 {
		t.Fatal("Spreadsheet ID is empty.")
	}
}
