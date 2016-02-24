package indexdevice

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
			loggerArgs.OutputfilePrefix()+"_indexDeviceManager",
		))
	}
	return logger
}

const (
	DIVIDER_LISTENER_INDEXDEVICE_ADDRESS = "127.0.0.1:8087"
)

type IndexDeviceManager interface {
	//	listen the tcp socket from indexDevice
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

type myIndexDeviceManager struct {
	//	connecded indexDevice map. bool means connecting
	connMap map[string]bool

	connList *list.List

	dataChan chan *cmn.ControlMessage

	loadBalancer base.LoadBalancer
}

func NewIndexDeviceManager() IndexDeviceManager {
	connMap := make(map[string]bool)
	dataChan := make(chan *cmn.ControlMessage)
	connList := list.New()
	loadBalancer := base.NewLoadBalancer()
	return &myIndexDeviceManager{
		connMap:      connMap,
		dataChan:     dataChan,
		connList:     connList,
		loadBalancer: loadBalancer,
	}
}

func (idm *myIndexDeviceManager) Listening() {
	go func() {
		var listener net.Listener
		listener, err := net.Listen(cmn.SERVER_NETWORK, DIVIDER_LISTENER_INDEXDEVICE_ADDRESS)
		if err != nil {
			genLogger().Fatalln(err)
			return
		}
		defer listener.Close()
		genLogger().Infoln(
			fmt.Sprintf("Got IndexDevice manager listener for the divider. (local address: %s) (indexDevice manager)\n",
				listener.Addr()))
		for {
			conn, err := listener.Accept()
			if err != nil {
				genLogger().Fatalln(err)
				continue
			}
			genLogger().Infoln(
				fmt.Sprintf("Established a connection with a client application. (remote address: %s) (indexDevice manager)\n",
					conn.RemoteAddr()))
			go idm.handleConn(conn)
		}
	}()
}

func (idm *myIndexDeviceManager) DataChan() chan *cmn.ControlMessage {
	return idm.dataChan
}

func (idm *myIndexDeviceManager) ConnList() *list.List {
	return idm.connList
}

func (idm *myIndexDeviceManager) GenLoadBalancing() uint32 {
	var clusters []uint32
	for e := idm.connList.Front(); e != nil; e = e.Next() {
		v := e.Value.(base.DeviceConn)
		clusters = append(clusters, v.Id())
	}
	idm.loadBalancer.Clusters(clusters)
	return idm.loadBalancer.GenLoadBalancing()
}
