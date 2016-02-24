package indexdevicectrlm

import (
	"bytes"
	cmn "circle/common"
	args "circle/indexdevice/args"
	base "circle/indexdevice/base"
	cardsm "circle/indexdevice/cardsmanager"
	pc "circle/indexdevice/pagecache"
	socket "circle/indexdevice/socket"
	"errors"
	"sync/atomic"
	"time"
)

var (
	indexDeviceCtrlMStatusMap map[uint32]string = map[uint32]string{
		INDEXDEVICE_CONTROL_MODEL_UNSTART: "unstart",
		INDEXDEVICE_CONTROL_MODEL_RUNNING: "running",
		INDEXDEVICE_CONTROL_MODEL_STOPED:  "stoped",
	}
)

const (
	INDEXDEVICE_CONTROL_MODEL_UNSTART uint32 = 0
	INDEXDEVICE_CONTROL_MODEL_RUNNING uint32 = 1
	INDEXDEVICE_CONTROL_MODEL_STOPED  uint32 = 2
)

type IndexDeviceControlModel interface {
	// 	initialize
	//	accept param type: [base.IndexDeviceArgs]
	Init(
		idargs interface{},
	) error

	//	gen local computer info
	OsInfo() string

	//	get runtime index device args
	IndexDeviceArgs() base.IndexDeviceArgs

	//	get rnning
	Running() bool

	//	socket between database and divider
	//	send string to divider after dial divider
	Send(data string)

	//	Query keyword
	QueryKeyword(doc string)
}

type myIndexDeviceControlModel struct {
	osInfo          string
	status          uint32
	socket          socket.Socket
	cardsManager    cardsm.CardsManager
	indexDeviceArgs base.IndexDeviceArgs
	pageCache       pc.PageCache
	//	map[md5(url)] bool
	mapPage map[string]bool
}

func GenIndexDeviceControlModel() IndexDeviceControlModel {
	return &myIndexDeviceControlModel{}
}

func (indexDeviceCtrlM *myIndexDeviceControlModel) QueryKeyword(doc string) {
	//	...

	indexDeviceCtrlM.cardsManager.QueryKeyword(doc)
}

func (indexDeviceCtrlM *myIndexDeviceControlModel) activate() {
	go func() {
		for {
			select {
			case doc, ok := <-indexDeviceCtrlM.socket.ParseDataChan():
				if ok {
					switch doc.Function {
					case cmn.CRAWLER_QUERYINDEX:
						indexDeviceCtrlM.queryIndexAnalyze(doc)
					case cmn.CRAWLER_SAVEPAGE:
						indexDeviceCtrlM.parsePageAnalyze(doc)
					case cmn.INDEXDEVICE_UPDATEPAGE:
						indexDeviceCtrlM.parsePageAnalyze(doc)
					}
				}
			case doc, ok := <-indexDeviceCtrlM.pageCache.SendDataChan():
				if ok {
					switch doc.Function {
					case cmn.INDEXDEVICE_SAVEPAGE:
						indexDeviceCtrlM.socket.Send(doc.String())
					}
				}
			}
		}
	}()
}

func (indexDeviceCtrlM *myIndexDeviceControlModel) analyzePageUpdate() {
	go func() {
		for {
			//	test code
			time.Sleep(20 * time.Second)

			// time.Sleep(24 * time.Hour)
			for e := indexDeviceCtrlM.cardsManager.PageList().Front(); e != nil; e = e.Next() {
				v := e.Value.(base.MyPage)
				if v.Page.IsUpdate() {
					indexDeviceCtrlM.updatePageAnalyze(v.Page.Url())
				}
			}
		}
	}()
}

func (indexDeviceCtrlM *myIndexDeviceControlModel) Init(
	idargs interface{},
) error {
	//	set runtime args
	if idargs != nil {
		switch idargs.(type) {
		case base.IndexDeviceArgs:
			err := args.Reset(idargs.(base.IndexDeviceArgs))
			if err != nil {
				return err
			}
		default:
			return errors.New("Index Device Args Cannot recognize.")
		}
	}
	indexDeviceCtrlM.indexDeviceArgs = args.IndexDeviceArgs

	indexDeviceCtrlM.mapPage = make(map[string]bool)

	//	socket start
	socket := socket.NewSocket(indexDeviceCtrlM.indexDeviceArgs.DividerAddr())
	if socket == nil {
		return errors.New("Socket Cannot recognize.")
	}
	socket.Run()
	indexDeviceCtrlM.socket = socket

	//	cardsManager
	cardsManager := cardsm.GenCardsManager(args.FileCacheUrl)
	indexDeviceCtrlM.cardsManager = cardsManager

	//	database cache
	pageCache := pc.GenPageCache()
	pageCache.Run()
	indexDeviceCtrlM.pageCache = pageCache

	//	activate data
	indexDeviceCtrlM.activate()
	indexDeviceCtrlM.analyzePageUpdate()

	atomic.StoreUint32(&indexDeviceCtrlM.status, INDEXDEVICE_CONTROL_MODEL_RUNNING)
	return nil
}

func (indexDeviceCtrlM *myIndexDeviceControlModel) Send(data string) {
	if atomic.LoadUint32(&indexDeviceCtrlM.status) != INDEXDEVICE_CONTROL_MODEL_RUNNING {
		return
	}
	indexDeviceCtrlM.socket.Send(data)
}

func (indexDeviceCtrlM *myIndexDeviceControlModel) IndexDeviceArgs() base.IndexDeviceArgs {
	return indexDeviceCtrlM.indexDeviceArgs
}

func (indexDeviceCtrlM *myIndexDeviceControlModel) Running() bool {
	return atomic.LoadUint32(&indexDeviceCtrlM.status) == INDEXDEVICE_CONTROL_MODEL_RUNNING
}

func (indexDeviceCtrlM *myIndexDeviceControlModel) OsInfo() string {
	if len(indexDeviceCtrlM.osInfo) == 0 {
		indexDeviceCtrlM.genOsinfo()
	}
	return indexDeviceCtrlM.osInfo
}

func (indexDeviceCtrlM *myIndexDeviceControlModel) genOsinfo() {
	var ret bytes.Buffer
	ret.WriteString("───────────────────────────── Local Computer ──────────────────────────────\n")
	ret.WriteString(cmn.GenOs())
	ret.WriteString(cmn.GenGoVersion())
	ret.WriteString(cmn.GenCpuNums())
	ret.WriteString(cmn.GenIpAdress())
	ret.WriteString(cmn.GenCpu())
	ret.WriteString(cmn.GenMem())
	ret.WriteString(cmn.GenDisk())
	indexDeviceCtrlM.osInfo = ret.String()
}
