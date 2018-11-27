package herschel

import (
	sheets "google.golang.org/api/sheets/v4"
)

// CreateNewSpreadsheet creates new spreadsheet
func (c *Client) CreateNewSpreadsheet(title string) (string, error) {
	resp, err := c.service.Spreadsheets.Create(&sheets.Spreadsheet{Properties: &sheets.SpreadsheetProperties{
		Title: title,
	}}).Do()
	if err != nil {
		return "", err
	}
	return resp.SpreadsheetId, nil
}
