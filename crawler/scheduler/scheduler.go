/* crawler scheduler */
package scheduler

import (
	cmn "circle/common"
	anlz "circle/crawler/analyzer"
	args "circle/crawler/args"
	base "circle/crawler/base"
	dm "circle/crawler/datamanager"
	dl "circle/crawler/downloader"
	ipl "circle/crawler/itempipeline"
	mdw "circle/crawler/middleware"
	rc "circle/crawler/requestcache"
	uc "circle/crawler/urlcache"
	logging "circle/logging"
	"errors"
	"fmt"
	"net/http"
	"strings"
	"sync/atomic"
	"time"
)

var (
	schedulerIdGenertor cmn.IdGenertor = cmn.NewIdGenertor()

	schedStatusMap map[uint32]string = map[uint32]string{
		SCHEDULER_UNSTART: "unstart",
		SCHEDULER_RUNNING: "running",
		SCHEDULER_STOPED:  "stoped",
	}
)

const (
	//	scheduler status
	SCHEDULER_UNSTART uint32 = 0
	SCHEDULER_RUNNING uint32 = 1
	SCHEDULER_STOPED  uint32 = 2

	//	implement code
	DOWNLOADER_CODE   = "downloader"
	ANALYZER_CODE     = "analyzer"
	ITEMPIPELINE_CODE = "item_pipeline"
	SCHEDULER_CODE    = "scheduler"
)

type GenHttpClient func() *http.Client

type Scheduler interface {
	Id() uint32

	//	create scheduler and init implements
	Start(
		channelArgs base.ChannelArgs,
		schePoolArgs base.SchePoolArgs,
		spiderArgs base.SpiderArgs,
		httpClientGenerator GenHttpClient,
		respParsers []anlz.ParseResponse,
		itemProcessors []ipl.ProcessItem,
	) (err error)

	//	accept request form request cache
	Accept(req base.Request) error

	//	stop scheduler and stop all implements
	Stop() bool

	//	jude scheduler running
	Running() bool

	//	get error chan
	//	if return nil, the chan is close or scheduler stoped
	ErrorChan() <-chan error

	//	get send message chan
	SendChan() <-chan *cmn.ControlMessage

	//	get accept message chan
	AcceptChan() chan<- *cmn.ControlMessage

	//	jude scduler idle
	Idle() bool

	//	start monitor
	//	when monitor stop, return chan (check idle times)
	StartMonitoring(args base.MonitorArgs) <-chan uint64

	Summary(prefix string) SchedSummary

	Snap(text string)
}

type myScheduler struct {
	id                  uint32
	schePoolArgs        base.SchePoolArgs
	channelArgs         base.ChannelArgs
	spiderArgs          base.SpiderArgs
	chanman             mdw.ChannelManager
	stopSign            mdw.StopSign
	dlpool              dl.DownloaderPool
	analyzerPool        anlz.AnalyzerPool
	itemPipeline        ipl.ItemPipeline
	running             uint32
	reqCache            rc.RequestCache
	httpClientGenerator GenHttpClient
	urlCache            uc.UrlCache
	logger              logging.Logger
	dataManager         dm.DataManager
	urlMap              map[string]bool

	// new interface {add rule}
	respParsers    []anlz.ParseResponse
	itemProcessors []ipl.ProcessItem
}

func NewScheduler() Scheduler {
	id := genSchedulerId()
	var loggerArgs cmn.LoggerArgs = args.CrawlerArgs.LoggerArgs()
	var logger = cmn.NewLogger(cmn.NewLoggerArgs(
		loggerArgs.ConsoleLog(),
		loggerArgs.OutputfileLog(),
		loggerArgs.OutputfilePath(),
		loggerArgs.OutputfilePrefix()+"_scheduler"))
	return &myScheduler{
		id:     id,
		logger: logger,
	}
}

func (sched *myScheduler) startDownloading() {
	go func() {
		for {
			req, ok := <-sched.getReqChan()
			if !ok {
				break
			}
			go sched.download(req)
		}
	}()
}

func (sched *myScheduler) activateAnalyzers(respParsers []anlz.ParseResponse) {
	go func() {
		for {
			resp, ok := <-sched.getRespChan()
			if !ok {
				break
			}
			go sched.analyze(respParsers, resp)
		}
	}()
}

