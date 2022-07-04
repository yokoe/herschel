package herschel

// Read returns a slice of cell values in sheet.
func (client *Client) Read(spreadsheetID string, sheetTitle string) ([][]interface{}, error) {
	resp, err := client.service.Spreadsheets.Values.Get(spreadsheetID, sheetTitle).Do()
	if err != nil {
		return nil, err
	}

	return resp.Values, nil
}

// ReadTable returns a table with values read from the spreadsheet set.
func (client *Client) ReadTable(spreadsheetID string, sheetTitle string) (*Table, error) {
	values, err := client.Read(spreadsheetID, sheetTitle)
	if err != nil {
		return nil, err
	}

	maxCols := 0
	for _, row := range values {
		cols := len(row)
		if cols > maxCols {
			maxCols = cols
		}
	}

	t := NewTable(len(values), maxCols)
	for i, row := range values {
		for j, col := range row {
			t.PutValue(i, j, col)
		}
	}
	return t, nil
}

// SheetTitles returns a slice of sheet titles.
func (client Client) SheetTitles(spreadsheetID string) ([]string, error) {
	return getSheetTitles(client, spreadsheetID)
}
