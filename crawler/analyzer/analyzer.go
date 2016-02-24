/* crawler scheduler analyzerPool analyzer */
package analyzer

import (
	cmn "circle/common"
	args "circle/crawler/args"
	base "circle/crawler/base"
	"circle/logging"
	"errors"
	"fmt"
	"net/url"
)

var (
	analyzerIdGenertor cmn.IdGenertor = cmn.NewIdGenertor()
	logger             logging.Logger
)

func genLogger() logging.Logger {
	if logger == nil {
		var loggerArgs cmn.LoggerArgs = args.CrawlerArgs.LoggerArgs()

		logger = cmn.NewLogger(cmn.NewLoggerArgs(
			loggerArgs.ConsoleLog(),
			loggerArgs.OutputfileLog(),
			loggerArgs.OutputfilePath(),
			loggerArgs.OutputfilePrefix()+"_analyzer",
		))
	}
	return logger
}

type Analyzer interface {
	Id() uint32
	Analyze(
		respParsers []ParseResponse,
		resp base.Response) ([]base.Data, []error)
}

type myAnalyzer struct {
	id uint32
}

func NewAnalyzer() Analyzer {
	return &myAnalyzer{id: genAnalyzerId()}
}

func genAnalyzerId() uint32 {
	return analyzerIdGenertor.GenUint32()
}

func (analyzer *myAnalyzer) Id() uint32 {
	return analyzer.id
}

func (analyzer *myAnalyzer) Analyze(
	respParsers []ParseResponse,
	resp base.Response) ([]base.Data, []error) {
	if respParsers == nil {
		err := errors.New("The response parser list is invalid!")
		return nil, []error{err}
	}
	httpResp := resp.HttpResp()
	if httpResp == nil {
		err := errors.New("The http response is invalid!")
		return nil, []error{err}
	}
	var reqUrl *url.URL = httpResp.Request.URL
	genLogger().Infof("Parse the response (reqUrl=%s)... \n", reqUrl)

	//	parse http response
	dataList := make([]base.Data, 0)
	errorList := make([]error, 0)

	for i, respParser := range respParsers {
		if respParser == nil {
			err := errors.New(fmt.Sprintf("The document parser [%d] is invalid!\n", i))
			errorList = append(errorList, err)
			continue
		}
		pDataList, pErrorList := respParser(resp)

		if pDataList != nil {
			for _, pData := range pDataList {
				dataList = appendDataList(dataList, pData, resp.Depth())
			}
		}
		if pErrorList != nil {
			for _, pError := range pErrorList {
				errorList = appendErrorList(errorList, pError)
			}
		}
	}
	return dataList, errorList
}

func appendDataList(dataList []base.Data, data base.Data, respDepth uint32) []base.Data {
	if data == nil {
		return dataList
	}
	req, ok := data.(*base.Request)
	//	the data isn't http request
	if !ok {
		return append(dataList, data)
	}
	//	old response take a new http request, and the depth add one
	newDepth := respDepth + 1
	if req.MaxDepth() != 0 {
		if req.Depth() != newDepth {
			req = base.NewMaxRequest(req.HttpReq(), req.PrimaryDomain(), newDepth, req.MaxDepth())
		}
	} else {
		if req.Depth() != newDepth {
			req = base.NewRequest(req.HttpReq(), req.PrimaryDomain(), newDepth)
		}
	}
	return append(dataList, req)
}

func appendErrorList(errorList []error, err error) []error {
	if err == nil {
		return errorList
	}
	return append(errorList, err)
}
