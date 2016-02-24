package common

import (
	"circle/logging"
)

func NewLogger(args LoggerArgs) logging.Logger {
	logger := logging.NewSimpleLogger()
	logger.SetConsoleLog(args.ConsoleLog())
	if outputfileLog := args.OutputfileLog(); outputfileLog {
		logger.SetOutputfileLog(outputfileLog)
		if outputfilePath := args.OutputfilePath(); len(outputfilePath) > 0 {
			logger.SetOutputfilePath(outputfilePath)
		}
		if outputfilePrefix := args.OutputfilePrefix(); len(outputfilePrefix) > 0 {
			logger.SetOutputfilePrefix(outputfilePrefix)
		}
	}
	logger.Init()
	return logger
}
