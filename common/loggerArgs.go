package common

import (
	"errors"
	"fmt"
	"os"
)

const (
	/*日志参数容器的描述模板*/
	loggerArgsTemplate string = "{ consoleLog: %v, outputfileLog: %v," +
		" outputfilePath: %s, outputfilePrefix: %s }"
)

type LoggerArgs struct {
	consoleLog       bool
	outputfileLog    bool
	outputfilePath   string
	outputfilePrefix string // 前缀
	description      string
}

func NewLoggerArgs(
	consoleLog,
	outputfileLog bool,
	outputfilePath,
	outputfilePrefix string) LoggerArgs {
	return LoggerArgs{
		consoleLog:       consoleLog,
		outputfileLog:    outputfileLog,
		outputfilePath:   outputfilePath,
		outputfilePrefix: outputfilePrefix,
	}
}

func (args *LoggerArgs) Check() error {
	if len(args.outputfilePath) > 0 {
		if !pathIsExist(args.outputfilePath) {
			return errors.New("The logger output file path is error!\n")
		}
	}
	return nil
}

func (args *LoggerArgs) String() string {
	if len(args.description) == 0 {
		args.genDescription()
	}
	return args.description
}

func (args *LoggerArgs) ConsoleLog() bool {
	return args.consoleLog
}

func (args *LoggerArgs) OutputfileLog() bool {
	return args.outputfileLog
}

func (args *LoggerArgs) OutputfilePath() string {
	return args.outputfilePath
}

func (args *LoggerArgs) OutputfilePrefix() string {
	return args.outputfilePrefix
}

func (args *LoggerArgs) genDescription() {
	outputfilePath := args.outputfilePath
	prefix := args.outputfilePrefix
	if len(args.outputfilePath) == 0 {
		outputfilePath = "<null>"
	}
	if len(prefix) == 0 {
		prefix = "<null>"
	}
	args.description = fmt.Sprintf(loggerArgsTemplate,
		args.consoleLog,
		args.outputfileLog,
		outputfilePath,
		prefix)
}

func pathIsExist(path string) bool {
	_, err := os.Stat(path)
	return err == nil || os.IsExist(err)
}
