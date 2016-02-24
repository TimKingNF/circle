package base

import (
	"errors"
	"fmt"
)

type SpiderArgs struct {
	crossDomain bool   // 是否跨域
	crawlDepth  uint32 // 能允许的被爬取的网页的最大深度值，大于该值的网页将会被忽略
	description string // 描述
}

const (
	spiderArgsTemplate string = "{ crossDomain: %v, crawlDepth: %d }"
)

func NewSpiderArgs(crossDomain bool, crawlDepth uint32) SpiderArgs {
	return SpiderArgs{
		crossDomain: crossDomain,
		crawlDepth:  crawlDepth,
	}
}

func (args *SpiderArgs) Check() error {
	if args.crawlDepth == 0 {
		return errors.New("The scheduler crawl depth can not be 0!\n")
	}
	return nil
}

func (args *SpiderArgs) String() string {
	if len(args.description) == 0 {
		args.genDescription()
	}
	return args.description
}

func (args *SpiderArgs) CrossDomain() bool {
	return args.crossDomain
}

func (args *SpiderArgs) CrawlDepth() uint32 {
	return args.crawlDepth
}

func (args *SpiderArgs) genDescription() {
	args.description = fmt.Sprintf(spiderArgsTemplate,
		args.crossDomain,
		args.crawlDepth)
}
