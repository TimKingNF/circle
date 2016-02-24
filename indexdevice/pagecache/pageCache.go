package pagecache

import (
	cmn "circle/common"
	base "circle/indexdevice/base"
	"encoding/json"
	"sync"
	"time"
)

type PageCache interface {
	//	start database cache run
	Run()

	//	cache page
	CachePage(page base.SavePageItem)

	//	page data send to indexDeviceCtrlM chan
	SendDataChan() <-chan *cmn.ControlMessage
}

type myPageCache struct {
	pageMap      []cmn.Page
	sendDataChan chan *cmn.ControlMessage
	mutex        sync.Mutex
}

func GenPageCache() PageCache {
	return &myPageCache{
		pageMap:      make([]cmn.Page, 0),
		sendDataChan: make(chan *cmn.ControlMessage),
	}
}

func (pc *myPageCache) Run() {
	go func() {
		for {
			if len(pc.pageMap) > 0 {
				pc.sendPageCacheToDatabase()
			}
			//	this time must be lt scheduler idle time
			time.Sleep(30 * time.Minute)

			//	test code
			// time.Sleep(10 * time.Second)
		}
	}()
}

func (pc *myPageCache) CachePage(item base.SavePageItem) {
	pc.mutex.Lock()
	defer pc.mutex.Unlock()

	page := cmn.GenPage(item.ParentUrl)
	page.SetKeywords(item.Keywords)
	page.SetTitle(item.Title)
	page.SetDescription(item.Description)
	page.Snap(item.Html)
	page.SetTopic(item.Topic)
	pc.pageMap = append(pc.pageMap, page)
}

func (pc *myPageCache) SendDataChan() <-chan *cmn.ControlMessage {
	return pc.sendDataChan
}

func (pc *myPageCache) sendPageCacheToDatabase() {
	go func() {
		pc.mutex.Lock()
		defer pc.mutex.Unlock()

		ret, _ := json.Marshal(pc.pageMap)
		//	send page to divider to database
		var controlMessage = &cmn.ControlMessage{
			Function: cmn.INDEXDEVICE_SAVEPAGE,
			Doc:      string(ret),
			Sender:   cmn.DEVICE_INDEXDEVICE,
			Accepter: cmn.DEVICE_DATABASE,
		}
		pc.pageMap = make([]cmn.Page, 0)
		pc.sendDataChan <- controlMessage
	}()
}
