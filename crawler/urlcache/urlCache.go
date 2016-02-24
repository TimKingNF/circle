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
	URLCACHE_QUERY_NOFIND  UrlCacheQueryResult = 3
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

	resultMapLen uint8

	stopSign StopSign

	mutex sync.Mutex
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
		resultMapLen: 0,
	}
}

func (uc *myUrlCache) Run() {
	go func() {
		for {
			if uc.stopSign.Signed() {
				if uc.resultMapLen == 0 {
					uc.stopSign.Reset()
				}
			} else {
				if len(uc.urlMap) > 0 {
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

					//	this time must be lt scheduler idle time
					time.Sleep(10 * time.Second)
					continue
				}
			}
			time.Sleep(10 * time.Millisecond)
		}
	}()
}

func (uc *myUrlCache) Downloading(url string) UrlCacheQueryResult {
	urlMd5 := uc.addUrl(url)
	if len(urlMd5) == 0 {
		return URLCACHE_QUERY_NOFIND
	}
	var runTime = time.Now()
	for {
		if uc.stopSign.Signed() {
			if b, ok := uc.resultMap[urlMd5]; !ok {
				return URLCACHE_QUERY_FALSE
			} else {
				uc.resultMapLen--
				if b == 1 {
					return URLCACHE_QUERY_TRUE
				}
				return URLCACHE_QUERY_FALSE
			}
		}
		if time.Now().Sub(runTime) > 5*time.Minute {
			return URLCACHE_QUERY_TIMEOUT
		}
	}
	return URLCACHE_QUERY_FALSE
}

func (uc *myUrlCache) SendDataChan() <-chan *cmn.ControlMessage {
	return uc.sendDataChan
}

func (uc *myUrlCache) addUrl(url string) string {
	for {
		if !uc.stopSign.Signed() {
			uc.mutex.Lock()
			defer uc.mutex.Unlock()

			urlMd5 := cmn.Str2Md5(url)
			if _, ok := uc.urlMap[urlMd5]; !ok {
				uc.urlMap[urlMd5] = url
				return urlMd5
			} else {
				return ""
			}
		}
		time.Sleep(10 * time.Millisecond)
	}
}

func (uc *myUrlCache) ParseQueryResult(doc string) {
	go func() {
		var dat myQueryResult
		if err := json.Unmarshal([]byte(doc), &dat); err == nil {
			uc.resultMap = dat
			uc.resultMapLen = uint8(len(dat))
			uc.urlMap = make(map[string]string)
			uc.stopSign.Sign()
		}
	}()
}
