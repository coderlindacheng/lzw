package common

import (
	"github.com/tealeg/xlsx"
	"math"
	"github.com/coderlindacheng/balabalago/errors"
	"fmt"
)

//func CalMaxCellsCount(sheets []*xlsx.Sheet) int {
//	minCellsCount := math.MaxInt32
//	for i, s := range sheets {
//		if i > 0 {
//			break
//		}
//		for _, r := range s.Rows {
//			cellsCount := len(r.Cells)
//			if minCellsCount > cellsCount {
//				minCellsCount = cellsCount
//			}
//		}
//	}
//
//	return minCellsCount
//}

func ReadSheet(fileName string, readfunc func(maxCellsCount int) func(rowNum, cellNum, maxCellsCount int, s *xlsx.Sheet, r *xlsx.Row, c *xlsx.Cell) error) error {
	xlsx, err := xlsx.OpenFile(fileName)
	if err != nil {
		return errors.NewWrapper(err, fmt.Sprintf("读取评分标准表出错,请确认在该目录下存在 %s", fileName))
	}
	sheets := xlsx.Sheets
	maxCellsCount := CalMaxCellsCount(sheets)

	for i, s := range sheets {
		//只读取第一个表
		if i > 0 {
			break
		}

		for rowNum, r := range s.Rows {
			realReadFunc := readfunc(maxCellsCount)
			for cellNum, c := range r.Cells {
				if cellNum > maxCellsCount {
					break
				}
				err := realReadFunc(rowNum, cellNum, maxCellsCount, s, r, c)
				if err != nil {
					return err
				}
			}
		}
	}

	return nil
}
