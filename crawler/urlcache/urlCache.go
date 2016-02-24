package urlcache

import (
	cmn "circle/common"
	"encoding/json"
	"sync"
	"time"
)

type myQueryResult map[string]uint8
type UrlCacheQueryResult uint8

const (
	URLCACHE_QUERY_FALSE   UrlCacheQueryResult = 0
	URLCACHE_QUERY_TRUE    UrlCacheQueryResult = 1
	URLCACHE_QUERY_TIMEOUT UrlCacheQueryResult = 2
)

type UrlCache interface {
	//	start url cache
	Run()

	//	url can be download, this function will be wait before ParseQueryResult() run
	Downloading(url string) UrlCacheQueryResult

	//	url send to scheduler.dataManager chan
	SendDataChan() <-chan *cmn.ControlMessage

	//	query result parse, return downloading() result and clean urlMap
	ParseQueryResult(doc string)
}

type myUrlCache struct {
	//	will be send urlMap
	//	map[ md5(url) ] url
	urlMap map[string]string

	//	received result from divider
	resultMap map[string]uint8

	sendDataChan chan *cmn.ControlMessage

	stopSign StopSign
	mutex    sync.Mutex
}

const (
	templates       = "{%s}"
	paramsTemplates = "\"%s\": \"%s\"," // md5(url) : url
)

func GenUrlCache() UrlCache {
	stopSign := NewStopSign()
	return &myUrlCache{
		urlMap:       make(map[string]string),
		sendDataChan: make(chan *cmn.ControlMessage),
		resultMap:    make(map[string]uint8),
		stopSign:     stopSign,
	}
}

func (uc *myUrlCache) sendUrlMapToDivider() {
	go func() {
		uc.mutex.Lock()
		defer uc.mutex.Unlock()

		urlMapJson, _ := json.Marshal(uc.urlMap)

		var controlMess = &cmn.ControlMessage{
			Function: cmn.CRAWLER_QUERYINDEX,
			Doc:      string(urlMapJson),
			Sender:   cmn.DEVICE_CRAWLER,
			Accepter: cmn.DEVICE_INDEXDEVICE,
		}

		//	reset stop sign
		uc.stopSign.Reset()

		uc.sendDataChan <- controlMess
	}()
}

func (uc *myUrlCache) Run() {
	go func() {
		for {
			if len(uc.urlMap) > 0 && !uc.stopSign.Signed() {
				uc.sendUrlMapToDivider()
			}
			//	this time must be lt scheduler idle time
			time.Sleep(10 * time.Second)
		}
	}()
}

func (uc *myUrlCache) Downloading(url string) UrlCacheQueryResult {
	urlMd5 := uc.addUrl(url)

	var runTime = time.Now()
	for {
		if uc.stopSign.Signed() {
			if b, ok := uc.resultMap[urlMd5]; !ok {
				return URLCACHE_QUERY_FALSE
			} else {
				if b == 1 {
					return URLCACHE_QUERY_TRUE
				}
				return URLCACHE_QUERY_FALSE
			}
		}
		if time.Now().Sub(runTime) > 5*time.Minute {
			return URLCACHE_QUERY_TIMEOUT
		}
		time.Sleep(10 * time.Millisecond)
	}
	return URLCACHE_QUERY_FALSE
}

func (uc *myUrlCache) SendDataChan() <-chan *cmn.ControlMessage {
	return uc.sendDataChan
}

func (uc *myUrlCache) addUrl(url string) string {
	uc.mutex.Lock()
	defer uc.mutex.Unlock()

	urlMd5 := cmn.Str2Md5(url)
	uc.urlMap[urlMd5] = url
	return urlMd5
}

func (uc *myUrlCache) ParseQueryResult(doc string) {
	go func() {
		var dat myQueryResult
		if err := json.Unmarshal([]byte(doc), &dat); err == nil {
			uc.resultMap = dat
			uc.urlMap = make(map[string]string)
			uc.stopSign.Sign()
		}
	}()
}
