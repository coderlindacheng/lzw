package standar

import (
	"github.com/tealeg/xlsx"
	"strings"
	. "github.com/coderlindacheng/balabalago/special_string"
	"fmt"
	"github.com/coderlindacheng/balabalago/errors"
	. "github.com/coderlindacheng/balabalago/types"
	"github.com/coderlindacheng/lzw/excel/common"
	"strconv"
)

const (
	MALE  = "男"
	FEMAL = "女"

	FILE_NAME = "评分标准.xlsx"

	SCORE_ROW_NAME = "分值" //评分标准表的字段名,一定要有这个字段,而且要特殊处理的

	UNIT_TYPE_TIME = "时间" //单位类型 时间

	SCORT_DENOMINATOR = 1000000 //分数的先扩大,再缩小的倍数,最多只有100分,所以多点0,最后算出来的结果精确点

	DESC_SORTING = 0 //降序排序
	ASC_SORTING  = 1 //升序排序

	DESC_SORTING_STR = "降序" //降序排序
	ASC_SORTING_STR  = "升序" //升序排序

)

type Sport struct {
	Name            string
	Sex             string
	Policy          UnitTypePolicy
	ScoreTranslator []TripleInt //left:底分 mid:输入 right:没一点数值的分差 (第一个是100分哦)
	UniqueKey       string
	sorting         int
}

func (p *Sport) OutPut(fileName, s string, sheetNum, rowNum, cellNum int) (string, error) {
	v, err := p.Policy(fileName, s, sheetNum, rowNum, cellNum)
	if err != nil {
		return s, err
	}

	f := func(tripleInt TripleInt, dv int) (string, error) {
		var finalScore int
		if dv > 0 {
			finalScore = tripleInt.Left() + dv*tripleInt.Right()
		} else if dv == 0 {
			finalScore = tripleInt.Left()
		} else {
			return s, errors.NewOnlyStr(fmt.Sprintf("%s 第%v行第%v列 出现了一些奇怪的问题", fileName, rowNum, cellNum))
		}
		return strconv.FormatFloat(float64(finalScore)/SCORT_DENOMINATOR, 'f', 2, 32), nil
	}

	if p.sorting == ASC_SORTING {
		for _, tripleInt := range p.ScoreTranslator {
			if v <= tripleInt.Mid() {
				return f(tripleInt, tripleInt.Mid()-v)
			}
		}
	} else {
		for _, tripleInt := range p.ScoreTranslator {
			if v >= tripleInt.Mid() {
				return f(tripleInt, v-tripleInt.Mid())
			}
		}

	}

	return s, errors.NewOnlyStr(fmt.Sprintf("%s 第%v行第%v列 出现了一些奇怪的问题", fileName, rowNum, cellNum))
}

var Sports map[string]Sport = map[string]Sport{}

/*
	默认的单位解析策略

	parma:
	 	任意字符串
	return :
		int 解析后的值得
 */

