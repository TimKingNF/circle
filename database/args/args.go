package args

import (
	cmn "circle/common"
	base "circle/database/base"
)

var (
	connected_divider_address = "127.0.0.1:8086"

	connected_sql_address = "127.0.0.1:27017"

	connected_db_name = "circle"

	DatabaseArgs base.DatabaseArgs = base.NewDatabaseArgs(
		connected_divider_address,
		connected_sql_address,
		connected_db_name)

	consoleLog       = true
	outputfileLog    = true
	outputfilePath   = ""
	outputfilePrefix = "db"

	LoggerArgs cmn.LoggerArgs = cmn.NewLoggerArgs(
		consoleLog,
		outputfileLog,
		outputfilePath,
		outputfilePrefix)
)

const (
	DEVICE_NAME = "database"
)

func Reset(args base.DatabaseArgs) error {
	if err := args.Check(); err != nil {
		return err
	}
	DatabaseArgs = args
	return nil
}
