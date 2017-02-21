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
	MALE  = "男"
	FEMAL = "女"

	FILE_NAME = "评分标准.xlsx"

	SCORE_ROW_NAME = "分值" //评分标准表的字段名,一定要有这个字段,而且要特殊处理的

	UNIT_TYPE_TIME_STR = "时间" //单位类型 时间(填表的时候填这个)
	UNIT_TYPE_TIME     = 1    //单位类型 时间 (程序用)
)

type Sport struct {
	Name     string
	Sex      string
	UnitTypPolicy func(s *string)
	Result   *[]int
	Score    *[]pair.IntPair
}

var Sports map[string]Sport = map[string]Sport{}



func ReadSheet() {
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
			var sport *Sport
			for j, c := range r.Cells {
				if j > minCellsCount {
					break
				}
				if j == 0 {
					rowName, _ := c.String()
					name, sex, uniType := parseRowName(rowName)
					if name != SCORE_ROW_NAME {
						sport = &Sport{name: &name, sex: &sex}
					}
				} else {
					if sport == nil {
						cInt, err := c.Int()
						if err != nil {
							log.Panicf("%s 读取分值的时候出错了 err = %s", fileName, err)
						}
						score[j] = cInt
						continue
					} else {

					}
				}

			}
		}
	}
}


func parseRowName(rowName string) (name string, sex string, unitType int) {
	ss := strings.Split(rowName, "#")
	name = ss[0]
	if length := len(ss); length > 1 {
		sex = ss[1]
		if sex != "" && sex != MALE && sex != FEMAL {
			log.Panicf("%s 组名的格式应该是 运动名称#性别#单位 而性别只能填 %s 或者 %s", FILE_NAME, MALE, FEMAL)
		}
		if length > 2 {
			switch ut := ss[2]; ut {
			case UNIT_TYPE_TIME_STR:
				unitType = UNIT_TYPE_TIME
			default:
				log.Panicf("%s 组名的格式应该是 运动名称#性别#单位 单位类型是个不知道什么来的 %s", FILE_NAME, MALE, ut)
			}

		}
	}
	if name == "" {
		log.Panicf("%s 有一个组名是空的", FILE_NAME)
	}

	return
}
