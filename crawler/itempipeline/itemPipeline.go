/* crawler scheduler itemPipeline */
package itempipeline

import (
	cmn "circle/common"
	args "circle/crawler/args"
	base "circle/crawler/base"
	logging "circle/logging"
	"encoding/json"
	"errors"
	"fmt"
	"sync/atomic"
)

const summaryTemplate = "failFast: %v, processorNumber: %d," +
	"sent: %d, accepted: %d, processed: %d, processingNumber: %d"

var (
	logger logging.Logger
)

func genLogger() logging.Logger {
	if logger == nil {
		var loggerArgs cmn.LoggerArgs = args.CrawlerArgs.LoggerArgs()

		logger = cmn.NewLogger(cmn.NewLoggerArgs(
			loggerArgs.ConsoleLog(),
			loggerArgs.OutputfileLog(),
			loggerArgs.OutputfilePath(),
			loggerArgs.OutputfilePrefix()+"_itempipeline"))
	}
	return logger
}

type ItemPipeline interface {
	Send(item base.Item) []error // 发送条目

	FailFast() bool // 判断是否有某一步骤失败

	SetFailFast(failFast bool) // 设置是否失败

	Count() []uint64 // 获得已发送、已接受和已处理的条目的计数值

	ProcessingNumber() uint64 // 获取正在被处理的条目的数量

	DataChan() chan *cmn.ControlMessage // item 处理后的数据的chan 通道

	UpdateDataChan() chan *cmn.ControlMessage // 更新操作 之后 item 处理数据的 chan 通道

	Summary() string // 获取该条目处理管道的摘要信息、介绍内容
}

type myItemPipeline struct {
	itemProcessors   []ProcessItem
	failFast         bool   // fast fail sign
	sent             uint64 // sent item count
	accepted         uint64 // accepted item count
	processed        uint64 // process item count
	processingNumber uint64 // processing item count
	dataChan         chan *cmn.ControlMessage
	updateDataChan   chan *cmn.ControlMessage
}

func NewItemPipeline(itemProcessors []ProcessItem) ItemPipeline {
	if itemProcessors == nil {
		panic(errors.New(fmt.Sprintln("Invalid item processor list")))
	}
	innerItemProcessors := make([]ProcessItem, 0)
	for i, ip := range itemProcessors {
		if ip == nil {
			panic(errors.New(fmt.Sprintf("Invalid item processtor[%d]!\n", i)))
		}
		innerItemProcessors = append(innerItemProcessors, ip)
	}
	dataChan := make(chan *cmn.ControlMessage)
	updateDataChan := make(chan *cmn.ControlMessage)
	return &myItemPipeline{
		itemProcessors: itemProcessors,
		dataChan:       dataChan,
		updateDataChan: updateDataChan,
	}
}

func (ip *myItemPipeline) Send(item base.Item) []error {
	//	原子操作保证并发调用的时候 processingNumber的安全
	atomic.AddUint64(&ip.processingNumber, 1)
	//	n= -1
	//	atomic.AddUint64(number, ^uint64(-n-1))
	defer atomic.AddUint64(&ip.processingNumber, ^uint64(0))
	atomic.AddUint64(&ip.sent, 1)
	errs := make([]error, 0)
	if item == nil {
		errs = append(errs, errors.New("The item is invalid!"))
		return errs
	}
	atomic.AddUint64(&ip.accepted, 1)
	var currentItem base.Item = item
	var isfail bool
	for _, itemProcessor := range ip.itemProcessors {
		processItem, err := itemProcessor(currentItem)
		if err != nil {
			errs = append(errs, err)
			if ip.failFast {
				isfail = true
				break
			}
			genLogger().Infoln("processed Failed [item]:\n", err, "\n prcessItemFunc:\n", itemProcessor)
		}
		if processItem != nil {
			currentItem = processItem
		}
	}
	if !isfail {
		currentItemJson, err := json.Marshal(currentItem)
		if err != nil {
			errs = append(errs, err)
		}
		genLogger().Infoln("processed[item]:\n", string(currentItemJson))
		fmt.Println("processed page: ", currentItem["title"])

		ip.parseData(currentItem, string(currentItemJson))
	}

	atomic.AddUint64(&ip.processed, 1)
	return errs
}

func (ip *myItemPipeline) FailFast() bool {
	return ip.failFast
}

func (ip *myItemPipeline) SetFailFast(failFast bool) {
	ip.failFast = failFast
}

func (ip *myItemPipeline) Count() []uint64 {
	counts := make([]uint64, 3)
	counts[0] = atomic.LoadUint64(&ip.sent)
	counts[1] = atomic.LoadUint64(&ip.accepted)
	counts[2] = atomic.LoadUint64(&ip.processed)
	return counts
}

func (ip *myItemPipeline) ProcessingNumber() uint64 {
	return atomic.LoadUint64(&ip.processingNumber)
}

func (ip *myItemPipeline) Summary() string {
	counts := ip.Count()
	summary := fmt.Sprintf(summaryTemplate,
		ip.failFast, len(ip.itemProcessors),
		counts[0], counts[1], counts[2],
		ip.ProcessingNumber())
	return summary
}

func (ip *myItemPipeline) DataChan() chan *cmn.ControlMessage {
	return ip.dataChan
}

func (ip *myItemPipeline) UpdateDataChan() chan *cmn.ControlMessage {
	return ip.updateDataChan
}

func (ip *myItemPipeline) parseData(item base.Item, doc string) {
	if _, ok := item["tag"]; ok {
		switch item["tag"].(string) {
		case "html":
			if _, ok := item["update"]; ok {
				ip.updatePageAnalyze(doc)
			} else {
				ip.savePageAnalyze(doc)
			}
		}
	}
}
