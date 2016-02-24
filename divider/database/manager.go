package database

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
			loggerArgs.OutputfilePrefix()+"_databaseManager",
		))
	}
	return logger
}

const (
	DIVIDER_LISTENER_DATABASE_ADDRESS = "127.0.0.1:8086"
)

type DatabaseManager interface {
	//	listen the tcp socket from database
	Listening()

	//	handle the connection
	handleConn(conn net.Conn)

	//	get data chan
	DataChan() chan string

	//	get conn list
	ConnList() *list.List

	//	load balancing get id
	GenLoadBalancing() uint32
}

type myDatabaseManager struct {
	//	connecded database map. bool means connecting
	connMap map[string]bool

	connList *list.List

	dataChan chan string

	loadBalancer base.LoadBalancer
}

func NewDatabaseManager() DatabaseManager {
	connMap := make(map[string]bool)
	dataChan := make(chan string)
	connList := list.New()
	loadBalancer := base.NewLoadBalancer()
	return &myDatabaseManager{
		connMap:      connMap,
		dataChan:     dataChan,
		connList:     connList,
		loadBalancer: loadBalancer,
	}
}

func (dbm *myDatabaseManager) Listening() {
	go func() {
		var listener net.Listener
		listener, err := net.Listen(cmn.SERVER_NETWORK, DIVIDER_LISTENER_DATABASE_ADDRESS)
		if err != nil {
			genLogger().Fatalln(err)
			return
		}
		defer listener.Close()
		genLogger().Infoln(
			fmt.Sprintf("Got Database manager listener for the divider. (local address: %s) (database manager)\n",
				listener.Addr()))
		for {
			conn, err := listener.Accept()
			if err != nil {
				genLogger().Fatalln(err)
				continue
			}
			genLogger().Infoln(
				fmt.Sprintf("Established a connection with a client application. (remote address: %s) (database manager)\n",
					conn.RemoteAddr()))
			go dbm.handleConn(conn)
		}
	}()
}

func (dbm *myDatabaseManager) DataChan() chan string {
	return dbm.dataChan
}

func (dbm *myDatabaseManager) ConnList() *list.List {
	return dbm.connList
}

func (dbm *myDatabaseManager) GenLoadBalancing() uint32 {
	var clusters []uint32
	for e := dbm.connList.Front(); e != nil; e = e.Next() {
		v := e.Value.(base.DeviceConn)
		clusters = append(clusters, v.Id())
	}
	dbm.loadBalancer.Clusters(clusters)
	return dbm.loadBalancer.GenLoadBalancing()
}
