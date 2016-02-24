/* crawler control model */
package crawlerctrlm

import (
	"bytes"
	cmn "circle/common"
	anlz "circle/crawler/analyzer"
	args "circle/crawler/args"
	base "circle/crawler/base"
	ipl "circle/crawler/itempipeline"
	sched "circle/crawler/scheduler"
	socket "circle/crawler/socket"
	"errors"
	"net/http"
	"sync"
	"sync/atomic"
	"time"
)

var (
	crawlerCtrlMStatusMap map[uint32]string = map[uint32]string{
		CRAWLER_CONTROL_MODEL_UNSTART:     "unstart",
		CRAWLER_CONTROL_MODEL_INITIALIZED: "initialized",
		CRAWLER_CONTROL_MODEL_RUNNING:     "running",
		CRAWLER_CONTROL_MODEL_STOPED:      "stoped",
	}
)

const (
	CRAWLER_CONTROL_MODEL_UNSTART     uint32 = 0
	CRAWLER_CONTROL_MODEL_INITIALIZED uint32 = 1
	CRAWLER_CONTROL_MODEL_RUNNING     uint32 = 2
	CRAWLER_CONTROL_MODEL_STOPED      uint32 = 3
)

type CrawlerControlModel interface {
	//	initialize and start
	Init(
		args interface{},
	) error
	Start()

	//	gen local computer info
	OsInfo() string

	//	gen crawler som args of running
	//	eg: log, monitor, channel, pool, spider, divider addr
	CrawlerArgs() base.CrawlerArgs

	//	gen scheuler of crawler
	Scheduler() sched.Scheduler

	//	gen crawler status
	Running() bool

	//	accept a url and add url to download lists
	Accept(url string)

	//	socket between crawler and divider
	//	send string to divider after dial divider
	Send(data string)
}

type myCrawlerControlModel struct {
	osInfo      string
	scheduler   sched.Scheduler
	status      uint32
	crawlerArgs base.CrawlerArgs
	socket      socket.Socket
	waitGroup   sync.WaitGroup
}

func (crawlerCtrlM *myCrawlerControlModel) activate() {
	go func() {
		for {
			select {
			case doc, ok := <-crawlerCtrlM.scheduler.SendChan():
				if ok {
					crawlerCtrlM.Send(doc.String())
				}
			case doc, ok := <-crawlerCtrlM.socket.ParseDataChan():
				if ok {
					crawlerCtrlM.scheduler.AcceptChan() <- doc
				}
			}
		}
	}()
}

func GenCrawlerControlModel() CrawlerControlModel {
	return &myCrawlerControlModel{}
}

func (crawlerCtrlM *myCrawlerControlModel) Init(
	cargs interface{},
) error {
	//	set runtime args
	if cargs != nil {
		switch cargs.(type) {
		case base.CrawlerArgs:
			err := args.Reset(cargs.(base.CrawlerArgs))
			if err != nil {
				return err
			}
		default:
			return errors.New("CrawlerArgs Cannot recognize.")
		}
	}
	crawlerCtrlM.crawlerArgs = args.CrawlerArgs

	if crawlerCtrlM.scheduler != nil && crawlerCtrlM.scheduler.Running() {
		crawlerCtrlM.scheduler.Stop()
	}

	scheduler := sched.NewScheduler()

	httpClientGenerator := genHttpClient
	respParsers := getResponseParsers()
	itemProcessors := getItemProcessors()

	err := scheduler.Start(
		crawlerCtrlM.crawlerArgs.ChannelArgs(),
		crawlerCtrlM.crawlerArgs.SchePoolArgs(),
		crawlerCtrlM.crawlerArgs.SpiderArgs(),
		httpClientGenerator,
		respParsers,
		itemProcessors,
	)
	if err != nil {
		return err
	}
	crawlerCtrlM.scheduler = scheduler

	socket := socket.NewSocket(crawlerCtrlM.crawlerArgs.DividerAddr())
	if socket == nil {
		return errors.New("Socket Cannot recognize.")
	}
	socket.Run()
	crawlerCtrlM.socket = socket

	crawlerCtrlM.activate()

	atomic.StoreUint32(&crawlerCtrlM.status, CRAWLER_CONTROL_MODEL_INITIALIZED)
	return nil
}

func (crawlerCtrlM *myCrawlerControlModel) Start() {
	if atomic.LoadUint32(&crawlerCtrlM.status) != CRAWLER_CONTROL_MODEL_INITIALIZED {
		return
	}
	atomic.StoreUint32(&crawlerCtrlM.status, CRAWLER_CONTROL_MODEL_RUNNING)

	checkCountChan := crawlerCtrlM.scheduler.StartMonitoring(
		crawlerCtrlM.crawlerArgs.MonitorArgs(),
	)

	crawlerCtrlM.waitGroup.Add(1)
	crawlerCtrlM.running(checkCountChan)
	crawlerCtrlM.waitGroup.Wait()
}

