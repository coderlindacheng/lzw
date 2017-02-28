package source

import (
	"fmt"
	"github.com/tealeg/xlsx"
	"github.com/coderlindacheng/lzw/excel/common"
	"github.com/coderlindacheng/balabalago/errors"
)

func Read(fileName string, ) (f common.ReadSheetFunc, err error) {
	var recores [][]string
	var rowNames []string
	var maxRow int
	var maxCell int
	f = func(rowNum, cellNum int, s *xlsx.Sheet, r *xlsx.Row, c *xlsx.Cell) error {
		if rowNum == 0 {
			if recores == nil {
				maxRow = s.MaxRow
				maxCell = len(r.Cells)
				recores = make([][]string, maxCell)
				rowNames = make([]string, maxCell)
			}

			if recores[cellNum] == nil {
				recores[cellNum] = make([]string, maxRow)
			}

			s, err := c.String()
			if err != nil {
				return errors.NewWrapper(err, fmt.Sprintf("%s 读取数据时出错了 读到的是 %v ", fileName, s))
			}
			recores[cellNum][rowNum] = s
			rowNames[cellNum] = s
		} else {
			rowName:=rowNames[cellNum]
			s, err := c.String()
			if err != nil {
				return errors.NewWrapper(err, fmt.Sprintf("%s 读取数据时出错了 读到的是 %v ", fileName, s))
			}
		}

		return nil
	}

	return
}
