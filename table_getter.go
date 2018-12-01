package herschel

// GetStringValue returns value of cell as string.
func (t *Table) GetStringValue(row int, col int) string {
	v := t.GetValue(row, col)
	if v != nil {
		if s, ok := v.(string); ok {
			return s
		}
	}
	return ""
}

// GetIntValue returns value of cell as int.
func (t *Table) GetIntValue(row int, col int) int {
	v := t.GetValue(row, col)
	if v != nil {
		if i, ok := v.(int); ok {
			return i
		}
	}
	return 0
}

// GetInt64Value returns value of cell as int.
func (t *Table) GetInt64Value(row int, col int) int64 {
	v := t.GetValue(row, col)
	if v != nil {
		if i, ok := v.(int64); ok {
			return i
		}
		if i, ok := v.(int); ok {
			return int64(i)
		}
	}
	return 0
}