func (crawlerCtrlM *myCrawlerControlModel) Accept(url string) {
	if !cmn.IsUrl(url) {
		return
	}

	for {
		if running := atomic.LoadUint32(&crawlerCtrlM.status); running == CRAWLER_CONTROL_MODEL_STOPED {
			//	initialize
			err := crawlerCtrlM.Init(crawlerCtrlM.crawlerArgs)
			if err != nil {
				return
			}
			crawlerCtrlM.scheduler.Snap("The scheduler restart !\n")
			//	restart
			go crawlerCtrlM.Start()
		} else if running == CRAWLER_CONTROL_MODEL_UNSTART {
			return
		} else {
			break
		}
		time.Sleep(10 * time.Millisecond)
	}

	httpReq, err := http.NewRequest("GET", url, nil)
	pd, err := cmn.GetPrimaryDomain(url)
	if err != nil {
		return
	}
	req := base.NewRequest(httpReq, pd, 0)

	crawlerCtrlM.scheduler.Accept(*req)
}

func (crawlerCtrlM *myCrawlerControlModel) Send(data string) {
	status := atomic.LoadUint32(&crawlerCtrlM.status)
	if status == CRAWLER_CONTROL_MODEL_UNSTART || status == CRAWLER_CONTROL_MODEL_STOPED {
		return
	}
	crawlerCtrlM.socket.Send(data)
}

func (crawlerCtrlM *myCrawlerControlModel) OsInfo() string {
	if len(crawlerCtrlM.osInfo) == 0 {
		crawlerCtrlM.genOsInfo()
	}
	return crawlerCtrlM.osInfo
}

func (crawlerCtrlM *myCrawlerControlModel) Running() bool {
	return atomic.LoadUint32(&crawlerCtrlM.status) == CRAWLER_CONTROL_MODEL_RUNNING
}

func (crawlerCtrlM *myCrawlerControlModel) CrawlerArgs() base.CrawlerArgs {
	return crawlerCtrlM.crawlerArgs
}

func (crawlerCtrlM *myCrawlerControlModel) Scheduler() sched.Scheduler {
	return crawlerCtrlM.scheduler
}

func (crawlerCtrlM *myCrawlerControlModel) genOsInfo() {
	var ret bytes.Buffer
	ret.WriteString("───────────────────────────── Local Computer ──────────────────────────────\n")
	ret.WriteString(cmn.GenOs())
	ret.WriteString(cmn.GenGoVersion())
	ret.WriteString(cmn.GenCpuNums())
	ret.WriteString(cmn.GenIpAdress())
	ret.WriteString(cmn.GenCpu())
	ret.WriteString(cmn.GenMem())
	ret.WriteString(cmn.GenDisk())
	crawlerCtrlM.osInfo = ret.String()
}

//	************************************************************************************
//	****************				运行函数							****************
//	************************************************************************************

func (crawlerCtrlM *myCrawlerControlModel) running(checkCountChan <-chan uint64) {
	if atomic.LoadUint32(&crawlerCtrlM.status) != CRAWLER_CONTROL_MODEL_RUNNING {
		return
	}

	go func(<-chan uint64) {
		defer crawlerCtrlM.waitGroup.Done()
		for {
			if _, ok := <-checkCountChan; ok {
				//	wait time out of idle time, so send completed
				crawlerCtrlM.socket.Close()
				atomic.StoreUint32(&crawlerCtrlM.status, CRAWLER_CONTROL_MODEL_STOPED)
				break
			}
			time.Sleep(time.Microsecond)
		}
	}(checkCountChan)
}

//	************************************************************************************
//	****************				内建函数							****************
//	************************************************************************************

func getResponseParsers() []anlz.ParseResponse {
	parsers := []anlz.ParseResponse{
		anlz.ParseForATag,    // gen a.href
		anlz.ParseForHtmlTag, // gen html text
	}
	return parsers
}

func getItemProcessors() []ipl.ProcessItem {
	itemProcessors := []ipl.ProcessItem{
		ipl.GenKeywordsFromPage,
	}
	return itemProcessors
}

//	************************************************************************************
//	****************				itempipeline处理函数				****************
//	************************************************************************************

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

//	************************************************************************************
//	****************				http客户端生成函数   				****************
//	************************************************************************************

func genHttpClient() *http.Client {
	cj := sched.NewCookiejar()
	return &http.Client{Jar: cj}
}
