package common

import (
	"github.com/tealeg/xlsx"
	"math"
	"github.com/coderlindacheng/balabalago/errors"
	"fmt"
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

func ReadSheet(fileName string, function func(i, j, k int, s *xlsx.Sheet, r *xlsx.Row, c *xlsx.Cell) error) error {
	xlsx, err := xlsx.OpenFile(fileName)
	if err != nil {
		return errors.NewWrapper(err, fmt.Sprintf("读取评分标准表出错,请确认在该目录下存在%s \n err=%s \n", fileName))
	}
	sheets := xlsx.Sheets
	minCellsCount := CalMinCellsEatchRow(sheets)

	for i, s := range sheets {
		//只读取第一个表
		if i > 0 {
			break
		}

		for j, r := range s.Rows {
			for k, c := range r.Cells {
				if j > minCellsCount {
					break
				}
				err := function(i, j, k, s, r, c)
				if err != nil {
					return err
				}
			}
		}
	}

	return nil
}
