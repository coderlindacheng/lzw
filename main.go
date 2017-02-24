package main

import (
	"log"
	"os"
	"bufio"
	"fmt"
	"strings"
	"github.com/coderlindacheng/balabalago/pair"
	"github.com/tealeg/xlsx"
	"math"
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
//	}yi g
//	fmt.Print("请按任意键退出")
//	reader.ReadLine()
//}
type A struct {
	i int
	j int
	k int
}

func main() {
	var a A = A{}
	var b A = A{}
	passingValue(a) //这个执行效率是不是比下面的高?
	passingPointer(&b)
}

func passingValue(i A)  {
	//do something
}

func passingPointer(i *A)  {
	//do something
}
