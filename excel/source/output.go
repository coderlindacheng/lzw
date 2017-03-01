package source

import (
	"fmt"
	"github.com/coderlindacheng/balabalago/errors"
	"os"
	"github.com/tealeg/xlsx"
)

const OUTPUT_FILE_NAME = "./分数表.xlsx"

func Output(fileName string) error {

	if _, err := os.Stat(fileName); err == nil || os.IsExist(err) {
		err := os.RemoveAll(fileName)
		if err != nil {
			return errors.NewWrapper(err, fmt.Sprintf("试图删除 %s 文件 但是出错了", fileName))
		}
	}

	file := xlsx.NewFile()
	for i, pair := range Datas {
		if pair != nil {
			s, ok := pair.Left().(string)
			if ok {
				sheet, err := file.AddSheet(s)
				if err != nil {
					return errors.NewWrapper(err, fmt.Sprintf("创建第 %v sheet(%s) 表名重复了", i, s))
				}
				ss, ok := pair.Right().([][]string)
				if ok {
					for _, rowOut := range ss {
						row := sheet.AddRow()
						for _, cellOut := range rowOut {
							cell := row.AddCell()
							cell.Value = cellOut
						}
					}
				}
			} else {
				return errors.NewOnlyStr(fmt.Sprintf("创建第 %v 个sheet的时候,sheet的名称读取错误读到个%s", i, s))
			}

		}

	}

	err := file.Save(fileName)
	if err != nil {
		return errors.NewWrapper(err, fmt.Sprintf("保存文件 %s 的时候出错了,尼玛最后一步啦,能不能好好的", fileName))
	}
	return nil
}
