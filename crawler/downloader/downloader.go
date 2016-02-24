/* crawler scheduler downloaderPool downloader */
package downloader

import (
	cmn "circle/common"
	args "circle/crawler/args"
	base "circle/crawler/base"
	"circle/logging"
	"net/http"
)

var (
	downloaderIdGenertor cmn.IdGenertor = cmn.NewIdGenertor()

	logger logging.Logger
)

func genLogger() logging.Logger {
	if logger == nil {
		var loggerArgs cmn.LoggerArgs = args.CrawlerArgs.LoggerArgs()

		logger = cmn.NewLogger(cmn.NewLoggerArgs(
			loggerArgs.ConsoleLog(),
			loggerArgs.OutputfileLog(),
			loggerArgs.OutputfilePath(),
			loggerArgs.OutputfilePrefix()+"_downloader"))
	}
	return logger
}

type Downloader interface {
	Id() uint32
	Download(req base.Request) (*base.Response, error)
}

type myDownloader struct {
	id         uint32
	httpClient http.Client
}

func NewDownloader(client *http.Client) Downloader {
	id := genDownloaderId()
	if client == nil {
		client = &http.Client{}
	}
	return &myDownloader{
		id:         id,
		httpClient: *client,
	}
}

func (dl *myDownloader) Id() uint32 {
	return dl.id
}

func (dl *myDownloader) Download(req base.Request) (*base.Response, error) {
	httpReq := req.HttpReq()
	httpResp, err := dl.httpClient.Do(httpReq)
	if err != nil {
		return nil, err
	}
	genLogger().Infof("Download the request (reqUrl=%s)... \n", httpReq.URL)
	if req.MaxDepth() != 0 {
		return base.NewUpdateResponse(httpResp, req.PrimaryDomain(), req.Depth()), nil
	}
	return base.NewResponse(httpResp, req.PrimaryDomain(), req.Depth()), nil
}

func genDownloaderId() uint32 {
	return downloaderIdGenertor.GenUint32()
}
