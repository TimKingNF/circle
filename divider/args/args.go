package args

import (
	cmn "circle/common"
)

var (
	consoleLog       = true
	outputfileLog    = true
	outputfilePath   = ""
	outputfilePrefix = "divider"

	LoggerArgs cmn.LoggerArgs = cmn.NewLoggerArgs(
		consoleLog,
		outputfileLog,
		outputfilePath,
		outputfilePrefix)
)

const (
	DEVICE_NAME = cmn.DEVICE_DIVIDER
)
