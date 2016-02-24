/* crawler error type */
package base

import (
	"bytes"
	"fmt"
)

type ErrorType string

const (
	DOWNLOADER_ERROR     ErrorType = "Downloader Error"
	ANALYZER_ERROR       ErrorType = "Analyer Error"
	ITEM_PROCESSOR_ERROR ErrorType = "Item Processor Error"
	SCHEDULER_ERROR      ErrorType = "Scheduler Error"
)

type CrawlerError interface {
	Type() ErrorType
	Error() string
}

type myCrawlerError struct {
	errType    ErrorType
	errMsg     string
	fullErrMsg string
}

func NewCrawlerError(errType ErrorType, errMsg string) CrawlerError {
	return &myCrawlerError{errType: errType, errMsg: errMsg}
}

func (ce *myCrawlerError) Type() ErrorType {
	return ce.errType
}

func (ce *myCrawlerError) Error() string {
	if len(ce.fullErrMsg) <= 0 {
		ce.getFullErrMsg()
	}
	return ce.fullErrMsg
}

func (ce *myCrawlerError) getFullErrMsg() {
	var buffer bytes.Buffer
	buffer.WriteString("Crawler Error: ")
	if len(ce.errType) > 0 {
		buffer.WriteString(string(ce.errType))
		buffer.WriteString(": ")
	}
	buffer.WriteString(ce.errMsg)
	ce.fullErrMsg = fmt.Sprintf("%s\n", buffer.String())
	return
}
