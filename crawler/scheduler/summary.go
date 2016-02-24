/* crawler scheduler summary */
package scheduler

import (
	"bytes"
	base "circle/crawler/base"
	"fmt"
)

type SchedSummary interface {
	String() string
	Detail() string
	Same(other SchedSummary) bool
}

type mySchedSummary struct {
	prefix              string
	running             uint32
	channelArgs         base.ChannelArgs
	schePoolArgs        base.SchePoolArgs
	spiderArgs          base.SpiderArgs
	chanmanSummary      string
	reqCacheSummary     string
	itemPipelineSummary string
	stopSignSummary     string
	dlPoolLen           uint32
	dlPoolCap           uint32
	analyzerPoolLen     uint32
	analyzerPoolCap     uint32
	urlCount            int    // download url count
	urlDetail           string // download url summary
}

func NewSchedSummary(sched *myScheduler, prefix string) SchedSummary {
	if sched == nil {
		return nil
	}
	urlCount := len(sched.urlMap)
	var urlDetail string
	if urlCount > 0 {
		var buffer bytes.Buffer
		buffer.WriteString("\n")
		for k, _ := range sched.urlMap {
			buffer.WriteString(prefix)
			buffer.WriteString(prefix)
			buffer.WriteString(k)
			buffer.WriteString("\n")
		}
		urlDetail = buffer.String()
	} else {
		urlDetail = "\n"
	}
	return &mySchedSummary{
		prefix:              prefix,
		running:             sched.running,
		schePoolArgs:        sched.schePoolArgs,
		channelArgs:         sched.channelArgs,
		spiderArgs:          sched.spiderArgs,
		chanmanSummary:      sched.chanman.Summary(),
		reqCacheSummary:     sched.reqCache.Summary(),
		itemPipelineSummary: sched.itemPipeline.Summary(),
		stopSignSummary:     sched.stopSign.Summary(),
		dlPoolLen:           sched.dlpool.Used(),
		dlPoolCap:           sched.dlpool.Total(),
		analyzerPoolLen:     sched.analyzerPool.Used(),
		analyzerPoolCap:     sched.analyzerPool.Total(),
		urlCount:            urlCount,
		urlDetail:           urlDetail,
	}
}

func (ss *mySchedSummary) String() string {
	return ss.getSummary(false)
}

func (ss *mySchedSummary) Detail() string {
	return ss.getSummary(true)
}

func (ss *mySchedSummary) Same(other SchedSummary) bool {
	if other == nil {
		return false
	}
	otherSs, ok := interface{}(other).(*mySchedSummary)
	if !ok {
		return false
	}
	if ss.running != otherSs.running ||
		ss.schePoolArgs.DownloaderPoolSize() != otherSs.schePoolArgs.DownloaderPoolSize() ||
		ss.schePoolArgs.AnalyzerPoolSize() != otherSs.schePoolArgs.AnalyzerPoolSize() ||
		ss.channelArgs.ReqChanLen() != otherSs.channelArgs.ReqChanLen() ||
		ss.channelArgs.RespChanLen() != otherSs.channelArgs.RespChanLen() ||
		ss.channelArgs.ItemChanLen() != otherSs.channelArgs.ItemChanLen() ||
		ss.channelArgs.ErrorChanLen() != otherSs.channelArgs.ErrorChanLen() ||
		ss.spiderArgs.CrossDomain() != otherSs.spiderArgs.CrossDomain() ||
		ss.spiderArgs.CrawlDepth() != otherSs.spiderArgs.CrawlDepth() ||
		ss.dlPoolLen != otherSs.dlPoolLen ||
		ss.dlPoolCap != otherSs.dlPoolCap ||
		ss.analyzerPoolLen != otherSs.analyzerPoolLen ||
		ss.analyzerPoolCap != otherSs.analyzerPoolCap ||
		ss.urlCount != otherSs.urlCount ||
		ss.stopSignSummary != otherSs.stopSignSummary ||
		ss.reqCacheSummary != otherSs.reqCacheSummary ||
		ss.itemPipelineSummary != otherSs.itemPipelineSummary ||
		ss.chanmanSummary != otherSs.chanmanSummary {
		return false
	}
	return true
}

func (ss *mySchedSummary) getSummary(detail bool) string {
	prefix := ss.prefix
	template := prefix + "Running: %s \n" +
		prefix + "Pool base args: %s \n" +
		prefix + "Channel args: %s \n" +
		prefix + "Spider args: %s \n" +
		prefix + "Channels manager: %s \n" +
		prefix + "Request cache: %s \n" +
		prefix + "Downloader pool: %d/%d \n" +
		prefix + "Analyzer pool: %d/%d \n" +
		prefix + "Item pipeline: %s \n" +
		prefix + "Urls(%d): %s" +
		prefix + "Stop sign: %s \n"
	return fmt.Sprintf(template,
		schedStatusMap[ss.running],
		ss.schePoolArgs.String(),
		ss.channelArgs.String(),
		ss.spiderArgs.String(),
		ss.chanmanSummary,
		ss.reqCacheSummary,
		ss.dlPoolLen, ss.dlPoolCap,
		ss.analyzerPoolLen, ss.analyzerPoolCap,
		ss.itemPipelineSummary,
		ss.urlCount,
		func() string {
			if detail {
				return ss.urlDetail
			} else {
				return "<concealed>\n"
			}
		}(),
		ss.stopSignSummary)
}
