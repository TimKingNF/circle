package dividerctrlm

import (
	"bytes"
	cmn "circle/common"
	args "circle/divider/args"
	base "circle/divider/base"
	cwm "circle/divider/crawler"
	dbm "circle/divider/database"
	idm "circle/divider/indexdevice"
	logging "circle/logging"
	"fmt"
	"net"
)

var logger logging.Logger

func genLogger() logging.Logger {
	if logger == nil {
		var loggerArgs cmn.LoggerArgs = args.LoggerArgs
		logger = cmn.NewLogger(cmn.NewLoggerArgs(
			loggerArgs.ConsoleLog(),
			loggerArgs.OutputfileLog(),
			loggerArgs.OutputfilePath(),
			loggerArgs.OutputfilePrefix()+"_controlModel",
		))
	}
	return logger
}

type Divider interface {
	Init()

	CrawlerData() <-chan *cmn.ControlMessage

	DatabaseData() <-chan string

	IndexdeviceData() <-chan *cmn.ControlMessage
}

type myDivider struct {
	crawlerManager    cwm.CrawlerManager
	databaseManager   dbm.DatabaseManager
	indexdeviceManage idm.IndexDeviceManager
}

func NewDivider() Divider {
	return &myDivider{}
}

func (divider *myDivider) analyzeIndexDevice(doc *cmn.ControlMessage) {
	var buffer bytes.Buffer
	docString := doc.String()
	buffer.WriteString(docString)
	buffer.WriteByte(cmn.DELIMITER)
	docBytes := buffer.Bytes()

	var conn net.Conn

	switch doc.Function {
	case cmn.CRAWLER_QUERYINDEX:
		for e := divider.indexdeviceManage.ConnList().Front(); e != nil; e = e.Next() {
			v := e.Value.(base.DeviceConn)
			conn = v.Conn()
		}
	case cmn.CRAWLER_SAVEPAGE:
		//	intro to load balancing
		lbConnId := divider.indexdeviceManage.GenLoadBalancing()
		for e := divider.indexdeviceManage.ConnList().Front(); e != nil; e = e.Next() {
			v := e.Value.(base.DeviceConn)
			if doc.AccepterId != "" {
				if doc.AccepterId == fmt.Sprintf("%d", v.Id()) {
					conn = v.Conn()
				}
			} else {
				if v.Id() == lbConnId {
					conn = v.Conn()
				}
			}
		}
	case cmn.INDEXDEVICE_UPDATEPAGE:
		for e := divider.indexdeviceManage.ConnList().Front(); e != nil; e = e.Next() {
			v := e.Value.(base.DeviceConn)
			if doc.AccepterId != "" && doc.AccepterId == fmt.Sprintf("%d", v.Id()) {
				conn = v.Conn()
			}
		}
	}

	if conn != nil {
		genLogger().Infoln(fmt.Sprintf("Send Data To IndexDevice[%s]: %s (divider)\n",
			conn.RemoteAddr().String(),
			docString))

		conn.Write(docBytes)
	}
}

func (divider *myDivider) analyzeCrawler(doc *cmn.ControlMessage) {
	var buffer bytes.Buffer
	docString := doc.String()
	buffer.WriteString(docString)
	buffer.WriteByte(cmn.DELIMITER)
	docBytes := buffer.Bytes()

	var conn net.Conn

	switch doc.Function {
	case cmn.CRAWLER_QUERYINDEX_RESULT:
		for e := divider.crawlerManager.ConnList().Front(); e != nil; e = e.Next() {
			v := e.Value.(base.DeviceConn)
			if doc.AccepterId != "" && doc.AccepterId == fmt.Sprintf("%d", v.Id()) {
				conn = v.Conn()
			}
		}
	case cmn.INDEXDEVICE_UPDATEPAGE:
		lbConnId := divider.crawlerManager.GenLoadBalancing()
		for e := divider.crawlerManager.ConnList().Front(); e != nil; e = e.Next() {
			v := e.Value.(base.DeviceConn)
			if v.Id() == lbConnId {
				conn = v.Conn()
			}
		}
	}

	if conn != nil {
		genLogger().Infoln(fmt.Sprintf("Send Data To Crawler[%s]: %s (divider)\n",
			conn.RemoteAddr().String(),
			docString))

		conn.Write(docBytes)
	}
}

func (divider *myDivider) analyzeDatabase(doc *cmn.ControlMessage) {
	var buffer bytes.Buffer
	docString := doc.String()
	buffer.WriteString(docString)
	buffer.WriteByte(cmn.DELIMITER)
	docBytes := buffer.Bytes()
	var conn net.Conn

	//	intro to load balancing
	lbConnId := divider.databaseManager.GenLoadBalancing()
	for e := divider.databaseManager.ConnList().Front(); e != nil; e = e.Next() {
		v := e.Value.(base.DeviceConn)
		if v.Id() == lbConnId {
			conn = v.Conn()
		}
	}

	if conn != nil {
		genLogger().Infoln(fmt.Sprintf("Send Data To Database[%s]: %s (divider)\n",
			conn.RemoteAddr().String(),
			docString))

		conn.Write(docBytes)
	}
}

func (divider *myDivider) activate() {
	go func() {
		for {
			select {
			case doc, ok := <-divider.CrawlerData():
				if ok {
					//	parse function
					switch doc.Function {
					case cmn.CRAWLER_QUERYINDEX:
						divider.analyzeIndexDevice(doc)
					case cmn.CRAWLER_SAVEPAGE:
						divider.analyzeIndexDevice(doc)
					case cmn.INDEXDEVICE_UPDATEPAGE:
						divider.analyzeIndexDevice(doc)
					}
				}
			// case doc, ok := <-divider.DatabaseData():

			case doc, ok := <-divider.IndexdeviceData():
				if ok {
					//	parse function
					switch doc.Function {
					case cmn.CRAWLER_QUERYINDEX_RESULT:
						divider.analyzeCrawler(doc)
					case cmn.INDEXDEVICE_SAVEPAGE:
						divider.analyzeDatabase(doc)
					case cmn.INDEXDEVICE_UPDATEPAGE:
						divider.analyzeCrawler(doc)
					}
				}

			}
		}
	}()
}

func (divider *myDivider) Init() {
	crawlerManager := cwm.NewCrawlerManager()
	crawlerManager.Listening()
	divider.crawlerManager = crawlerManager

	databaseManager := dbm.NewDatabaseManager()
	databaseManager.Listening()
	divider.databaseManager = databaseManager

	indexdeviceManager := idm.NewIndexDeviceManager()
	indexdeviceManager.Listening()
	divider.indexdeviceManage = indexdeviceManager

	divider.activate()
}

func (divider *myDivider) CrawlerData() <-chan *cmn.ControlMessage {
	return divider.crawlerManager.DataChan()
}

func (divider *myDivider) DatabaseData() <-chan string {
	return divider.databaseManager.DataChan()
}

func (divider *myDivider) IndexdeviceData() <-chan *cmn.ControlMessage {
	return divider.indexdeviceManage.DataChan()
}
