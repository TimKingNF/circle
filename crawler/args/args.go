/* crawler args */
package args

import (
	cmn "circle/common"
	base "circle/crawler/base"
	"time"
)

var (
	channelArgs base.ChannelArgs = base.NewChannelArgs(10, 10, 10, 10)

	schePoolArgs base.SchePoolArgs = base.NewSchePoolArgs(3, 3)

	intervalNs    = 10 * time.Millisecond // 10 毫秒
	maxIdleCount  = uint(1000)
	autoStop      = true
	detailSummary = true

	moitorArgs base.MonitorArgs = base.NewMonitorArgs(intervalNs, maxIdleCount,
		autoStop, detailSummary)

	crawlDepth  = uint32(1)
	crossDomain = false

	spiderArgs base.SpiderArgs = base.NewSpiderArgs(crossDomain, crawlDepth)

	consoleLog       = true
	outputfileLog    = true
	outputfilePath   = ""
	outputfilePrefix = "crawler"

	loggerArgs cmn.LoggerArgs = cmn.NewLoggerArgs(consoleLog, outputfileLog,
		outputfilePath, outputfilePrefix)

	connected_divider_address = "127.0.0.1:8085"

	CrawlerArgs base.CrawlerArgs = base.NewCrawlerArgs(
		channelArgs,
		schePoolArgs,
		moitorArgs,
		spiderArgs,
		loggerArgs,
		connected_divider_address,
	)
)

const (
	DEVICE_NAME = cmn.DEVICE_CRAWLER
)

func Reset(args base.CrawlerArgs) error {
	if err := args.Check(); err != nil {
		return err
	}
	CrawlerArgs = args
	return nil
}
