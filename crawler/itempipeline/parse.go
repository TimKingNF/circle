package itempipeline

import (
	cmn "circle/common"
)

func (ip *myItemPipeline) savePageAnalyze(doc string) {
	var controlMess = &cmn.ControlMessage{
		Function: cmn.CRAWLER_SAVEPAGE,
		Doc:      doc,
		Sender:   cmn.DEVICE_CRAWLER,
		Accepter: cmn.DEVICE_INDEXDEVICE,
	}
	ip.dataChan <- controlMess
}

func (ip *myItemPipeline) updatePageAnalyze(doc string) {
	var controlMess = &cmn.ControlMessage{
		Function: cmn.INDEXDEVICE_UPDATEPAGE,
		Doc:      doc,
		Sender:   cmn.DEVICE_CRAWLER,
		Accepter: cmn.DEVICE_INDEXDEVICE,
	}
	ip.updateDataChan <- controlMess
}
