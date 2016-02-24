/* circle logging file control */
package logging

import (
	// "io"
	"os"
	"runtime"
	"strings"
)

func (logger *ConsoleLogger) writeLogFile(data []byte) error {
	logger.mutex.Lock()
	defer logger.mutex.Unlock()

	data = append(data, []byte("\r\n\r\n")...)
	_, err := logger.filePtr.Write(data)
	if err != nil {
		return err
	}
	return nil
}

func pathIsExist(path string) bool {
	_, err := os.Stat(path)
	return err == nil || os.IsExist(err)
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
