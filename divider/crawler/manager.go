package crawler

import (
	cmn "circle/common"
	args "circle/divider/args"
	base "circle/divider/base"
	logging "circle/logging"
	"container/list"
	"fmt"
	"net"
)

var (
	connIdGenertor cmn.IdGenertor = cmn.NewIdGenertor()

	logger logging.Logger
)

func genLogger() logging.Logger {
	if logger == nil {
		var loggerArgs cmn.LoggerArgs = args.LoggerArgs

		logger = cmn.NewLogger(cmn.NewLoggerArgs(
			loggerArgs.ConsoleLog(),
			loggerArgs.OutputfileLog(),
			loggerArgs.OutputfilePath(),
			loggerArgs.OutputfilePrefix()+"_crawlerManager",
		))
	}
	return logger
}

const (
	DIVIDER_LISTENER_CRAWLER_ADDRESS = "127.0.0.1:8085"
)

type CrawlerManager interface {
	//	listen the tcp socket from crawler
	Listening()

	//	handle the connection
	handleConn(conn net.Conn)

	//	get data chan
	DataChan() chan *cmn.ControlMessage

	//	get conn list
	ConnList() *list.List

	//	load balancing get id
	GenLoadBalancing() uint32
}

type myCrawlerManager struct {
	//	connecded crawler map. bool means connecting
	connMap map[string]bool

	connList *list.List

	dataChan chan *cmn.ControlMessage

	loadBalancer base.LoadBalancer
}

func NewCrawlerManager() CrawlerManager {
	connMap := make(map[string]bool)
	connList := list.New()
	dataChan := make(chan *cmn.ControlMessage)
	loadBalancer := base.NewLoadBalancer()
	return &myCrawlerManager{
		connMap:      connMap,
		dataChan:     dataChan,
		connList:     connList,
		loadBalancer: loadBalancer,
	}
}

func (cm *myCrawlerManager) Listening() {
	go func() {
		var listener net.Listener
		listener, err := net.Listen(cmn.SERVER_NETWORK, DIVIDER_LISTENER_CRAWLER_ADDRESS)
		if err != nil {
			genLogger().Fatalln(err)
			return
		}
		defer listener.Close()
		genLogger().Infoln(
			fmt.Sprintf("Got Crawerler manager listener for the divider. (local address: %s) (crawler manager)\n",
				listener.Addr()))
		for {
			conn, err := listener.Accept()
			if err != nil {
				genLogger().Fatalln(err)
				continue
			}
			genLogger().Infoln(
				fmt.Sprintf("Established a connection with a client application. (remote address: %s) (crawler manager)\n",
					conn.RemoteAddr()))
			go cm.handleConn(conn)
		}
	}()
}

func (cm *myCrawlerManager) DataChan() chan *cmn.ControlMessage {
	return cm.dataChan
}

func (cm *myCrawlerManager) ConnList() *list.List {
	return cm.connList
}

func (cm *myCrawlerManager) GenLoadBalancing() uint32 {
	var clusters []uint32
	for e := cm.connList.Front(); e != nil; e = e.Next() {
		v := e.Value.(base.DeviceConn)
		clusters = append(clusters, v.Id())
	}
	cm.loadBalancer.Clusters(clusters)
	return cm.loadBalancer.GenLoadBalancing()
}
