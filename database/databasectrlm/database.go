package databasectrlm

import (
	"bytes"
	cmn "circle/common"
	api "circle/database/api"
	args "circle/database/args"
	base "circle/database/base"
	socket "circle/database/socket"
	"errors"
	"sync/atomic"
)

var (
	databaseCtrlMStatusMap map[uint32]string = map[uint32]string{
		DATABASE_CONTROL_MODEL_UNSTART: "unstart",
		DATABASE_CONTROL_MODEL_RUNNING: "running",
		DATABASE_CONTROL_MODEL_STOPED:  "stoped",
	}
)

const (
	DATABASE_CONTROL_MODEL_UNSTART uint32 = 0
	DATABASE_CONTROL_MODEL_RUNNING uint32 = 1
	DATABASE_CONTROL_MODEL_STOPED  uint32 = 2
)

type DatabaseControlModel interface {
	//	initialize
	//	accept param type: [base.DatabaseArgs]
	Init(
		dbargs interface{},
	) error

	//	gen local computer info
	OsInfo() string

	//	socket between database and divider
	//	send string to divider after dial divider
	Send(data string)

	//	get runtime db args
	DbArgs() base.DatabaseArgs

	//	get running
	Running() bool
}

type myDatabaseControlModel struct {
	osInfo    string
	socket    socket.Socket
	status    uint32
	dbArgs    base.DatabaseArgs
	apiLayout api.ApiLayout
}

func GenDatabaseControlModel() DatabaseControlModel {
	return &myDatabaseControlModel{}
}

func (dbCtrlM *myDatabaseControlModel) activate() {
	go func() {
		for {
			select {
			case doc, ok := <-dbCtrlM.socket.ParseDataChan():
				if ok {
					switch doc.Function {
					case cmn.INDEXDEVICE_SAVEPAGE:
						dbCtrlM.savePageAnalyze(doc)
					}
				}
			}
		}
	}()
}

func (dbCtrlM *myDatabaseControlModel) Init(
	dbargs interface{},
) error {
	//	set runtime args
	if dbargs != nil {
		switch dbargs.(type) {
		case base.DatabaseArgs:
			err := args.Reset(dbargs.(base.DatabaseArgs))
			if err != nil {
				return err
			}
		default:
			return errors.New("Dbargs Cannot recognize.")
		}
	}
	dbCtrlM.dbArgs = args.DatabaseArgs

	//	api layout
	apiLayout := api.NewApiLayout(args.DatabaseArgs.SqlAddr(), args.DatabaseArgs.Dbname())
	dbCtrlM.apiLayout = apiLayout

	//	connect divider
	socket := socket.NewSocket(dbCtrlM.dbArgs.DividerAddr())
	if socket == nil {
		return errors.New("Socket Cannot recognize.")
	}
	socket.Run()
	dbCtrlM.socket = socket

	//	activate data
	dbCtrlM.activate()

	atomic.StoreUint32(&dbCtrlM.status, DATABASE_CONTROL_MODEL_RUNNING)
	return nil
}

func (dbCtrlM *myDatabaseControlModel) Send(data string) {
	if atomic.LoadUint32(&dbCtrlM.status) != DATABASE_CONTROL_MODEL_RUNNING {
		return
	}
	dbCtrlM.socket.Send(data)
}

func (dbCtrlM *myDatabaseControlModel) DbArgs() base.DatabaseArgs {
	return dbCtrlM.dbArgs
}

func (dbCtrlM *myDatabaseControlModel) Running() bool {
	return atomic.LoadUint32(&dbCtrlM.status) == DATABASE_CONTROL_MODEL_RUNNING
}

func (dbCtrlM *myDatabaseControlModel) OsInfo() string {
	if len(dbCtrlM.osInfo) == 0 {
		dbCtrlM.genOsinfo()
	}
	return dbCtrlM.osInfo
}

func (dbCtrlM *myDatabaseControlModel) genOsinfo() {
	var ret bytes.Buffer
	ret.WriteString("───────────────────────────── Local Computer ──────────────────────────────\n")
	ret.WriteString(cmn.GenOs())
	ret.WriteString(cmn.GenGoVersion())
	ret.WriteString(cmn.GenCpuNums())
	ret.WriteString(cmn.GenIpAdress())
	ret.WriteString(cmn.GenCpu())
	ret.WriteString(cmn.GenMem())
	ret.WriteString(cmn.GenDisk())
	dbCtrlM.osInfo = ret.String()
}
