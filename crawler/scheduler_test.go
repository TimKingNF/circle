/* crawler test scheduler */
package main

import (
	cmn "circle/common"
	anlz "circle/crawler/analyzer"
	args "circle/crawler/args"
	base "circle/crawler/base"
	ipl "circle/crawler/itempipeline"
	sched "circle/crawler/scheduler"
	"errors"
	"fmt"
	"net/http"
	"sync"
	"time"
)

func genHttpClient() *http.Client {
	cj := sched.NewCookiejar()
	return &http.Client{Jar: cj}
}

func getResponseParsers() []anlz.ParseResponse {
	parsers := []anlz.ParseResponse{
		anlz.ParseForHtml,
	}
	return parsers
}

func getItemProcessors() []ipl.ProcessItem {
	itemProcessors := []ipl.ProcessItem{
		processItem,
	}
	return itemProcessors
}

func main() {
	defer func() {
		if err := recover(); err != nil {
			fmt.Println(err)
		}
	}()

	scheduler := sched.NewScheduler()

	checkCountChan := scheduler.StartMonitoring(args.CrawlerArgs.MonitorArgs())

	httpClientGenerator := genHttpClient
	respParsers := getResponseParsers()
	itemProcessors := getItemProcessors()

	scheduler.Start(
		args.CrawlerArgs.ChannelArgs(),
		args.CrawlerArgs.SchePoolArgs(),
		args.CrawlerArgs.SpiderArgs(),
		httpClientGenerator,
		respParsers,
		itemProcessors,
	)

	startUrl := "http://news.d.cn/pc/view-38900.html"
	firstHttpReq, err := http.NewRequest("GET", startUrl, nil)
	pd, err := cmn.GetPrimaryDomain(startUrl)
	if err != nil {
		panic(err)
	}
	firstReq := base.NewRequest(firstHttpReq, pd, 0)

	scheduler.Accept(*firstReq)

	<-checkCountChan

	var waitgroup sync.WaitGroup
	waitgroup.Add(1)
	waitgroup.Wait()
}

func processItem(item base.Item) (result base.Item, err error) {
	if item == nil {
		return nil, errors.New("Invalid item!")
	}
	result = make(map[string]interface{})
	for k, v := range item {
		result[k] = v
	}
	if _, ok := result["number"]; !ok {
		result["number"] = len(result)
	}
	time.Sleep(10 * time.Millisecond)
	return result, nil
}
