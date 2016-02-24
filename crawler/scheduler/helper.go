/* crawler scheduler same helper method */
package scheduler

import (
	cmn "circle/common"
	anlz "circle/crawler/analyzer"
	base "circle/crawler/base"
	dl "circle/crawler/downloader"
	ipl "circle/crawler/itempipeline"
	mdw "circle/crawler/middleware"
	"fmt"
	"strings"
)

func getPrimaryDomain(host string) (string, error) {
	return cmn.GetPrimaryDomain(host)
}

func parseCode(code string) []string {
	ret := make([]string, 2)
	var codePrefix, id string
	index := strings.Index(code, "-")
	if index > 0 {
		codePrefix = code[:index]
		id = code[index+1:]
	} else {
		codePrefix = code
	}
	ret[0] = codePrefix
	ret[1] = id
	return ret
}

func generateChannelManager(channelArgs base.ChannelArgs) mdw.ChannelManager {
	return mdw.NewChannelManager(channelArgs)
}

func generateDownloaderPool(
	poolSize uint32,
	httpClientGenerator GenHttpClient) (dl.DownloaderPool, error) {
	var httpClient = httpClientGenerator()
	dlpool, err := dl.NewDownloaderPool(poolSize, func() dl.Downloader {
		return dl.NewDownloader(httpClient)
	})
	if err != nil {
		return nil, err
	}
	return dlpool, nil
}

func generateAnalyzerPool(poolSize uint32) (anlz.AnalyzerPool, error) {
	alzpool, err := anlz.NewAnalyzerPool(poolSize, func() anlz.Analyzer {
		return anlz.NewAnalyzer()
	})
	if err != nil {
		return nil, err
	}
	return alzpool, nil
}

func generateItemPipeline(itemProcessors []ipl.ProcessItem) ipl.ItemPipeline {
	return ipl.NewItemPipeline(itemProcessors)
}

//	gen implements code
func generateCode(prefix string, id uint32) string {
	return fmt.Sprintf("%s-%d", prefix, id)
}
