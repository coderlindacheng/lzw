package source

import (
	"fmt"
	"github.com/tealeg/xlsx"
	"github.com/coderlindacheng/lzw/excel/common"
	"github.com/coderlindacheng/balabalago/errors"
	"github.com/coderlindacheng/balabalago/special_string"
	"github.com/coderlindacheng/lzw/excel/standar"
	"github.com/coderlindacheng/balabalago/types"
	"strconv"
)

const (
	SEX_MARK        = "性别"
	INPUT_FILE_NAME = "./原始表.xlsx"
	SHEET_PREFIX    = "工作表"
)

var datasToOutput []*types.Pair = make([]*types.Pair, 0, 20) //随便定义一个默认值
func DatasToOutput() []*types.Pair {
	return datasToOutput
}

var sheetNameMap map[string]string = map[string]string{}

func Read(fileName string) (f common.ReadSheetFunc, err error) {
	var recores [][]string
	var rowNames []string
	var maxRow int
	var maxCell int
	var current []string
	var sex string
	var currentSheetNum int

	f = func(sheetNum, rowNum, cellNum int, s *xlsx.Sheet, r *xlsx.Row, c *xlsx.Cell) error {
		if rowNum == 0 {
			if recores == nil || currentSheetNum != sheetNum {
				sex = ""
				maxRow = s.MaxRow
				maxCell = len(r.Cells)
				recores = make([][]string, maxRow)
				rowNames = make([]string, maxCell)

				for i := 0; i < maxRow; i++ {
					recores[i] = make([]string, maxCell)
				}
				rowNames = recores[0]
				currentSheetNum = sheetNum
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

			key := rowName + special_string.POUND_SIGN + sex
			sprot, ok := standar.Sports()[key]
			if ok {
				outPutStr, err := sprot.OutPut(fileName, s, sheetNum, rowNum, cellNum)
				if err != nil {
					return err
				}
				current[cellNum] = outPutStr
			} else {
				current[cellNum] = s
			}

			if rowName == SEX_MARK && sex == "" {
				sex = s
			}

		}

		if rowNum == maxRow-1 && cellNum == maxCell-1 {
			sheetName := s.Name
			if _, exist := sheetNameMap[sheetName]; exist {
				sheetName = SHEET_PREFIX + special_string.SPACE + strconv.Itoa(len(sheetNameMap))
			}
			sheetNameMap[sheetName] = sheetName
			datasToOutput = append(datasToOutput, &types.Pair{sheetName, recores})
		}

		return nil
	}

	return
}
