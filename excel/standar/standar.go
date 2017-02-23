package standar

import (
	"github.com/tealeg/xlsx"
	"log"
	"github.com/coderlindacheng/balabalago/pair"
	"strings"
	"github.com/coderlindacheng/balabalago/time"
	"strconv"
	"github.com/coderlindacheng/balabalago/special_string"
	"fmt"
	"github.com/coderlindacheng/balabalago/errors"
)

const (
	MALE  = "男"
	FEMAL = "女"

	FILE_NAME = "评分标准.xlsx"

	SCORE_ROW_NAME = "分值" //评分标准表的字段名,一定要有这个字段,而且要特殊处理的

	UNIT_TYPE_TIME = "时间" //单位类型 时间
)

type UnitTypePolicy func(s string) (int, error)

type Sport struct {
	Name      string
	Sex       string
	Policy    UnitTypePolicy
	Result    []int
	Score     []pair.IntPair
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
func defaultUnitTypePolicy(s string) (int, error) {
	v, err := strconv.Atoi(s)
	if err != nil {
		log.Printf("%s 你填的是 %s 不是数字", FILE_NAME, s)
		return 0, err
	}
	return v, nil
}

/*
	以时间为单位的解析策略,最大到分钟,最小号毫秒,把输入的时间格式解析成毫秒

	parma:
	 	时间格式 , 分"秒%s毫秒x100
	return :
		int 解析后的毫秒值
 */
func timeUnitTypePolicy(s string) (int, error) {

	//2个index到最后一定要被查出来,所以不如一开始就查出来
	minIndex := strings.Index(s, special_string.QUOTE)
	secIndex := strings.Index(s, special_string.SINGLE_QUOTE)

	//"和%s号都没有,默认解析成分钟
	if minIndex == -1 && secIndex == -1 {
		v, err := strconv.Atoi(s)
		if err != nil {
			errStr:=fmt.Sprintf("%s 时间格式解析有误 时间格式应该为 分%s秒%s毫秒x100 你填的分钟数是 %s 不是数字", FILE_NAME, special_string.QUOTE, special_string.SINGLE_QUOTE, s)
			return 0, &errors.ErrorWrapper{&err, &errStr}
		}
		return v * time.MILLIS_PER_MINUTE, nil
	}

	if minIndex == 0 || secIndex == 0 {
		return 0, &errors.ErrorWrapper{Attach: &(fmt.Sprintf("%s 时间格式解析有误 时间格式应该为 分%s秒%s毫秒x100 %s和%s 不能出现在第一位", FILE_NAME, special_string.QUOTE, special_string.SINGLE_QUOTE, special_string.QUOTE, special_string.SINGLE_QUOTE))}
	}

	var minute int
	var second int
	var millis int

	//分
	if minIndex > 0 {
		minuteStr := s[0:minIndex]
		v, err := strconv.Atoi(minuteStr)
		if err != nil {
			return 0, &errors.ErrorWrapper{&err, &(fmt.Sprintf("%s 时间格式解析有误 时间格式应该为 分%s秒%s毫秒x100 你填的分钟数是 %s 不是数字", FILE_NAME, special_string.QUOTE, special_string.SINGLE_QUOTE, minuteStr))}
		}
		minute = v * time.MILLIS_PER_MINUTE
	}

	//秒
	if secIndex > 0 {
		if minIndex > 0 {
			if secIndex-minIndex <= 1 {
				return 0, &errors.ErrorWrapper{Attach: &(fmt.Sprintf("%s 时间格式解析有误 时间格式应该为 分%s秒%s毫秒x100 %s一定要出现在%s后面而且不能连续出现", FILE_NAME, special_string.QUOTE, special_string.SINGLE_QUOTE, special_string.SINGLE_QUOTE, special_string.QUOTE))}
			}
			secondStr := s[minIndex+1:secIndex]
			v, err := strconv.Atoi(secondStr)
			if err != nil {
				return 0, &errors.ErrorWrapper{&err, &(fmt.Sprintf("%s 时间格式解析有误 时间格式应该为 分%s秒%s毫秒x100 你填的秒数是 %s 不是数字", FILE_NAME, special_string.QUOTE, special_string.SINGLE_QUOTE, secondStr))}
			}
			second = v
		} else {
			secondStr := s[0:secIndex]
			v, err := strconv.Atoi(secondStr)
			if err != nil {
				return 0, &errors.ErrorWrapper{&err, &(fmt.Sprintf("%s 时间格式解析有误 时间格式应该为 分%s秒%s毫秒x100 你填的秒数是 %s 不是数字", FILE_NAME, special_string.QUOTE, special_string.SINGLE_QUOTE, secondStr))}
			}
			second = v
		}
		second = second * time.MILLIS_PER_SECOND
	}

	//毫秒
	if length := len(s); secIndex > 0 && secIndex < length-1 {
		millisStr := s[secIndex+1:length]
		v, err := strconv.Atoi(millisStr)
		if err != nil {
			return 0, &errors.ErrorWrapper{&err, &(fmt.Sprintf("%s 时间格式解析有误 时间格式应该为 分%s秒%s毫秒x100 你填的毫秒数是 %s 不是数字", FILE_NAME, special_string.QUOTE, special_string.SINGLE_QUOTE, millisStr))}
		}
		if millis >= 10 {
			return 0, &errors.ErrorWrapper{Attach: &(fmt.Sprintf("%s 时间格式解析有误 时间格式应该为 分%s秒%s毫秒x100 你填的毫秒数不能大于9", FILE_NAME, special_string.QUOTE, special_string.SINGLE_QUOTE, millisStr))}
		}
		millis = v * 100
	}

	return minute + second + millis, nil
}

func Read() func(i, j, k int, s *xlsx.Sheet, r *xlsx.Row, c *xlsx.Cell) error {
	score := []int{}
	var sport *Sport
	return func(i, j, k int, s *xlsx.Sheet, r *xlsx.Row, c *xlsx.Cell) error {
		if j == 0 {
			rowName, err := c.String()
			if err != nil {
				return &errors.ErrorWrapper{&err, &(fmt.Sprintf("%s 读取表头的时候出错了 读到%s err=%s", FILE_NAME, rowName, err))}
			}
			name, sex, policy, err := parseRowName(rowName)
			if err != nil {
				return err
			}
			if name != SCORE_ROW_NAME {
				sport = &Sport{Name: name, Sex: sex, Policy: policy, Result: []int{}, Score: []pair.IntPair{}, UniqueKey: name + special_string.POUND_SIGN + sex}
			} else {
				sport = nil
			}
		} else {
			if sport == nil {
				cInt, err := c.Int()
				if err != nil {
					return &errors.ErrorWrapper{&err, &(fmt.Sprintf("%s 读取分值的时候出错了 err = %s", FILE_NAME, err))}
				}
				score = append(score, cInt)
			} else {
				cStr, err := c.String()
				if err != nil {
					return &errors.ErrorWrapper{&err, &(fmt.Sprintf("%s %s 读取数据的时候出错了 err = %s", FILE_NAME, sport.UniqueKey, err))}
				}

				v := sport.Policy(cStr)
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
func parseRowName(rowName string) (name string, sex string, policy UnitTypePolicy, err error) {
	ss := strings.Split(rowName, special_string.POUND_SIGN)
	name = ss[0]
	policy = defaultUnitTypePolicy
	if length := len(ss); length > 1 {
		sex = ss[1]
		if sex != "" && sex != MALE && sex != FEMAL {
			err = &errors.ErrorWrapper{&err, &(fmt.Sprintf("%s 组名的格式应该是 运动名称%s性别%s单位 而性别只能填 %s 或者 %s", FILE_NAME, special_string.POUND_SIGN, special_string.POUND_SIGN, MALE, FEMAL))}
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
		err = &errors.ErrorWrapper{Attach: &(fmt.Sprintf("%s 有一个组名是空的", FILE_NAME))}
		return
	}

	return
}