func (sched *myScheduler) openItemPipeline() {
	go func() {
		sched.itemPipeline.SetFailFast(true)
		code := ITEMPIPELINE_CODE
		for item := range sched.getItemChan() {
			go func(item base.Item) {
				defer func() {
					if p := recover(); p != nil {
						errMsg := fmt.Sprintf("Fatal Item Processing Error: %s\n", p)
						sched.logger.Fatal(errMsg)
					}
				}()
				errs := sched.itemPipeline.Send(item)
				if errs != nil {
					for _, err := range errs {
						sched.sendError(err, code)
					}
				}
			}(item)
		}
	}()
}

func (sched *myScheduler) activateData() {
	go func() {
		for {
			select {
			case doc, ok := <-sched.urlCache.SendDataChan():
				if ok {
					//	parse function
					switch doc.Function {
					case cmn.CRAWLER_QUERYINDEX:
						sched.dataManager.SendDataChan() <- doc
					}
				}
			case doc, ok := <-sched.dataManager.AcceptDataChan():
				if ok {
					//	parse function
					switch doc.Function {
					case cmn.CRAWLER_QUERYINDEX_RESULT:
						sched.urlCache.ParseQueryResult(doc.Doc)
					case cmn.INDEXDEVICE_UPDATEPAGE:
						sched.updatePageAnalyze(doc.Doc)
						if updateMess, ok := <-sched.itemPipeline.UpdateDataChan(); ok {
							if updateMess.Function == cmn.INDEXDEVICE_UPDATEPAGE {
								updateMess.AccepterId = doc.SenderId
								sched.dataManager.SendDataChan() <- updateMess
							}
						}
					}
				}
			case doc, ok := <-sched.itemPipeline.DataChan():
				if ok {
					//	parse function
					switch doc.Function {
					case cmn.CRAWLER_SAVEPAGE:
						//	base.item will be send to divider and go save page control
						sched.dataManager.SendDataChan() <- doc
					}
				}
			}
		}
	}()
}

func (sched *myScheduler) schedule(interval time.Duration) {
	go func() {
		for {
			if sched.stopSign.Signed() {
				sched.stopSign.Deal(SCHEDULER_CODE)
				return
			}
			remainder := cap(sched.getReqChan()) - len(sched.getReqChan())
			var temp *base.Request
			for remainder > 0 {
				temp = sched.reqCache.Get()
				if temp == nil {
					break
				}
				if sched.stopSign.Signed() {
					sched.stopSign.Deal(SCHEDULER_CODE)
					return
				}
				sched.getReqChan() <- *temp
				remainder--
			}
			time.Sleep(interval)
		}
	}()
}

func (sched *myScheduler) Id() uint32 {
	return sched.id
}

