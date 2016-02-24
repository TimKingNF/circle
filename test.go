package main

import (
	cmn "circle/common"
	// cwlcm "circle/crawler/crawlerctrlm"
	"fmt"
	"net/url"
	"time"
)

func main() {
	fmt.Println(cmn.CrawlerStruct())

	// crawlerCtrlM := cwlcm.GenCrawlerControlModel()
	// fmt.Println(crawlerCtrlM.OsInfo())

	// test()
	r := "http://baidu.com"
	urls, _ := url.Parse("http://baidu.com")
	fmt.Println(urls.Host)
	fmt.Println(cmn.GetPrimaryDomain(r))
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
