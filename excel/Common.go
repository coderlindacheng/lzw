package excel

import (
	"github.com/tealeg/xlsx"
	"math"
)

func CalMinCellsEatchRow(sheets []*xlsx.Sheet) int {
	minCellsCount := math.MaxInt32
	for i, s := range sheets {
		if i > 0 {
			break
		}
		for _, r := range s.Rows {
			cellsCount := len(r.Cells)
			if minCellsCount > cellsCount {
				minCellsCount = cellsCount
			}
		}
	}

	return minCellsCount
}
