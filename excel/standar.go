package excel

import (
	"github.com/tealeg/xlsx"
	"log"
	"github.com/coderlindacheng/balabalago/pair"
	"github.com/coderlindacheng/balabalago"
	"strings"
	"github.com/coderlindacheng/balabalago/time"
	"strconv"
)

const (
	MALE  = "男"
	FEMAL = "女"

	FILE_NAME = "评分标准.xlsx"

	SCORE_ROW_NAME = "分值" //评分标准表的字段名,一定要有这个字段,而且要特殊处理的

	UNIT_TYPE_TIME = "时间" //单位类型 时间
)

type UnitTypePolicy func(s string) int

type Sport struct {
	Name      string
	Sex       string
	Policy    UnitTypePolicy
	Result    *[]int
	Score     *[]pair.IntPair
	UniqueKey string
}

var Sports map[string]Sport = map[string]Sport{}

/*
	默认的单位解析策略

	parma:
	 	任意字符串
	return :
		int 解析后的值得
 */
func defaultUnitTypePolicy(s string) int {
	v, err := strconv.Atoi(s)
	if err != nil {
		log.Panicf("%s 你填的是 %s 不是数字", FILE_NAME, s)
	}
	return v
}

/*
	以时间为单位的解析策略,最大到分钟,最小号毫秒,把输入的时间格式解析成毫秒

	parma:
	 	时间格式 , 分"秒'毫秒x100
	return :
		int 解析后的毫秒值
 */
func timeUnitTypePolicy(s string) int {

	//2个index到最后一定要被查出来,所以不如一开始就查出来
	minIndex := strings.Index(s, "\"")
	secIndex := strings.Index(s, "'")

	//"和'号都没有,默认解析成分钟
	if minIndex == -1 && secIndex == -1 {
		v, err := strconv.Atoi(s)
		if err != nil {
			log.Panicf("%s 时间格式解析有误 时间格式应该为 分\"秒'毫秒x100 你填的分钟数是 %s 不是数字", FILE_NAME, s)
		}
		return v * time.MILLIS_PER_MINUTE
	}

	if minIndex == 0 || secIndex == 0 {
		log.Panicf("%s 时间格式解析有误 时间格式应该为 分\"秒'毫秒x100 \"和' 不能出现在第一位", FILE_NAME)
	}

	var minute int
	var second int
	var millis int

	//分
	if minIndex > 0 {
		minuteStr := s[0:minIndex]
		v, err := strconv.Atoi(minuteStr)
		if err != nil {
			log.Panicf("%s 时间格式解析有误 时间格式应该为 分\"秒'毫秒x100 你填的分钟数是 %s 不是数字", FILE_NAME, minuteStr)
		}
		minute = v
	}

	//秒
	if secIndex > 0 {
		if minIndex > 0 {
			if secIndex-minIndex <= 1 {
				log.Panicf("%s 时间格式解析有误 时间格式应该为 分\"秒'毫秒x100 '一定要出现在\"后面而且不能连续出现", FILE_NAME)
			}
			secondStr := s[minIndex+1:secIndex]
			v, err := strconv.Atoi(secondStr)
			if err != nil {
				log.Panicf("%s 时间格式解析有误 时间格式应该为 分\"秒'毫秒x100 你填的秒数是 %s 不是数字", FILE_NAME, secondStr)
			}
			second = v
		} else {
			secondStr := s[0:secIndex]
			v, err := strconv.Atoi(secondStr)
			if err != nil {
				log.Panicf("%s 时间格式解析有误 时间格式应该为 分\"秒'毫秒x100 你填的秒数是 %s 不是数字", FILE_NAME, secondStr)
			}
			second = v
		}
		second = second * time.MILLIS_PER_SECOND
	}

	//毫秒
	if length := len(s); secIndex < length-1 {
		millisStr := s[secIndex+1:length]
		v, err := strconv.Atoi(millisStr)
		if err != nil {
			log.Panicf("%s 时间格式解析有误 时间格式应该为 分\"秒'毫秒x100 你填的毫秒数是 %s 不是数字", FILE_NAME, millisStr)
		}
		second = v * 100
	}

	return minute + second + millis
}

func ReadSheet() {
	defer utils.Pause()
	fileName := "评分标准.xlsx"
	xlsx, err := xlsx.OpenFile(fileName)
	if err != nil {
		log.Panicf("读取评分标准表出错,请确认在该目录下存在%s \n err=%s", fileName, err)
	}
	sheets := xlsx.Sheets
	minCellsCount := CalMinCellsEatchRow(sheets)
	score := []int{}
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
					name, sex, policy := parseRowName(rowName)
					if name != SCORE_ROW_NAME {
						sport = &Sport{Name: name, Sex: sex, Policy: policy, UniqueKey: name + sex}
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

/*
	评分表的表头解析

	parma:
		表头名字
	return:
		name 组名字
		sex 性别
		policy 该组的解析策略
 */
func parseRowName(rowName string) (name string, sex string, policy UnitTypePolicy) {
	ss := strings.Split(rowName, "#")
	name = ss[0]
	if length := len(ss); length > 1 {
		sex = ss[1]
		if sex != "" && sex != MALE && sex != FEMAL {
			log.Panicf("%s 组名的格式应该是 运动名称#性别#单位 而性别只能填 %s 或者 %s", FILE_NAME, MALE, FEMAL)
		}
		if length > 2 {
			switch ss[2] {
			case UNIT_TYPE_TIME:
				policy = timeUnitTypePolicy
			default:
				policy = defaultUnitTypePolicy
			}

		}
	}
	if name == "" {
		log.Panicf("%s 有一个组名是空的", FILE_NAME)
	}

	return
}