func Read(fileName string) (f common.ReadSheetFunc, err error) {
	var realMaxRow int
	var score []int
	var dscore []int
	var results [][]int
	var sports []*Sport
	var maxCell int
	f = func(sheetNum, rowNum, cellNum int, s *xlsx.Sheet, r *xlsx.Row, c *xlsx.Cell) error {
		if sheetNum > 1 {
			return errors.NewOnlyStr(fmt.Sprintf("%s 只读一个sheet", fileName))
		}

		if rowNum == 0 {
			if sports == nil {
				maxRow := s.MaxRow
				if maxRow < 3 {
					err = errors.NewOnlyStr(fmt.Sprintf("%s 表至少要有3行 表头和2行数据,要不没法评分", fileName))
				}
				realMaxRow = maxRow - 1 //因为第一行是表头,不能计算成数据,所以要减掉1
				score = make([]int, realMaxRow)
				dscore = make([]int, realMaxRow-1)
				maxCell = len(r.Cells)
				results = make([][]int, maxCell)
				sports = make([]*Sport, maxCell)
			}

			rowName, err := c.String()
			if err != nil {
				return errors.NewWrapper(err, fmt.Sprintf("%s 读取表头的时候出错了 读到%s", fileName, rowName))
			}
			name, sex, policy, sorting, err := parseRowName(fileName, rowName)
			if err != nil {
				return err
			}
			if name != SCORE_ROW_NAME {
				sports[cellNum] = &Sport{Name: name, Sex: sex, Policy: policy, ScoreTranslator: make([]TripleInt, realMaxRow), UniqueKey: name + POUND_SIGN + sex, sorting: sorting}
			}
		} else {
			sport := sports[cellNum]
			realRowNum := rowNum - 1
			if sport == nil {
				cInt, err := c.Int()
				if err != nil {
					return errors.NewWrapper(err, fmt.Sprintf("%s 读取分值的时候出错了", fileName))
				}
				score[realRowNum] = cInt * SCORT_DENOMINATOR
				if realRowNum > 0 {
					dscore[realRowNum-1] = score[realRowNum-1] - score[realRowNum]
				}
			} else {
				cStr, err := c.String()
				if err != nil {
					return errors.NewWrapper(err, fmt.Sprintf("%s %s 读取数据的时候出错了", fileName, sport.UniqueKey))
				}

				v, err := sport.Policy(fileName, cStr, sheetNum, rowNum, cellNum)
				if err != nil {
					return err
				}
				if results[cellNum] == nil {
					results[cellNum] = make([]int, realMaxRow)
				}

				results[cellNum][realRowNum] = v

			}
		}

		if rowNum == realMaxRow && cellNum+1 >= maxCell {
			for i, sport := range sports {
				if sport == nil {
					continue
				}
				var result []int = results[i]
				sport.ScoreTranslator[0] = TripleInt{score[0], result[0], 0}
				for j := 1; j < realMaxRow; j++ {
					currentResult := result[j]
					lastIndex := j - 1
					var dresule int
					var sortingStr string
					if sport.sorting == ASC_SORTING {
						dresule = currentResult - result[lastIndex]
						sortingStr = ASC_SORTING_STR
					} else {
						dresule = result[lastIndex] - currentResult
						sortingStr = DESC_SORTING_STR
					}
					if dresule <= 0 {
						return errors.NewOnlyStr(fmt.Sprintf("%s 第 %v 列所有的值必须是%s排序 问题出在第 %v 行", fileName, i, sortingStr, j))

					}
					perSR := dscore[lastIndex] / dresule
					sport.ScoreTranslator[j] = TripleInt{score[j], result[j], perSR}
				}
				Sports[sport.UniqueKey] = *sport
			}
		}
		return nil
	}

	return
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
func parseRowName(fileName, rowName string) (name string, sex string, policy UnitTypePolicy, sorting int, err error) {
	rowName = strings.TrimSpace(rowName)
	splitedStrArry := strings.Split(rowName, POUND_SIGN)
	name = splitedStrArry[0]
	policy = defaultUnitTypePolicy //默认字符解析策略
	sorting = DESC_SORTING         //默认降序排序

	if length := len(splitedStrArry); length > 1 {
		sex = splitedStrArry[1]
		if sex == "" || (sex != MALE && sex != FEMAL) {
			err = errors.NewOnlyStr(fmt.Sprintf("%s 组名的格式应该是 运动名称%s性别%s单位 而性别只能填 %s 或者 %s 你填的却是 %s", fileName, POUND_SIGN, POUND_SIGN, MALE, FEMAL, sex))
			return
		}
		if length > 2 {
			switch ut := splitedStrArry[2]; ut {
			case UNIT_TYPE_TIME:
				policy = timeUnitTypePolicy
				sorting = ASC_SORTING
			default:
				err = errors.NewOnlyStr(fmt.Sprintf("%s 不存在这个单位类型 %s 现在可填的只有 %s 这个类型", fileName, ut, UNIT_TYPE_TIME))
				return
			}

			if length > 3 {
				switch sortingStr := splitedStrArry[3]; sortingStr {
				case ASC_SORTING_STR:
					sorting = ASC_SORTING
				case DESC_SORTING_STR:
					sorting = DESC_SORTING
				default:
					err = errors.NewOnlyStr(fmt.Sprintf("%s 不存在这个排序类型 %s 而且注意所有的类型都要大写 升序(%s) 降序(%s)", fileName, sortingStr, ASC_SORTING_STR, DESC_SORTING_STR))
					return
				}

			}
		}

	}
	if name == "" {
		err = errors.NewOnlyStr(fmt.Sprintf("%s 有一个组名是空的", fileName))
		return
	}

	return
}
