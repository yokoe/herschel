package herschel

import (
	"encoding/csv"
	"io"
	"strconv"
)

// ToCSV writes table in csv format
func (t *Table) ToCSV(w io.Writer) error {
	cw := csv.NewWriter(w)
	defer cw.Flush()
	for i := 0; i < t.GetRows(); i++ {
		strValues := []string{}
		for j := 0; j < t.GetCols(); j++ {
			s := ""
			v := t.GetValue(i, j)
			if v != nil {
				if strValue, ok := v.(string); ok {
					s = strValue
				} else if i, ok := v.(int); ok {
					s = strconv.Itoa(i)
				} else if i, ok := v.(int64); ok {
					s = strconv.FormatInt(i, 10)
				}
			}
			strValues = append(strValues, s)
		}
		cw.Write(strValues)
	}
	return cw.Error()
}