func (sched *myScheduler) Start(
	channelArgs base.ChannelArgs,
	schePoolArgs base.SchePoolArgs,
	spiderArgs base.SpiderArgs,
	httpClientGenerator GenHttpClient,
	respParsers []anlz.ParseResponse,
	itemProcessors []ipl.ProcessItem,
) (err error) {
	defer func() {
		if p := recover(); p != nil {
			errMsg := fmt.Sprintf("Fatal Scheduler Error: %s\n", p)
			sched.logger.Fatal(errMsg)
			err = errors.New(errMsg)
		}
	}()
	if atomic.LoadUint32(&sched.running) == SCHEDULER_RUNNING {
		return errors.New("The scheduler has been started!\n")
	}

	//	init implements
	if err := channelArgs.Check(); err != nil {
		return err
	}
	sched.channelArgs = channelArgs

	if err := schePoolArgs.Check(); err != nil {
		return err
	}
	sched.schePoolArgs = schePoolArgs

	if err := spiderArgs.Check(); err != nil {
		return err
	}
	sched.spiderArgs = spiderArgs

	sched.chanman = generateChannelManager(sched.channelArgs)

	if httpClientGenerator == nil {
		return errors.New("The HTTP client generator list is invalid!")
	}
	sched.httpClientGenerator = httpClientGenerator

	dlpool, err := generateDownloaderPool(
		sched.schePoolArgs.DownloaderPoolSize(), httpClientGenerator)
	if err != nil {
		errMsg := fmt.Sprintf("Occur error when get page downloader pool: %s\n", err)
		return errors.New(errMsg)
	}
	sched.dlpool = dlpool

	analyzerPool, err := generateAnalyzerPool(sched.schePoolArgs.AnalyzerPoolSize())
	if err != nil {
		errMsg := fmt.Sprintf("Occur error when get analyzer pool: %s\n", err)
		return errors.New(errMsg)
	}
	sched.analyzerPool = analyzerPool

	if itemProcessors == nil {
		return errors.New("The item processor list is invalid!")
	}
	for i, ip := range itemProcessors {
		if ip == nil {
			return errors.New(fmt.Sprintf("The %dth item processor is invalid!", i))
		}
	}
	sched.itemProcessors = itemProcessors
	sched.itemPipeline = generateItemPipeline(itemProcessors)

	if sched.stopSign == nil {
		sched.stopSign = mdw.NewStopSign()
	} else {
		sched.stopSign.Reset()
	}

	sched.reqCache = rc.GenRequestCache()
	sched.urlCache = uc.GenUrlCache()
	sched.urlMap = make(map[string]bool)
	sched.dataManager = dm.GenDataManager()
	sched.respParsers = respParsers

	//	scheduler is running
	atomic.StoreUint32(&sched.running, SCHEDULER_RUNNING)

	//	scheduler executing
	sched.startDownloading()
	sched.activateAnalyzers(respParsers)
	sched.openItemPipeline()
	sched.urlCache.Run()
	sched.activateData()
	sched.schedule(10 * time.Millisecond)

	return nil
}

func (sched *myScheduler) Accept(req base.Request) error {
	for {
		if running := atomic.LoadUint32(&sched.running); running == SCHEDULER_RUNNING {
			break
		}
		time.Sleep(10 * time.Millisecond)
	}
	reqHttp := req.HttpReq()
	if reqHttp == nil {
		return errors.New("The HTTP request is invalid!")
	}
	code := generateCode(SCHEDULER_CODE, sched.id)
	sched.saveReqToCache(req, code)
	return nil
}

func (sched *myScheduler) Stop() bool {
	if atomic.LoadUint32(&sched.running) != SCHEDULER_RUNNING {
		return false
	}
	sched.stopSign.Sign()
	sched.chanman.Close()
	sched.reqCache.Close()

	atomic.StoreUint32(&sched.running, SCHEDULER_STOPED)
	return true
}

func (sched *myScheduler) Running() bool {
	return atomic.LoadUint32(&sched.running) == SCHEDULER_RUNNING
}

func (sched *myScheduler) ErrorChan() <-chan error {
	if sched.chanman.Status() != mdw.CHANNEL_MANAGER_STATUS_INITALIZED {
		return nil
	}
	return sched.getErrorChan()
}

func (sched *myScheduler) SendChan() <-chan *cmn.ControlMessage {
	if sched.chanman.Status() != mdw.CHANNEL_MANAGER_STATUS_INITALIZED {
		return nil
	}
	return sched.dataManager.SendDataChan()
}

func (sched *myScheduler) AcceptChan() chan<- *cmn.ControlMessage {
	if sched.chanman.Status() != mdw.CHANNEL_MANAGER_STATUS_INITALIZED {
		return nil
	}
	return sched.dataManager.AcceptDataChan()
}

func (sched *myScheduler) Idle() bool {
	idleDlpool := sched.dlpool.Used() == 0
	idleAnalyzerPool := sched.analyzerPool.Used() == 0
	idleItemPipeline := sched.itemPipeline.ProcessingNumber() == 0
	if idleDlpool && idleAnalyzerPool && idleItemPipeline {
		return true
	}
	return false
}

func (sched *myScheduler) Summary(prefix string) SchedSummary {
	return NewSchedSummary(sched, prefix)
}

func (sched *myScheduler) Snap(text string) {
	if !sched.Running() {
		return
	}
	sched.logger.Infoln(text)
}

