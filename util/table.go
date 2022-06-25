package util

import "fmt"

type Column struct {
	Name       string
	RightAlign bool
}

type Table struct {
	Columns []Column
}

func makeFormatString(tableConfig *Table, columnWidths []int) string {
	formatString := ""

	for i, column := range tableConfig.Columns {
		width := columnWidths[i]
		if !column.RightAlign {
			width = -width
		}
		formatString = formatString + fmt.Sprintf(" %%%ds", width)
	}

	return formatString[1:] + "\n"
}

func computeColumnWidths(tableConfig *Table, tableData [][]string) []int {
	columnWidths := make([]int, len(tableConfig.Columns))

	// Default to making column wide enough for the heading.
	for i, column := range tableConfig.Columns {
		columnWidths[i] = len(column.Name)
	}

	for _, row := range tableData {
		for i, col := range row {
			if len(col) > columnWidths[i] {
				columnWidths[i] = len(col)
			}
		}
	}

	return columnWidths
}

// TODO: Add some sorting functionality...
func (tableConfig *Table) Show(tableData [][]string) {
	columnWidths := computeColumnWidths(tableConfig, tableData)
	formatString := makeFormatString(tableConfig, columnWidths)

	headers := make([]any, len(tableConfig.Columns))
	for idx, col := range tableConfig.Columns {
		headers[idx] = col.Name
	}

	fmt.Printf(formatString, headers...)

	for _, row := range tableData {
		anyRow := make([]any, len(row))
		for i, col := range row {
			anyRow[i] = col
		}
		fmt.Printf(formatString, anyRow...)
	}
}
