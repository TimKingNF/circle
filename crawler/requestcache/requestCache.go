/* crawler scheduler request cache */
package requestcache

import (
	base "circle/crawler/base"
	"fmt"
	"sync"
)

type requestCacheStatus uint

var (
	statusMap = map[requestCacheStatus]string{
		REQUEST_CACHE_STATUS_RUNNING: "running",
		REQUEST_CACHE_STATUS_CLOSED:  "closed",
	}
)

const (
	REQUEST_CACHE_STATUS_RUNNING requestCacheStatus = 0
	REQUEST_CACHE_STATUS_CLOSED  requestCacheStatus = 1

	summaryTemplate = "status: %s," + "size: %d / %d"
)

type RequestCache interface {
	Put(req *base.Request) bool

	Get() *base.Request

	Capacity() int

	Length() int

	Close()

	Summary() string
}

type reqCacheBySlice struct {
	cache  []*base.Request
	mutex  sync.Mutex
	status requestCacheStatus
}

func GenRequestCache() RequestCache {
	rc := &reqCacheBySlice{
		cache: make([]*base.Request, 0),
	}
	return rc
}

func (rcache *reqCacheBySlice) Put(req *base.Request) bool {
	if req == nil {
		return false
	}
	if rcache.status == REQUEST_CACHE_STATUS_CLOSED {
		return false
	}
	rcache.mutex.Lock()
	defer rcache.mutex.Unlock()
	rcache.cache = append(rcache.cache, req)
	return true
}

func (rcache *reqCacheBySlice) Get() *base.Request {
	if rcache.Length() == 0 {
		return nil
	}
	if rcache.status == REQUEST_CACHE_STATUS_CLOSED {
		return nil
	}
	rcache.mutex.Lock()
	defer rcache.mutex.Unlock()
	req := rcache.cache[0]
	rcache.cache = rcache.cache[1:]
	return req
}

func (rcache *reqCacheBySlice) Capacity() int {
	return cap(rcache.cache)
}

func (rcache *reqCacheBySlice) Length() int {
	return len(rcache.cache)
}

func (rcache *reqCacheBySlice) Close() {
	if rcache.status == REQUEST_CACHE_STATUS_CLOSED {
		return
	}
	rcache.status = REQUEST_CACHE_STATUS_CLOSED
}

func (rcache *reqCacheBySlice) Summary() string {
	summary := fmt.Sprintf(summaryTemplate,
		statusMap[rcache.status],
		rcache.Length(),
		rcache.Capacity())
	return summary
}
