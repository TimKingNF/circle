/* crawler monitor args */
package base

import (
	"errors"
	"fmt"
	"time"
)

const (
	/*调度器监控参数容器的描述模板*/
	monitorArgsTemplate string = "{ intervalNs: %s, maxIdleCount: %d," +
		" autoStop: %v, detailSummary: %v }"
)

type MonitorArgs struct {
	intervalNs   time.Duration // 监控检查间隔
	maxIdleCount uint          // 最大空闲计数

	// 该参数用来指示该方法是否在调度器空闲一段时间 (即持续空闲时间，由intervalNs * maxIdleCount)之后自行停止调度器
	autoStop bool

	detailSummary bool   // 是否需要详细的摘要信息
	description   string // 描述
}

func NewMonitorArgs(intervalNs time.Duration, maxIdleCount uint, autoStop bool,
	detailSummary bool) MonitorArgs {
	return MonitorArgs{
		intervalNs:    intervalNs,
		maxIdleCount:  maxIdleCount,
		autoStop:      autoStop,
		detailSummary: detailSummary,
	}
}

func (args *MonitorArgs) Check() error {
	if args.intervalNs == 0 {
		return errors.New("The monitor check interval time can not be 0!\n")
	}
	if args.maxIdleCount == 0 {
		return errors.New("The monitor check times can not be 0!\n")
	}
	return nil
}

func (args *MonitorArgs) String() string {
	if len(args.description) == 0 {
		args.genDescription()
	}
	return args.description
}

/*获取调度器监控参数容器的检查时间间隔*/
func (args *MonitorArgs) IntervalNs() time.Duration {
	return args.intervalNs
}

/*获取调度器监控参数容器的最大空闲次数*/
func (args *MonitorArgs) MaxIdleCount() uint {
	return args.maxIdleCount
}

func (args *MonitorArgs) AutoStop() bool {
	return args.autoStop
}

func (args *MonitorArgs) DetailSummary() bool {
	return args.detailSummary
}

func (args *MonitorArgs) genDescription() {
	args.description = fmt.Sprintf(monitorArgsTemplate,
		args.intervalNs,
		args.maxIdleCount,
		args.autoStop,
		args.detailSummary)
}
