/* crawler middleware stopsign */
package middleware

import (
	"fmt"
	"sync"
)

type StopSign interface {
	Sign() bool                   // 发出停止信号，如果先前已经发出信号，则返回false
	Signed() bool                 // 判断信号是否已经被发出
	Reset()                       // 重置停止信号。相当于收回停止信号，并清除所有的停止信号处理记录
	Deal(code string)             // 处理停止信号
	DealCount(code string) uint32 // 获取某一个停止信号处理方的处理计数
	DealTotal() uint32            // 获取停止信号被处理的总计数
	Summary() string              // 获取摘要信息，包含所有的停止信号处理记录
}

type myStopSign struct {
	signed       bool
	dealCountMap map[string]uint32
	rwmutex      sync.Mutex
}

func NewStopSign() StopSign {
	ss := &myStopSign{
		dealCountMap: make(map[string]uint32),
	}
	return ss
}

func (ss *myStopSign) Sign() bool {
	ss.rwmutex.Lock()
	defer ss.rwmutex.Unlock()
	if ss.signed {
		return false
	}
	ss.signed = true
	return true
}

func (ss *myStopSign) Signed() bool {
	return ss.signed
}

func (ss *myStopSign) Deal(code string) {
	ss.rwmutex.Lock()
	defer ss.rwmutex.Unlock()
	if !ss.signed {
		return
	}
	if _, ok := ss.dealCountMap[code]; !ok {
		ss.dealCountMap[code] = 1
	} else {
		ss.dealCountMap[code] += 1
	}
}

func (ss *myStopSign) Reset() {
	ss.rwmutex.Lock()
	defer ss.rwmutex.Unlock()
	ss.signed = false
	ss.dealCountMap = make(map[string]uint32)
}

func (ss *myStopSign) DealCount(code string) uint32 {
	ss.rwmutex.Lock()
	defer ss.rwmutex.Unlock()
	v, ok := ss.dealCountMap[code]
	if !ok {
		return 0
	}
	return v
}

func (ss *myStopSign) DealTotal() uint32 {
	ss.rwmutex.Lock()
	defer ss.rwmutex.Unlock()
	var count uint32 = 0
	for k, _ := range ss.dealCountMap {
		count += ss.dealCountMap[k]
	}
	return uint32(count)
}

func (ss *myStopSign) Summary() string {
	summary := "signed："
	if ss.Signed() {
		summary += "true"
	} else {
		summary += "false"
	}
	summary += "\n"
	for k, _ := range ss.dealCountMap {
		summary += k + "：" + fmt.Sprintf("%d", ss.dealCountMap[k]) + "\n"
	}
	return summary
}
