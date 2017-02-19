package excel

import (
	"github.com/tealeg/xlsx"
	"log"
	"github.com/coderlindacheng/balabalago/pair"
	"github.com/coderlindacheng/balabalago"
	"strings"
	"fmt"
)

const (
	male     = "男"
	femal    = "女"
	fileName = "评分标准.xlsx"
)

var Sports map[string]Sport = make(map[string]Sport)

type Sport struct {
	name   string
	sex    string
	result []int
	score  []pair.IntPair
}

func readStandard() {
	defer utils.Pause()
	fileName := "评分标准.xlsx"
	xlsx, err := xlsx.OpenFile(fileName)
	if err != nil {
		log.Panicf("读取评分标准表出错,请确认在该目录下存在%s \n err=%s", fileName, err)
	}
	sheets := xlsx.Sheets
	minCellsCount := CalMinCellsEatchRow(sheets)
	score := [minCellsCount]int{}
	for i, s := range sheets {
		//只读取第一个表
		if i > 0 {
			break
		}

		for _, r := range s.Rows {
			var rowName string
			var sport Sport = nil
			for j, c := range r.Cells {
				if j > minCellsCount {
					break
				}
				if j == 0 {
					rowName, _ = c.String()
				} else {
					if rowName == "分值" {
						cInt, err := c.Int()
						if err != nil {
							log.Panicf("%s 读取分值的时候出错了 err = %s", fileName, err)
						}
						score[j] = cInt
						continue
					} else {
						sport = Sport{}
						ss := strings.Split(rowName, "#")
						sport.name = ss[0]
						if length := len(ss); length > 1 {
							sport.sex = ss[1]
							if length > 2 {

							}
						}
						if sport.name == "" {
							log.Panicf("%s 有一个组名是空的", fileName)
						}

						if sport.sex != "" && sport.sex != male && sport.sex != femal {
							log.Panicf("%s 组名的格式应该是 运动名称#性别#单位 而性别只能填 %s 或者 %s", fileName, male, femal)
						}

					}
				}

			}
		}
	}
}
