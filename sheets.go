package herschel

import (
	sheets "google.golang.org/api/sheets/v4"
)

func getSheetID(client Client, spreadsheetID string, sheetName string) (int64, bool, error) {
	spreadsheet, err := client.service.Spreadsheets.Get(spreadsheetID).Do()
	if err != nil {
		return 0, false, err
	}

	for _, sheet := range spreadsheet.Sheets {
		if sheet.Properties.Title == sheetName {
			return sheet.Properties.SheetId, true, nil
		}
	}
	return 0, false, nil
}

func addSheet(client Client, spreadsheetID string, title string) error {
	req := sheets.Request{
		AddSheet: &sheets.AddSheetRequest{
			Properties: &sheets.SheetProperties{
				Title: title,
			},
		},
	}

	return client.batchUpdate(spreadsheetID, []*sheets.Request{&req})
}

func deleteSheetByID(client Client, spreadsheetID string, sheetID int64) error {
	req := sheets.Request{
		DeleteSheet: &sheets.DeleteSheetRequest{
			SheetId: sheetID,
		},
	}

	return client.batchUpdate(spreadsheetID, []*sheets.Request{&req})
}

func recreateSheet(client Client, spreadsheetID string, title string) error {
	sheetID, found, err := getSheetID(client, spreadsheetID, title)
	if err != nil {
		return err
	}
	if found {
		if err := deleteSheetByID(client, spreadsheetID, sheetID); err != nil {
			return err
		}
	}
	return addSheet(client, spreadsheetID, title)
}
