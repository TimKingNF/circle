/* circle logging console */
package logging

import (
	"bytes"
	"errors"
	"fmt"
	"log"
	"os"
	"sync"
	"time"
)

type ConsoleLogger struct {
	position         Position
	consoleLog       bool
	outputfileLog    bool
	outputfilePath   string
	outputfilePrefix string
	mutex            sync.Mutex
	filePtr          *os.File
	logsdate         time.Time
}

func (logger *ConsoleLogger) Init() error {
	if logger.outputfileLog {
		var sincetime = time.Now()
		var prefix bytes.Buffer
		prefix.WriteString(logger.outputfilePath)
		prefix.WriteString("/")
		prefix.WriteString(sincetime.String()[:4])
		prefix.WriteString("_")
		prefix.WriteString(sincetime.String()[5:7])
		prefix.WriteString("_")
		prefix.WriteString(sincetime.String()[8:10])
		var outputfilePrefix = prefix.String()
		if pathIsExist(outputfilePrefix) != true {
			if err := os.MkdirAll(outputfilePrefix, 0777); err != nil {
				if os.IsPermission(err) {
					return errors.New("SaveFile: not power")
				}
				return err
			}
		}
		var buffer bytes.Buffer
		buffer.WriteString(outputfilePrefix)
		buffer.WriteString("/")
		buffer.WriteString(logger.outputfilePrefix)
		buffer.WriteString(".log")
		filePath := buffer.String()

		f, err := os.OpenFile(filePath, os.O_APPEND|os.O_CREATE, 0666)
		//	linux
		//	f, err := os.OpenFile(filePath, os.O_APPEND|os.O_CREATE|os.O_RDWR, 0666)
		if err != nil {
			return err
		}
		logger.filePtr = f
		logger.logsdate = sincetime
	}
	return nil
}

func (logger *ConsoleLogger) File() os.File {
	return *logger.filePtr
}

func (logger *ConsoleLogger) GetPosition() Position {
	return logger.position
}

func (logger *ConsoleLogger) SetPosition(pos Position) {
	logger.position = pos
}

func (logger *ConsoleLogger) ConsoleLog() bool {
	return logger.consoleLog
}

func (logger *ConsoleLogger) OutputfileLog() bool {
	return logger.outputfileLog
}

func (logger *ConsoleLogger) OutputfilePath() string {
	return logger.outputfilePath
}

func (logger *ConsoleLogger) OutputfilePrefix() string {
	return logger.outputfilePrefix
}

func (logger *ConsoleLogger) SetConsoleLog(b bool) {
	logger.consoleLog = b
}

func (logger *ConsoleLogger) SetOutputfileLog(b bool) {
	logger.outputfileLog = b
}

func (logger *ConsoleLogger) SetOutputfilePath(path string) error {
	if pathIsExist(path) != true {
		if err := os.MkdirAll(path, 0777); err != nil {
			if os.IsPermission(err) {
				return errors.New("SaveFile: not power")
			}
			return err
		}
	}
	logger.outputfilePath = path
	return nil
}

func (logger *ConsoleLogger) SetOutputfilePrefix(str string) {
	logger.outputfilePrefix = str
}

func (logger *ConsoleLogger) Check() {
	logger.mutex.Lock()
	defer logger.mutex.Unlock()

	if !compareTime(logger.logsdate, time.Now()) {
		logger.filePtr.Close()
		logger.Init()
	}
}

func (logger *ConsoleLogger) Error(v ...interface{}) string {

	logger.Check()
	content := generateLogContent(getErrorLogTag(), logger.GetPosition(), "", v...)
	if logger.consoleLog {
		log.Print(content)
	}
	if logger.outputfileLog {
		sincetime := time.Now().String()
		logger.writeLogFile([]byte(fmt.Sprintf("%s %s", sincetime[:19], content)))
	}
	return content
}

func (logger *ConsoleLogger) Errorf(format string, v ...interface{}) string {

	logger.Check()
	content := generateLogContent(getErrorLogTag(), logger.GetPosition(), format, v...)
	if logger.consoleLog {
		log.Print(content)
	}
	if logger.outputfileLog {
		sincetime := time.Now().String()
		logger.writeLogFile([]byte(fmt.Sprintf("%s %s", sincetime[:19], content)))
	}
	return content
}

func (logger *ConsoleLogger) Errorln(v ...interface{}) string {

	logger.Check()
	content := generateLogContent(getErrorLogTag(), logger.GetPosition(), "", v...)
	if logger.consoleLog {
		log.Println(content)
	}
	if logger.outputfileLog {
		sincetime := time.Now().String()
		logger.writeLogFile([]byte(fmt.Sprintf("%s %s", sincetime[:19], content)))
	}
	return content
}

func (logger *ConsoleLogger) Fatal(v ...interface{}) string {

	logger.Check()
	content := generateLogContent(getFatalLogTag(), logger.GetPosition(), "", v...)
	if logger.consoleLog {
		log.Print(content)
	}
	if logger.outputfileLog {
		sincetime := time.Now().String()
		logger.writeLogFile([]byte(fmt.Sprintf("%s %s", sincetime[:19], content)))
	}
	return content
}

