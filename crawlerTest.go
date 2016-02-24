/* crawler test crawler control model */
package main

import (
	cmn "circle/common"
	base "circle/crawler/base"
	cwlcm "circle/crawler/crawlerctrlm"
	"fmt"
	"sync"
	"time"
)

var (
	channelArgs = base.NewChannelArgs(10, 10, 10, 10)

	schePoolArgs = base.NewSchePoolArgs(3, 3)

	intervalNs    = 10 * time.Millisecond    // 10 毫秒
	maxIdleCount  = uint(24 * 60 * 60 * 100) // 24小时
	autoStop      = true
	detailSummary = true
	moitorArgs    = base.NewMonitorArgs(intervalNs, maxIdleCount,
		autoStop, detailSummary)

	crawlDepth  = uint32(1) //	 crawlerDepth 0 L download page itself
	crossDomain = true
	spiderArgs  = base.NewSpiderArgs(crossDomain, crawlDepth)

	consoleLog       = true
	outputfileLog    = true
	outputfilePath   = ""
	outputfilePrefix = "crawler"
	loggerArgs       = cmn.NewLoggerArgs(consoleLog, outputfileLog,
		outputfilePath, outputfilePrefix)

	dividerAddr = "127.0.0.1:8085"

	crawlerArgs base.CrawlerArgs = base.NewCrawlerArgs(
		channelArgs,
		schePoolArgs,
		moitorArgs,
		spiderArgs,
		loggerArgs,
		dividerAddr,
	)
)

func main() {
	defer func() {
		if err := recover(); err != nil {
			fmt.Println(err)
		}
	}()
	crawlerCtrlM := cwlcm.GenCrawlerControlModel()

	err := crawlerCtrlM.Init(crawlerArgs)
	if err != nil {
		panic(fmt.Sprintf("crawler control model initialize failed, Error: %s", err))
	}

	//	crawlerCtrlM.scheduler.Snap test
	// crawlerCtrlM.Scheduler().Snap("123123")

	//	crawlerCtrlM initialized test
	/*args := crawlerCtrlM.CrawlerArgs()
	fmt.Println((&args).String())*/

	//	socket test
	/*for i := 0; i < 5; i++ {
		crawlerCtrlM.Send("CIRCLE|crawler-divider|queryIndex|i love u")
		time.Sleep(time.Second)
	}*/
	// crawlerCtrlM.Start()

	secondUrl := "http://news.d.cn/pc/view-38900.html"
	// startUrl := "http://news.d.cn/evaluation.html"

	crawlerCtrlM.Accept(startUrl)

	//	crawlerCtrlM.scheduler.start test
	//	restart crawlerCtrlM.scheduler
	// go test1(crawlerCtrlM) // the second round
	// crawlerCtrlM.Start() // the first round

	var waitgroup sync.WaitGroup
	waitgroup.Add(1)
	waitgroup.Wait()
}

func test1(crawlerCtrlM cwlcm.CrawlerControlModel) {
	time.Sleep(15000 * time.Millisecond)
	if !crawlerCtrlM.Scheduler().Running() {
		crawlerCtrlM.Accept("http://127.0.0.1:8080")
	}
}
