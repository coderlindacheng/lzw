package standar

import (
	"fmt"
	"strconv"
	"log"
	"strings"
	"github.com/coderlindacheng/balabalago/errors"
	. "github.com/coderlindacheng/balabalago/time"
	. "github.com/coderlindacheng/balabalago/special_string"
)

type UnitTypePolicy func(string) (int, error)

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
	minIndex := strings.Index(s, QUOTE)
	secIndex := strings.Index(s, SINGLE_QUOTE)

	//"和%s号都没有,默认解析成分钟
	if minIndex == -1 && secIndex == -1 {
		v, err := strconv.Atoi(s)
		if err != nil {
			return 0, errors.NewWrapper(err, fmt.Sprintf("%s 时间格式解析有误 时间格式应该为 分%s秒%s毫秒x100 你填的分钟数是 %s 不是数字", FILE_NAME, QUOTE, SINGLE_QUOTE, s))
		}
		return v * MILLIS_PER_MINUTE, nil
	}

	if minIndex == 0 || secIndex == 0 {
		return 0, errors.NewOnlyStr(fmt.Sprintf("%s 时间格式解析有误 时间格式应该为 分%s秒%s毫秒x100 %s和%s 不能出现在第一位", FILE_NAME, QUOTE, SINGLE_QUOTE, QUOTE, SINGLE_QUOTE))
	}

	var minute int
	var second int
	var millis int

	//分
	if minIndex > 0 {
		minuteStr := s[0:minIndex]
		v, err := strconv.Atoi(minuteStr)
		if err != nil {
			return 0, errors.NewWrapper(err, fmt.Sprintf("%s 时间格式解析有误 时间格式应该为 分%s秒%s毫秒x100 你填的分钟数是 %s 不是数字", FILE_NAME, QUOTE, SINGLE_QUOTE, minuteStr))
		}
		minute = v * MILLIS_PER_MINUTE
	}

	//秒
	if secIndex > 0 {
		if minIndex > 0 {
			if secIndex-minIndex <= 1 {
				return 0, errors.NewOnlyStr(fmt.Sprintf("%s 时间格式解析有误 时间格式应该为 分%s秒%s毫秒x100 %s一定要出现在%s后面而且不能连续出现", FILE_NAME, QUOTE, SINGLE_QUOTE, SINGLE_QUOTE, QUOTE))
			}
			secondStr := s[minIndex+1:secIndex]
			v, err := strconv.Atoi(secondStr)
			if err != nil {
				return 0, errors.NewWrapper(err, fmt.Sprintf("%s 时间格式解析有误 时间格式应该为 分%s秒%s毫秒x100 你填的秒数是 %s 不是数字", FILE_NAME, QUOTE, SINGLE_QUOTE, secondStr))
			}
			second = v
		} else {
			secondStr := s[0:secIndex]
			v, err := strconv.Atoi(secondStr)
			if err != nil {
				return 0, errors.NewWrapper(err, fmt.Sprintf("%s 时间格式解析有误 时间格式应该为 分%s秒%s毫秒x100 你填的秒数是 %s 不是数字", FILE_NAME, QUOTE, SINGLE_QUOTE, secondStr))
			}
			second = v
		}
		second = second * MILLIS_PER_SECOND
	}

	//毫秒
	if length := len(s); secIndex > 0 && secIndex < length-1 {
		millisStr := s[secIndex+1:length]
		v, err := strconv.Atoi(millisStr)
		if err != nil {
			return 0, errors.NewWrapper(err, fmt.Sprintf("%s 时间格式解析有误 时间格式应该为 分%s秒%s毫秒x100 你填的毫秒数是 %s 不是数字", FILE_NAME, QUOTE, SINGLE_QUOTE, millisStr))
		}
		if millis >= 10 {
			return 0, errors.NewWrapper(err, fmt.Sprintf("%s 时间格式解析有误 时间格式应该为 分%s秒%s毫秒x100 你填的毫秒数不能大于9", FILE_NAME, QUOTE, SINGLE_QUOTE, millisStr))
		}
		millis = v * 100
	}

	return minute + second + millis, nil
}