func (sched *myScheduler) StartMonitoring(
	monitorArgs base.MonitorArgs,
) <-chan uint64 {
	return Monitoring(
		sched,
		monitorArgs.IntervalNs(),
		monitorArgs.MaxIdleCount(),
		monitorArgs.AutoStop(),
		monitorArgs.DetailSummary(),
		func(level byte, content string) {
			if len(content) == 0 {
				return
			}
			switch level {
			case SCHEDULER_MONITOR_RECORD_COMMON:
				sched.logger.Infoln(content)
			case SCHEDULER_MONITOR_RECORD_NOTICE:
				sched.logger.Warnln(content)
			case SCHEDULER_MONITOR_RECORD_ERROR:
				sched.logger.Infoln(content)
			}
		})
}

func (sched *myScheduler) getReqChan() chan base.Request {
	reqChan, err := sched.chanman.ReqChan()
	if err != nil {
		panic(err)
	}
	return reqChan
}

func (sched *myScheduler) getRespChan() chan base.Response {
	respChan, err := sched.chanman.RespChan()
	if err != nil {
		panic(err)
	}
	return respChan
}

func (sched *myScheduler) getErrorChan() chan error {
	errChan, err := sched.chanman.ErrorChan()
	if err != nil {
		panic(err)
	}
	return errChan
}

func (sched *myScheduler) getItemChan() chan base.Item {
	itemChan, err := sched.chanman.ItemChan()
	if err != nil {
		panic(err)
	}
	return itemChan
}

func (sched *myScheduler) sendResp(resp base.Response, code string) bool {
	if sched.stopSign.Signed() {
		sched.stopSign.Deal(code)
		return false
	}
	sched.getRespChan() <- resp
	return true
}

func (sched *myScheduler) saveReqToCache(req base.Request, code string) {
	httpReq := req.HttpReq()
	// fmt.Println(httpReq.URL.String())
	if httpReq == nil {
		sched.logger.Warnln("Ignore the request! It's HTTP request is invalid!")
		return
	}
	reqUrl := httpReq.URL
	if reqUrl == nil {
		sched.logger.Warnln("Ignore the request! It's HTTP request is invalid!")
		return
	}
	if strings.ToLower(reqUrl.Scheme) != "http" {
		sched.logger.Warnf("Ignore the request! It's url scheme '%s', but should be 'http'!\n",
			reqUrl.Scheme)
		return
	}
	downloadRet := sched.urlCache.Downloading(reqUrl.String())
	if downloadRet == uc.URLCACHE_QUERY_TIMEOUT {
		sched.logger.Warnf("Ignore the request! Download Query is timeout. (requestUrl=%s)\n", reqUrl)
		return
	}
	if downloadRet == uc.URLCACHE_QUERY_FALSE {
		// sched.logger.Warnf("Ignore the request! It's url is repeated. (requestUrl=%s)\n", reqUrl)
		return
	}
	if !sched.spiderArgs.CrossDomain() {
		if pd, _ := getPrimaryDomain(httpReq.URL.String()); pd != req.PrimaryDomain() {
			sched.logger.Warnf("Ignore the request! It's host '%s' not in primary domain '%s'."+
				"(requestUrl=%s)\n", httpReq.URL.String(), req.PrimaryDomain(), reqUrl)
			return
		}
	}
	if req.Depth() > sched.spiderArgs.CrawlDepth() {
		sched.logger.Warnf("Ignore the request! It's depth %d greater than %d. (requestUrl=%s)\n",
			req.Depth(), sched.spiderArgs.CrawlDepth(), reqUrl)
		return
	}
	if sched.stopSign.Signed() {
		sched.stopSign.Deal(code)
		return
	}
	ok := sched.reqCache.Put(&req)
	if !ok {
		return
	}
	sched.urlMap[reqUrl.String()] = true
	return
}

func (sched *myScheduler) updateReqToCache(req base.Request, code string) {
	httpReq := req.HttpReq()
	if httpReq == nil {
		sched.logger.Warnln("Ignore the update request! It's HTTP request is invalid!")
		return
	}
	reqUrl := httpReq.URL
	if reqUrl == nil {
		sched.logger.Warnln("Ignore the update request! It's HTTP request is invalid!")
		return
	}
	if strings.ToLower(reqUrl.Scheme) != "http" {
		sched.logger.Warnf("Ignore the update request! It's url scheme '%s', but should be 'http'!\n",
			reqUrl.Scheme)
		return
	}
	if req.MaxDepth() != 0 && req.Depth() == req.MaxDepth() {
		sched.logger.Warnf("Ignore the update request! It's depth %d greater than %d. (requestUrl=%s)\n",
			req.Depth(), req.MaxDepth(), reqUrl)
		return
	}
	if sched.stopSign.Signed() {
		sched.stopSign.Deal(code)
		return
	}
	ok := sched.reqCache.Put(&req)
	if !ok {
		return
	}
	return
}

