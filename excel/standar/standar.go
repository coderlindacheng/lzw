package standar

import (
	"github.com/tealeg/xlsx"
	"strings"
	. "github.com/coderlindacheng/balabalago/special_string"
	"fmt"
	"github.com/coderlindacheng/balabalago/errors"
	. "github.com/coderlindacheng/balabalago/types"
)

const (
	MALE  = "男"
	FEMAL = "女"

	FILE_NAME = "评分标准.xlsx"

	SCORE_ROW_NAME = "分值" //评分标准表的字段名,一定要有这个字段,而且要特殊处理的

	UNIT_TYPE_TIME = "时间" //单位类型 时间

	SCORT_DENOMINATOR = 1000000 //分数的先扩大,再缩小的倍数,最多只有100分,所以多点0,最后算出来的结果精确点
)

type Sport struct {
	Name            string
	Sex             string
	Policy          UnitTypePolicy
	ScoreTranslator []TripleInt
	UniqueKey       string
}

var Sports map[string]Sport = map[string]Sport{}

/*
	默认的单位解析策略

	parma:
	 	任意字符串
	return :
		int 解析后的值得
 */

func Read(maxCellsCount int) func(rowNum, cellNum, maxCellsCount int, s *xlsx.Sheet, r *xlsx.Row, c *xlsx.Cell) error {
	score := make([]int, maxCellsCount)
	result := make([]int, maxCellsCount)
	var sport *Sport
	return func(rowNum, cellNum, maxCellsCount int, s *xlsx.Sheet, r *xlsx.Row, c *xlsx.Cell) error {
		if rowNum == 0 {
			rowName, err := c.String()
			if err != nil {
				return errors.NewWrapper(err, fmt.Sprintf("%s 读取表头的时候出错了 读到%s", FILE_NAME, rowName))
			}
			name, sex, policy, err := parseRowName(rowName)
			if err != nil {
				return err
			}
			if name != SCORE_ROW_NAME {
				sport = &Sport{Name: name, Sex: sex, Policy: policy, ScoreTranslator: make([]TripleInt, maxCellsCount), UniqueKey: name + POUND_SIGN + sex}
			} else {
				sport = nil
			}
		} else {
			if sport == nil {
				cInt, err := c.Int()
				if err != nil {
					return errors.NewWrapper(err, fmt.Sprintf("%s 读取分值的时候出错了", FILE_NAME))
				}
				score[cellNum] = cInt * SCORT_DENOMINATOR
			} else {
				cStr, err := c.String()
				if err != nil {
					return errors.NewWrapper(err, fmt.Sprintf("%s %s 读取数据的时候出错了", FILE_NAME, sport.UniqueKey))
				}

				v, err := sport.Policy(cStr)
				if err != nil {
					return err
				}

				result[cellNum] = v
			}
		}

		if cellNum == maxCellsCount {
			sport.ScoreTranslator[maxCellsCount] = TripleInt{score[maxCellsCount], result[maxCellsCount], 0}
			for i := maxCellsCount; i > 0; i-- {
				currentScore := score[i]
				currentResult := result[i]
				nextIndex := i - 1
				dscore := currentScore - score[nextIndex]
				dresule := currentResult - result[nextIndex]
				perSR := dscore / dresule
				sport.ScoreTranslator[nextIndex] = TripleInt{score[nextIndex], result[nextIndex], perSR}
			}
			Sports[sport.UniqueKey] = *sport
		}
		return nil
	}
}

/*
	评分表的表头解析

	parma:
		表头名字
	return:
		name 组名字
		sex 性别
		policy 该组的解析策略
 */
func parseRowName(rowName string) (name string, sex string, policy UnitTypePolicy, err error) {
	ss := strings.Split(rowName, POUND_SIGN)
	name = ss[0]
	policy = defaultUnitTypePolicy
	if length := len(ss); length > 1 {
		sex = ss[1]
		if sex == "" || (sex != MALE && sex != FEMAL) {
			err = errors.NewOnlyStr(fmt.Sprintf("%s 组名的格式应该是 运动名称%s性别%s单位 而性别只能填 %s 或者 %s 你填的却是 %s", FILE_NAME, POUND_SIGN, POUND_SIGN, MALE, FEMAL, sex))
			return
		}
		if length > 2 {
			switch ss[2] {
			case UNIT_TYPE_TIME:
				policy = timeUnitTypePolicy
			default:
			//do nothing
			}

		}

	}
	if name == "" {
		err = errors.NewOnlyStr(fmt.Sprintf("%s 有一个组名是空的", FILE_NAME))
		return
	}

	return
}
