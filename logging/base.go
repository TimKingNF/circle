/* circle logging logger */
package logging

import (
	"fmt"
	"log"
	"os"
	"runtime"
	"strings"
)

type Position uint

const (
	POSITION_SINGLE Position = 1 // 单例
)

var outputfilePath = getfilepath() + "/../logs"

func init() {
	log.SetFlags(log.LstdFlags)
}

type Logger interface {
	GetPosition() Position
	SetPosition(pos Position)
	ConsoleLog() bool
	OutputfileLog() bool
	OutputfilePath() string
	OutputfilePrefix() string
	SetConsoleLog(b bool)
	SetOutputfileLog(b bool)
	SetOutputfilePath(str string) error
	SetOutputfilePrefix(str string)

	Init() error
	File() os.File

	Error(v ...interface{}) string
	Errorf(format string, v ...interface{}) string
	Errorln(v ...interface{}) string
	Fatal(v ...interface{}) string
	Fatalf(format string, v ...interface{}) string
	Fatalln(v ...interface{}) string
	Info(v ...interface{}) string
	Infof(format string, v ...interface{}) string
	Infoln(v ...interface{}) string
	Panic(v ...interface{}) string
	Panicf(format string, v ...interface{}) string
	Panicln(v ...interface{}) string
	Warn(v ...interface{}) string
	Warnf(format string, v ...interface{}) string
	Warnln(v ...interface{}) string
}

func getInvokerLocation(skipNumber int) string {
	pc, file, line, ok := runtime.Caller(skipNumber)
	if !ok {
		return ""
	}
	simpleFileName := ""
	if index := strings.LastIndex(file, "/"); index > 0 {
		simpleFileName = file[index+1 : len(file)]
	}
	funcPath := ""
	funcPtr := runtime.FuncForPC(pc)
	if funcPtr != nil {
		funcPath = funcPtr.Name()
	}
	return fmt.Sprintf("%s : (%s:%d)", funcPath, simpleFileName, line)
}

func generateLogContent(
	logTag LogTag,
	pos Position,
	format string,
	v ...interface{}) string {
	skipNumber := int(pos) + 2
	baseInfo :=
		fmt.Sprintf("%s %s - ", logTag.Prefix(), getInvokerLocation(skipNumber))
	var result string
	if len(format) > 0 {
		result = fmt.Sprintf((baseInfo + format), v...)
	} else {
		vLen := len(v)
		params := make([]interface{}, (vLen + 1))
		params[0] = baseInfo
		for i := 1; i <= vLen; i++ {
			params[i] = v[i-1]
		}
		result = fmt.Sprint(params...)
	}
	return result
}

func NewSimpleLogger() Logger {
	logger := &ConsoleLogger{
		consoleLog:       true,
		outputfileLog:    false,
		outputfilePath:   outputfilePath,
		outputfilePrefix: "",
	}
	logger.SetPosition(POSITION_SINGLE)
	return logger
}