func (logger *ConsoleLogger) Fatalf(format string, v ...interface{}) string {

	logger.Check()
	content := generateLogContent(getFatalLogTag(), logger.GetPosition(), format, v...)
	if logger.consoleLog {
		log.Print(content)
	}
	if logger.outputfileLog {
		sincetime := time.Now().String()
		logger.writeLogFile([]byte(fmt.Sprintf("%s %s", sincetime[:19], content)))
	}
	return content
}

func (logger *ConsoleLogger) Fatalln(v ...interface{}) string {

	logger.Check()
	content := generateLogContent(getFatalLogTag(), logger.GetPosition(), "", v...)
	if logger.consoleLog {
		log.Println(content)
	}
	if logger.outputfileLog {
		sincetime := time.Now().String()
		logger.writeLogFile([]byte(fmt.Sprintf("%s %s", sincetime[:19], content)))
	}
	return content
}

func (logger *ConsoleLogger) Info(v ...interface{}) string {

	logger.Check()
	content := generateLogContent(getInfoLogTag(), logger.GetPosition(), "", v...)
	if logger.consoleLog {
		log.Print(content)
	}
	if logger.outputfileLog {
		sincetime := time.Now().String()
		logger.writeLogFile([]byte(fmt.Sprintf("%s %s", sincetime[:19], content)))
	}
	return content
}

func (logger *ConsoleLogger) Infof(format string, v ...interface{}) string {

	logger.Check()
	content := generateLogContent(getInfoLogTag(), logger.GetPosition(), format, v...)
	if logger.consoleLog {
		log.Print(content)
	}
	if logger.outputfileLog {
		sincetime := time.Now().String()
		logger.writeLogFile([]byte(fmt.Sprintf("%s %s", sincetime[:19], content)))
	}
	return content
}

func (logger *ConsoleLogger) Infoln(v ...interface{}) string {

	logger.Check()
	content := generateLogContent(getInfoLogTag(), logger.GetPosition(), "", v...)
	if logger.consoleLog {
		log.Println(content)
	}
	if logger.outputfileLog {
		sincetime := time.Now().String()
		logger.writeLogFile([]byte(fmt.Sprintf("%s %s", sincetime[:19], content)))
	}
	return content
}

func (logger *ConsoleLogger) Panic(v ...interface{}) string {

	logger.Check()
	content := generateLogContent(getPanicLogTag(), logger.GetPosition(), "", v...)
	if logger.consoleLog {
		log.Print(content)
	}
	if logger.outputfileLog {
		sincetime := time.Now().String()
		logger.writeLogFile([]byte(fmt.Sprintf("%s %s", sincetime[:19], content)))
	}
	return content
}

func (logger *ConsoleLogger) Panicf(format string, v ...interface{}) string {

	logger.Check()
	content := generateLogContent(getPanicLogTag(), logger.GetPosition(), format, v...)
	if logger.consoleLog {
		log.Print(content)
	}
	if logger.outputfileLog {
		sincetime := time.Now().String()
		logger.writeLogFile([]byte(fmt.Sprintf("%s %s", sincetime[:19], content)))
	}
	return content
}

func (logger *ConsoleLogger) Panicln(v ...interface{}) string {

	logger.Check()
	content := generateLogContent(getPanicLogTag(), logger.GetPosition(), "", v...)
	if logger.consoleLog {
		log.Println(content)
	}
	if logger.outputfileLog {
		sincetime := time.Now().String()
		logger.writeLogFile([]byte(fmt.Sprintf("%s %s", sincetime[:19], content)))
	}
	return content
}

func (logger *ConsoleLogger) Warn(v ...interface{}) string {

	logger.Check()
	content := generateLogContent(getWarnLogTag(), logger.GetPosition(), "", v...)
	if logger.consoleLog {
		log.Print(content)
	}
	if logger.outputfileLog {
		sincetime := time.Now().String()
		logger.writeLogFile([]byte(fmt.Sprintf("%s %s", sincetime[:19], content)))
	}
	return content
}

func (logger *ConsoleLogger) Warnf(format string, v ...interface{}) string {

	logger.Check()
	content := generateLogContent(getWarnLogTag(), logger.GetPosition(), format, v...)
	if logger.consoleLog {
		log.Print(content)
	}
	if logger.outputfileLog {
		sincetime := time.Now().String()
		logger.writeLogFile([]byte(fmt.Sprintf("%s %s", sincetime[:19], content)))
	}
	return content
}

func (logger *ConsoleLogger) Warnln(v ...interface{}) string {

	logger.Check()
	content := generateLogContent(getWarnLogTag(), logger.GetPosition(), "", v...)
	if logger.consoleLog {
		log.Println(content)
	}
	if logger.outputfileLog {
		sincetime := time.Now().String()
		logger.writeLogFile([]byte(fmt.Sprintf("%s %s", sincetime[:19], content)))
	}
	return content
}

func compareTime(time1, time2 time.Time) bool {
	if time1.Year() != time2.Year() {
		return false
	} else if time1.Month() != time2.Month() {
		return false
	} else if time1.Day() != time2.Day() {
		return false
	}
	return true
}
