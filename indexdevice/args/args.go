package args

import (
	cmn "circle/common"
	base "circle/indexdevice/base"
	"runtime"
	"strings"
)

var (
	connected_divider_address = "127.0.0.1:8087"

	IndexDeviceArgs base.IndexDeviceArgs = base.NewIndexDeviceArgs(
		connected_divider_address)

	consoleLog       = true
	outputfileLog    = true
	outputfilePath   = ""
	outputfilePrefix = "indexDevice"

	LoggerArgs cmn.LoggerArgs = cmn.NewLoggerArgs(
		consoleLog,
		outputfileLog,
		outputfilePath,
		outputfilePrefix)

	FileCacheUrl = getfilepath() + "/../../cache"
)

const (
	DEVICE_NAME = cmn.DEVICE_INDEXDEVICE
)

func Reset(args base.IndexDeviceArgs) error {
	if err := args.Check(); err != nil {
		return err
	}
	IndexDeviceArgs = args
	return nil
}

func getfilepath() (filepath string) {
	_, file, _, ok := runtime.Caller(0)
	if ok {
		if index := strings.LastIndex(file, "/"); index > 0 {
			filepath = file[:index]
		}
	}
	return
}
