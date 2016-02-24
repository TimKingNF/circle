package scheduler

import (
	cmn "circle/common"
	base "circle/crawler/base"
	"net/http"
)

func (sched *myScheduler) updatePageAnalyze(url string) {
	go func() {
		httpReq, err := http.NewRequest("GET", url, nil)
		pd, err := cmn.GetPrimaryDomain(url)
		if err != nil {
			return
		}
		req := base.NewMaxRequest(httpReq, pd, 0, 1)
		code := generateCode(SCHEDULER_CODE, sched.id)
		sched.updateReqToCache(*req, code)
	}()
}
