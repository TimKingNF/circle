package urlcache

import (
	"sync"
)

type StopSign interface {
	Sign() bool   // 发出停止信号，如果先前已经发出信号，则返回false
	Signed() bool // 判断信号是否已经被发出
	Reset()       // 重置停止信号。相当于收回停止信号，并清除所有的停止信号处理记录
}

type myStopSign struct {
	signed  bool
	rwmutex sync.Mutex
}

func NewStopSign() StopSign {
	return &myStopSign{}
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

func (ss *myStopSign) Reset() {
	ss.rwmutex.Lock()
	defer ss.rwmutex.Unlock()
	ss.signed = false
}
