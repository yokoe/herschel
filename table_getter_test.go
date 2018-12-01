package herschel

import (
	"testing"
)

func TestTable_GetStringValue(t *testing.T) {
	table := NewTable(1, 4)
	table.PutValue(0, 0, "Hello")
	table.PutValue(0, 1, "World")
	table.PutValue(0, 2, 123)

	type args struct {
		row int
		col int
	}
	tests := []struct {
		name string
		args args
		want string
	}{
		{"StringAtFirstCell", args{0, 0}, "Hello"},
		{"StringAtSecondCell", args{0, 1}, "World"},
		{"CellWithIntValue", args{0, 2}, ""},
		{"EmptyCell", args{0, 3}, ""},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := table.GetStringValue(tt.args.row, tt.args.col); got != tt.want {
				t.Errorf("Table.GetStringValue() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestTable_GetIntValue(t *testing.T) {
	table := NewTable(1, 4)
	table.PutValue(0, 0, 123)
	table.PutValue(0, 1, 456)
	table.PutValue(0, 2, "Hello")

	type args struct {
		row int
		col int
	}
	tests := []struct {
		name string
		args args
		want int
	}{
		{"IntAtFirstCell", args{0, 0}, 123},
		{"IntAtSecondCell", args{0, 1}, 456},
		{"CellWithStringValue", args{0, 2}, 0},
		{"EmptyCell", args{0, 3}, 0},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := table.GetIntValue(tt.args.row, tt.args.col); got != tt.want {
				t.Errorf("Table.GetIntValue() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestTable_GetInt64Value(t *testing.T) {
	table := NewTable(1, 4)
	table.PutValue(0, 0, int64(9223372036854775806))
	table.PutValue(0, 1, 456)
	table.PutValue(0, 2, "Hello")

	type args struct {
		row int
		col int
	}
	tests := []struct {
		name string
		args args
		want int64
	}{
		{"Int64AtFirstCell", args{0, 0}, 9223372036854775806},
		{"IntAtSecondCell", args{0, 1}, 456},
		{"CellWithStringValue", args{0, 2}, 0},
		{"EmptyCell", args{0, 3}, 0},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := table.GetInt64Value(tt.args.row, tt.args.col); got != tt.want {
				t.Errorf("Table.GetInt64Value() = %v, want %v", got, tt.want)
			}
		})
	}
}
