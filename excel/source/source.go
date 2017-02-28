package source

import (
	"fmt"
	"github.com/tealeg/xlsx"
	"github.com/coderlindacheng/lzw/excel/common"
	"github.com/coderlindacheng/balabalago/errors"
	"github.com/coderlindacheng/balabalago/special_string"
	"github.com/coderlindacheng/lzw/excel/standar"
)

const SEX_MARK = "性别"

func Read(fileName string, ) (f common.ReadSheetFunc, err error) {
	var recores [][]string
	var rowNames []string
	var maxRow int
	var maxCell int
	var current []string
	var sex string
	f = func(rowNum, cellNum int, s *xlsx.Sheet, r *xlsx.Row, c *xlsx.Cell) error {
		if rowNum == 0 {
			if recores == nil {
				sex = ""
				maxRow = s.MaxRow
				maxCell = len(r.Cells)
				recores = make([][]string, maxRow)
				rowNames = make([]string, maxCell)

				for i := 0; i < maxRow; i++ {
					recores[i] = make([]string, maxCell)
				}
				rowNames = recores[0]
			}

			s, err := c.String()
			if err != nil {
				return errors.NewWrapper(err, fmt.Sprintf("%s 读取数据时出错了 读到的是 %s ", fileName, s))
			}
			rowNames[cellNum] = s
			if s == SEX_MARK {
				sex = SEX_MARK
			}
			if cellNum+1 == maxCell && sex == "" {
				return errors.NewOnlyStr(fmt.Sprintf("%s 表头没有 %s 这个字段", fileName, SEX_MARK))
			}
		} else {
			if cellNum == 0 {
				current = recores[rowNum]
				sex = ""
			}

			rowName := rowNames[cellNum]
			s, err := c.String()
			if err != nil {
				return errors.NewWrapper(err, fmt.Sprintf("%s 读取数据时出错了 读到的是 %s ", fileName, s))
			}
			if sex != "" {
				key := rowName + special_string.POUND_SIGN + sex
				sprot, ok := standar.Sports[key]
				if ok {
					v,err:=sprot.Policy(fileName,s)

				} else {
					return errors.NewWrapper(err, fmt.Sprintf("%s 读取数据时出错了  %v 在 %s 中读不到数据", fileName, current, rowName))
				}
			} else {
				if rowName == SEX_MARK {
					sex = s
				}
			}
			current[cellNum] = s
		}

		return nil
	}

	return
}
