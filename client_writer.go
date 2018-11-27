package herschel

import (
	"fmt"

	"google.golang.org/api/sheets/v4"
)

// Write writes values to spreadsheet
func (client Client) Write(spreadsheetID string, sheetTitle string, values [][]interface{}) error {
	return client.updateCellValues(spreadsheetID, sheetTitle, values)
}

// WriteTable writes values of table to spreadsheet
func (client Client) WriteTable(spreadsheetID string, sheetTitle string, table *Table) error {
	if err := client.Write(spreadsheetID, sheetTitle, table.Values()); err != nil {
		return err
	}
	return client.setCellFormats(spreadsheetID, sheetTitle, table)
}

// AddSheet adds new sheet with title
func (client Client) AddSheet(spreadsheetID string, sheetTitle string) error {
	return addSheet(client, spreadsheetID, sheetTitle)
}

// DeleteSheet deletes a sheet with title.
func (client Client) DeleteSheet(spreadsheetID string, sheetTitle string) error {
	sheetID, exists, err := getSheetID(client, spreadsheetID, sheetTitle)
	if err != nil {
		return err
	}
	if !exists {
		return nil
	}
	return deleteSheetByID(client, spreadsheetID, sheetID)
}

// RecreateSheet deletes a sheet with title and adds new one.
func (client Client) RecreateSheet(spreadsheetID string, sheetTitle string) error {
	return recreateSheet(client, spreadsheetID, sheetTitle)
}

// ClearSheetValues clears values of sheet.
func (client Client) ClearSheetValues(spreadsheetID string, sheetTitle string) error {
	_, err := client.service.Spreadsheets.Values.Clear(spreadsheetID, sheetTitle, &sheets.ClearValuesRequest{}).Do()
	return err
}

// UpdateSheetGridLimits updates grid limits of sheet.
func (client Client) UpdateSheetGridLimits(spreadsheetID string, sheetTitle string, rows int, columns int) error {
	sheetID, exists, err := getSheetID(client, spreadsheetID, sheetTitle)
	if err != nil {
		return err
	}
	if !exists {
		return fmt.Errorf("sheet with title %s not found", sheetTitle)
	}

	return client.batchUpdate(spreadsheetID, []*sheets.Request{
		&sheets.Request{
			UpdateSheetProperties: &sheets.UpdateSheetPropertiesRequest{
				Fields: "*",
				Properties: &sheets.SheetProperties{
					SheetId: sheetID,
					Title:   sheetTitle,
					GridProperties: &sheets.GridProperties{
						ColumnCount: int64(columns),
						RowCount:    int64(rows),
					},
				},
			},
		},
	})
}
