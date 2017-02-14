package main



import (
	"log"
	"os"
	"bufio"
	"fmt"
	"strings"
	"github.com/coderlindacheng/balabalago/pair"
	"github.com/tealeg/xlsx"
)


//func main() {
//
//	scoreAreaOne := []float64{100.0, 90.0}
//	timeAreaOne := []float64{3 * 60 + 21.0, 3 * 60 + 35.0}
//	scoreAreaTwo := []float64{90.0, 0.0}
//	timeAreaTwo := []float64{3 * 60 + 35.0, 5 * 60 + 5.0}
//
//	//inputTime := 3*60+50.0
//	reader := bufio.NewReader(os.Stdin)
//
//	fmt.Println("请输入时间格式为 \"分:秒\"(00:00)")
//	data, _, err := reader.ReadLine()
//	if err != nil {
//		fmt.Println("输入有误,请重新输入")
//		os.Exit(0)
//	}
//	sdata := string(data)
//	fmt.Println("你输入的是:" + sdata)
//	temp := strings.Split(sdata, ":")
//	var r1, r2 float64
//	var e1, e2 error
//
//	r1, e1 = strconv.ParseFloat(temp[0], 64)
//	if e1 != nil {
//		fmt.Println("输入有误,请重新输入")
//	}
//	r2, e2 = strconv.ParseFloat(temp[1], 64)
//	if e2 != nil {
//		fmt.Println("输入有误,请重新输入")
//	}
//	inputTime := r1 * 60 + r2
//	var maxScore float64
//	var minScored float64
//	var minTime float64
//	var maxTime float64
//	if inputTime >= timeAreaOne[0]&&inputTime <= timeAreaOne[1] {
//		//分数区间(用来计算一定区间段内每一分对应的秒数)
//		maxScore = scoreAreaOne[0]
//		minScored = scoreAreaOne[1]
//		//时间区间
//		minTime = timeAreaOne[0]
//		maxTime = timeAreaOne[1]
//		fmt.Print("这个时间对应的得分是")
//		fmt.Println(maxScore - (inputTime - minTime) / ((maxTime - minTime) / (maxScore - minScored)))
//	} else if inputTime >= timeAreaTwo[0]&&inputTime <= timeAreaTwo[1] {
//		//分数区间(用来计算一定区间段内每一分对应的秒数)
//		maxScore = scoreAreaTwo[0]
//		minScored = scoreAreaTwo[1]
//		//时间区间
//		minTime = timeAreaTwo[0]
//		maxTime = timeAreaTwo[1]
//		fmt.Print("这个时间对应的得分是 ")
//		fmt.Println(maxScore - (inputTime - minTime) / ((maxTime - minTime) / (maxScore - minScored)))
//	} else {
//		fmt.Println("你输入的数据不在可以计算的范围内")
//	}
//	fmt.Print("请按任意键退出")
//	reader.ReadLine()
//}
type sport struct {
	name  string
	sex   string
	result []int
	score  []pair.IntPair
}

func readStandard() {
	defer pause()
	fileName := "评分标准.xlsx"
	xlsx, err := xlsx.OpenFile(fileName)
	if err != nil {
		log.Panicf("读取评分标准表出错,请确认在改目录下存在%s,\n", fileName)
	}
	sheet:=xlsx.Sheets
	if len(sheet) != 1 {
		log.Panicf("评分表只能有一个sheet\n")

	}
	for _, s := range sheet{
		for _, r := range s.Rows {
			row:=row{}
			for i, c := range r.Cells {
				switch i {
				case 0:

				}
				var v float64
				minute := strings.Split(c.Value, "\"")
				mLen:=len(minute)
				if mLen >2{
					log.Panicf()
					pause()
				}
				if > 1 {
					v+=float64(minute) * 60 * 1000
				}else {

				}
			}
		}
	}
}

func pause() {
	reader := bufio.NewReader(os.Stdin)
	fmt.Println("请按任意键继续")
	reader.ReadLine();
	os.Exit(0)
}
