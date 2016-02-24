/* crawler init args */
package base

import (
	cmn "circle/common"
	"errors"
	"fmt"
)

type CrawlerArgs struct {
	channelArgs  ChannelArgs    // 通道参数容器
	schePoolArgs SchePoolArgs   // 调度器内部池基本参数容器
	monitorArgs  MonitorArgs    // 监控参数容器
	spiderArgs   SpiderArgs     // 抓取参数容器
	loggerArgs   cmn.LoggerArgs // 日志参数容器
	dividerAddr  string         // 远程分配器地址
	description  string         // 描述
}

const (
	crawlerCtrlMArgsTemplate string = "{ channelArgs: %s,\n schePoolArgs: %s,\n" +
		" monitorArgs: %s,\n spiderArgs: %s,\n loggerArgs: %s,\n dividerAddr: %s }"
)

func NewCrawlerArgs(channelArgs ChannelArgs,
	schePoolArgs SchePoolArgs,
	monitorArgs MonitorArgs,
	spiderArgs SpiderArgs,
	loggerArgs cmn.LoggerArgs,
	dividerAddr string) CrawlerArgs {
	return CrawlerArgs{
		channelArgs:  channelArgs,
		schePoolArgs: schePoolArgs,
		monitorArgs:  monitorArgs,
		spiderArgs:   spiderArgs,
		loggerArgs:   loggerArgs,
		dividerAddr:  dividerAddr,
	}
}

func (args *CrawlerArgs) Check() error {
	if e := args.channelArgs.Check(); e != nil {
		return e
	}
	if e := args.schePoolArgs.Check(); e != nil {
		return e
	}
	if e := args.monitorArgs.Check(); e != nil {
		return e
	}
	if e := args.spiderArgs.Check(); e != nil {
		return e
	}
	if e := args.loggerArgs.Check(); e != nil {
		return e
	}
	if len(args.dividerAddr) == 0 {
		return errors.New("The divider addr can not be nil.(crawler)")
	}
	return nil
}

func (args *CrawlerArgs) String() string {
	if len(args.description) == 0 {
		args.genDescription()
	}
	return args.description
}

func (args *CrawlerArgs) ChannelArgs() ChannelArgs {
	return args.channelArgs
}

func (args *CrawlerArgs) SchePoolArgs() SchePoolArgs {
	return args.schePoolArgs
}

func (args *CrawlerArgs) MonitorArgs() MonitorArgs {
	return args.monitorArgs
}

func (args *CrawlerArgs) SpiderArgs() SpiderArgs {
	return args.spiderArgs
}

func (args *CrawlerArgs) LoggerArgs() cmn.LoggerArgs {
	return args.loggerArgs
}

func (args *CrawlerArgs) DividerAddr() string {
	return args.dividerAddr
}

func (args *CrawlerArgs) genDescription() {
	args.description = fmt.Sprintf(crawlerCtrlMArgsTemplate,
		args.channelArgs.String(),
		args.schePoolArgs.String(),
		args.monitorArgs.String(),
		args.spiderArgs.String(),
		args.loggerArgs.String(),
		args.dividerAddr)
}
