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

type UnitTypePolicy func(fileName, s string,sheetNum,rowNum,cellNum int) (int, error)

func defaultUnitTypePolicy(fileName, s string,sheetNum, rowNum, cellNum int) (int, error) {
	v, err := strconv.Atoi(s)
	if err != nil {
		log.Printf("%s 你填的是 %s 不是数字 第%v个表第%v行%v列", fileName, s,sheetNum, rowNum, cellNum)
		return 0, err
	}
	return v, nil
}

/*
	以时间为单位的解析策略,最大到分钟,最小号毫秒,把输入的时间格式解析成毫秒

	args:
	 	时间格式 , 分"秒%s毫秒x100
	return :
		int 解析后的毫秒值
 */
func timeUnitTypePolicy(fileName, s string,sheetNum,rowNum,cellNum int) (int, error) {

	//2个index到最后一定要被查出来,所以不如一开始就查出来
	minIndex := strings.Index(s, QUOTE)
	secIndex := strings.Index(s, SINGLE_QUOTE)

	//"和%s号都没有,默认解析成分钟
	if minIndex == -1 && secIndex == -1 {
		v, err := strconv.Atoi(s)
		if err != nil {
			return 0, errors.NewWrapper(err, fmt.Sprintf("%s 时间格式解析有误 时间格式应该为 分%s秒%s毫秒x100 你填的分钟数是 %s 不是数字 第%v个表第%v行%v列", fileName, QUOTE, SINGLE_QUOTE,s, sheetNum,rowNum,cellNum))
		}
		return v * MILLIS_PER_MINUTE, nil
	}

	if minIndex == 0 || secIndex == 0 {
		return 0, errors.NewOnlyStr(fmt.Sprintf("%s 时间格式解析有误 时间格式应该为 分%s秒%s毫秒x100 %s和%s 不能出现在第一位 第%v个表第%v行%v列", fileName, QUOTE, SINGLE_QUOTE, QUOTE, SINGLE_QUOTE,sheetNum,rowNum,cellNum))
	}

	var minute int
	var second int
	var millis int

	//分
	if minIndex > 0 {
		minuteStr := s[:minIndex]
		v, err := strconv.Atoi(minuteStr)
		if err != nil {
			return 0, errors.NewWrapper(err, fmt.Sprintf("%s 时间格式解析有误 时间格式应该为 分%s秒%s毫秒x100 你填的分钟数是 %s 不是数字 第%v个表第%v行%v列", fileName, QUOTE, SINGLE_QUOTE, minuteStr,sheetNum,rowNum,cellNum))
		}
		minute = v * MILLIS_PER_MINUTE
	}

	var secondStr string

	//秒
	if secIndex > 0 && minIndex > 0 {
		if secIndex-minIndex <= 1 {
			return 0, errors.NewOnlyStr(fmt.Sprintf("%s 时间格式解析有误 时间格式应该为 分%s秒%s毫秒x100 %s一定要出现在%s后面而且不能连续出现 第%v个表第%v行%v列", fileName, QUOTE, SINGLE_QUOTE, SINGLE_QUOTE, QUOTE,sheetNum,rowNum,cellNum))
		}
		secondStr = s[minIndex+1:secIndex]
	} else if secIndex > 0 && minIndex < 0 {
		secondStr = s[:secIndex]
	} else if secIndex < 0 && minIndex > 0 {
		secondStr = s[minIndex+1:]
	}

	if secondStr != "" {
		v, err := strconv.Atoi(secondStr)
		if err != nil {
			return 0, errors.NewWrapper(err, fmt.Sprintf("%s 时间格式解析有误 时间格式应该为 分%s秒%s毫秒x100 你填的秒数是 %s 不是数字 第%v个表第%v行%v列", fileName, QUOTE, SINGLE_QUOTE, secondStr,sheetNum,rowNum,cellNum))
		}
		second = v * MILLIS_PER_SECOND
	}

	//毫秒
	if length := len(s); secIndex > 0 && secIndex < length-1 {
		millisStr := s[secIndex+1:]
		v, err := strconv.Atoi(millisStr)
		if err != nil {
			return 0, errors.NewWrapper(err, fmt.Sprintf("%s 时间格式解析有误 时间格式应该为 分%s秒%s毫秒x100 你填的毫秒数是 %s 不是数字 第%v个表第%v行%v列", fileName, QUOTE, SINGLE_QUOTE, millisStr,sheetNum,rowNum,cellNum))
		}
		if millis >= 10 {
			return 0, errors.NewWrapper(err, fmt.Sprintf("%s 时间格式解析有误 时间格式应该为 分%s秒%s毫秒x100 你填的毫秒数不能大于9 第%v个表第%v行%v列", fileName, QUOTE, SINGLE_QUOTE,sheetNum,rowNum,cellNum))
		}
		millis = v * 100
	}
	if minute > 60*MILLIS_PER_MINUTE {
		return 0, errors.NewOnlyStr(fmt.Sprintf("%s 以时间为单位的数据分钟是不能大于或者等于60的 第%v个表第%v行%v列", fileName,sheetNum,rowNum,cellNum))
	}
	if second > 60*MILLIS_PER_SECOND {
		return 0, errors.NewOnlyStr(fmt.Sprintf("%s 以时间为单位的数据秒数是不能大于或者等于60的 第%v个表第%v行%v列", fileName,sheetNum,rowNum,cellNum))
	}
	if millis > MILLIS_PER_SECOND {
		return 0, errors.NewOnlyStr(fmt.Sprintf("%s 以时间为单位的数据毫秒数字是不能大于或者等于10 第%v个表第%v行%v列", fileName,sheetNum,rowNum,cellNum))
	}
	return minute + second + millis, nil
}
