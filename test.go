package main

import (
	cmn "circle/common"
	// cwlcm "circle/crawler/crawlerctrlm"
	"fmt"
	"time"
)

func main() {
	// fmt.Println(cmn.CrawlerStruct())
	fmt.Println(cmn.CircleCondition())

	// crawlerCtrlM := cwlcm.GenCrawlerControlModel()
	// fmt.Println(crawlerCtrlM.OsInfo())

	test()
}

func test() {
	d1 := time.Now()
	args := cmn.NewLoggerArgs(false, true, "", "crawler")
	// fmt.Println(args.String())
	var logger = cmn.NewLogger(args)
	for i := 0; i < 100000; i++ {
		logger.Error("nihao")
	}
	d2 := time.Now()
	fmt.Println(d2.Sub(d1))
}