func (sched *myScheduler) sendItem(item base.Item, code string) bool {
	if sched.stopSign.Signed() {
		sched.stopSign.Deal(code)
		return false
	}
	sched.getItemChan() <- item
	return true
}

func (sched *myScheduler) sendError(err error, code string) bool {
	if err != nil {
		return false
	}
	codePrefix := parseCode(code)[0]
	var errorType base.ErrorType
	switch codePrefix {
	case DOWNLOADER_CODE:
		errorType = base.DOWNLOADER_ERROR
	case ANALYZER_CODE:
		errorType = base.ANALYZER_ERROR
	case ITEMPIPELINE_CODE:
		errorType = base.ITEM_PROCESSOR_ERROR
	case SCHEDULER_CODE:
		errorType = base.SCHEDULER_ERROR
	}
	cError := base.NewCrawlerError(errorType, err.Error())
	if sched.stopSign.Signed() {
		sched.stopSign.Deal(code)
		return false
	}
	go func() {
		sched.getErrorChan() <- cError
	}()
	return true
}

func (sched *myScheduler) download(req base.Request) {
	defer func() {
		if p := recover(); p != nil {
			errMsg := fmt.Sprintf("Fatal Downloader Error: %s\n", p)
			sched.logger.Fatal(errMsg)
		}
	}()

	downloader, err := sched.dlpool.Take()
	if err != nil {
		errMsg := fmt.Sprintf("Downloader pool error: %s", err)
		sched.sendError(errors.New(errMsg), SCHEDULER_CODE)
		return
	}
	defer func() {
		err := sched.dlpool.Return(downloader)
		if err != nil {
			errMsg := fmt.Sprintf("Downloader pool error: %s", err)
			sched.sendError(errors.New(errMsg), SCHEDULER_CODE)
		}
	}()
	code := generateCode(DOWNLOADER_CODE, downloader.Id())
	resp, err := downloader.Download(req)
	if resp != nil {
		sched.sendResp(*resp, code)
	}
	if err != nil {
		sched.sendError(err, code)
	}
}

func (sched *myScheduler) analyze(
	respParsers []anlz.ParseResponse, resp base.Response) {
	defer func() {
		if p := recover(); p != nil {
			errMsg := fmt.Sprintf("Fatal Analyzer Error: %s\n", p)
			sched.logger.Fatal(errMsg)
		}
	}()
	analyzer, err := sched.analyzerPool.Take()
	if err != nil {
		errMsg := fmt.Sprintf("Analyzer pool error: %s", err)
		sched.sendError(errors.New(errMsg), SCHEDULER_CODE)
		return
	}
	defer func() {
		err := sched.analyzerPool.Return(analyzer)
		if err != nil {
			errMsg := fmt.Sprintf("Analyzer pool error: %s", err)
			sched.sendError(errors.New(errMsg), SCHEDULER_CODE)
		}
	}()
	code := generateCode(ANALYZER_CODE, analyzer.Id())
	dataList, errs := analyzer.Analyze(respParsers, resp)
	if dataList != nil {
		for _, data := range dataList {
			if data == nil {
				continue
			}
			switch d := data.(type) {
			case *base.Request:
				sched.saveReqToCache(*d, code)
			case *base.Item:
				sched.sendItem(*d, code)
			default:
				errMsg := fmt.Sprintf("Unsupported data type '%T'! (value=%v)\n", d, d)
				sched.sendError(errors.New(errMsg), code)
			}
		}
	}
	if errs != nil {
		for _, err := range errs {
			if err == nil {
				continue
			}
			errMsg := fmt.Sprintf("Analyzer processing error: %s", err)
			sched.sendError(errors.New(errMsg), code)
		}
	}
}

func genSchedulerId() uint32 {
	return schedulerIdGenertor.GenUint32()
}
