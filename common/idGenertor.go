/* id genertor */
package common

import (
	"math"
	"sync"
)

type IdGenertor interface {
	GenUint32() uint32

	Start(started uint32)
}

type cyclicIdGenertor struct {
	sn    uint32     // 当前的 id
	ended bool       // 前一个id 是否已经为其类型所能表示的最大值
	mutex sync.Mutex // 互斥锁
}

//	从 started 开始
func (gen *cyclicIdGenertor) Start(started uint32) {
	gen.mutex.Lock()
	defer gen.mutex.Unlock()
	gen.sn = started
}

func NewIdGenertor() IdGenertor {
	return &cyclicIdGenertor{}
}

func (gen *cyclicIdGenertor) GenUint32() uint32 {
	gen.mutex.Lock()
	defer gen.mutex.Unlock()
	if gen.ended {
		defer func() {
			gen.ended = false
		}()
		gen.sn = 0
		return gen.sn
	}
	id := gen.sn
	if id < math.MaxUint32 {
		gen.sn++
	} else {
		gen.ended = true
	}
	return id
}
