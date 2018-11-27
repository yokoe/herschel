package herschel

import (
	"fmt"
	"image/color"

	sheets "google.golang.org/api/sheets/v4"
)

func (client Client) setCellFormats(spreadsheetID string, sheetName string, table *Table) error {
	sheetID, exists, err := getSheetID(client, spreadsheetID, sheetName)
	if err != nil {
		return err
	}
	if !exists {
		return fmt.Errorf("sheet not found with name: %s", sheetName)
	}
	// Background color
	return client.updateCellFormats(spreadsheetID, sheetID, table)
}

func (client Client) updateCellFormats(spreadsheetID string, sheetID int64, table *Table) error {
	requests := []*sheets.Request{}

	if table.FrozenRowCount > 0 {
		req := sheets.Request{
			UpdateSheetProperties: &sheets.UpdateSheetPropertiesRequest{
				Properties: &sheets.SheetProperties{
					SheetId: sheetID,
					GridProperties: &sheets.GridProperties{
						FrozenRowCount: table.FrozenRowCount,
					},
				},
				Fields: "gridProperties.frozenRowCount",
			},
		}
		requests = append(requests, &req)
	}

	if table.FrozenColumnCount > 0 {
		req := sheets.Request{
			UpdateSheetProperties: &sheets.UpdateSheetPropertiesRequest{
				Properties: &sheets.SheetProperties{
					SheetId: sheetID,
					GridProperties: &sheets.GridProperties{
						FrozenColumnCount: table.FrozenColumnCount,
					},
				},
				Fields: "gridProperties.frozenColumnCount",
			},
		}
		requests = append(requests, &req)
	}

	for col := 0; col < table.cols; col++ {
		for row := 0; row < table.rows; row++ {
			c := table.getBackgroundColor(row, col)
			if c != color.Transparent {
				r, g, b, a := c.RGBA()
				req := sheets.Request{
					RepeatCell: &sheets.RepeatCellRequest{
						Range: &sheets.GridRange{SheetId: sheetID,
							StartColumnIndex: int64(col),
							EndColumnIndex:   int64(col) + 1,
							StartRowIndex:    int64(row),
							EndRowIndex:      int64(row) + 1,
						},
						Cell: &sheets.CellData{
							UserEnteredFormat: &sheets.CellFormat{
								BackgroundColor: &sheets.Color{Alpha: float64(a) / 65536, Blue: float64(b) / 65536, Green: float64(g) / 65536, Red: float64(r) / 65536},
							},
						},
						Fields: "userEnteredFormat(backgroundColor)",
					},
				}
				requests = append(requests, &req)
			}

			p := table.getNumberFormatPattern(row, col)
			if len(p) > 0 {
				formatType := table.getNumberFormatType(row, col)
				if len(formatType) == 0 {
					formatType = "NUMBER"
				}
				req := sheets.Request{
					RepeatCell: &sheets.RepeatCellRequest{
						Range: &sheets.GridRange{SheetId: sheetID,
							StartColumnIndex: int64(col),
							EndColumnIndex:   int64(col) + 1,
							StartRowIndex:    int64(row),
							EndRowIndex:      int64(row) + 1,
						},
						Cell: &sheets.CellData{
							UserEnteredFormat: &sheets.CellFormat{
								NumberFormat: &sheets.NumberFormat{
									Type:    formatType,
									Pattern: p,
								},
							},
						},
						Fields: "userEnteredFormat(numberFormat)",
					},
				}
				requests = append(requests, &req)
			}
		}
	}

	return client.batchUpdate(spreadsheetID, requests)
}
