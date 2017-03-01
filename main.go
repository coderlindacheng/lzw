package main

import (
	"github.com/coderlindacheng/balabalago"
	"github.com/coderlindacheng/lzw/excel/common"
	"github.com/coderlindacheng/lzw/excel/standar"
	"log"
	"github.com/coderlindacheng/lzw/excel/source"
)

func main() {
	defer utils.Pause()

	if err := common.ReadSheet(standar.FILE_NAME, standar.Read); err != nil {
		log.Panicln(err)
	}
	if err := common.ReadSheet(source.INPUT_FILE_NAME, source.Read); err != nil {
		log.Panicln(err)
	}
	if err := source.Output(source.OUTPUT_FILE_NAME); err != nil {
		log.Panicln(err)
	}

}
