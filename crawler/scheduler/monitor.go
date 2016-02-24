/* crawler monitor */
package scheduler

import (
	"errors"
	"fmt"
	"runtime"
	"time"
)

//	@param level: log level
type Record func(level byte, content string)

const (
	SCHEDULER_MONITOR_RECORD_COMMON byte = 0
	SCHEDULER_MONITOR_RECORD_NOTICE byte = 1
	SCHEDULER_MONITOR_RECORD_ERROR  byte = 2

	summaryForMonitoring = "Monitor - Collected information[%d]:\n" +
		"  Goroutine number: %d\n" +
		"  Scheduler:\n%s" +
		"  Escaped time: %s\n"

	msgReachMaxIdleCount = "The scheduler has been idle for a period of time" +
		" (about %s)." +
		" Now consider what stop it.\n"

	msgStopScheduler = "Stop scheduler...%s.\n"
)

//	monitor
//	when monitor stop, return chan (check idle times)
func Monitoring(
	scheduler Scheduler,

	intervalNs time.Duration, // check time interval

	maxIdleCount uint, // max ilde count

	autoStop bool, // after time=(intervalNs * maxIdleCount), auto stop scheduler

	detailSummary bool, // more summary

	record Record,
) <-chan uint64 {
	if scheduler == nil {
		panic(errors.New("The scheduler is invalid!"))
	}
	//	when param too small
	if intervalNs < time.Millisecond {
		intervalNs = time.Millisecond
	}
	if maxIdleCount < 1000 {
		maxIdleCount = 1000
	}
	stopNotifier := make(chan byte, 1)
	reportError(scheduler, record, stopNotifier)
	recordSummary(scheduler, detailSummary, record, stopNotifier)
	checkCountChan := make(chan uint64, 2)
	checkStatus(
		scheduler,
		intervalNs,
		maxIdleCount,
		autoStop,
		checkCountChan,
		record,
		stopNotifier)
	return checkCountChan
}

func checkStatus(
	scheduler Scheduler,
	intervalNs time.Duration,
	maxIdleCount uint,
	autoStop bool,
	checkCountChan chan<- uint64,
	record Record,
	stopNotifier chan<- byte) {
	var checkCount uint64
	go func() {
		defer func() {
			stopNotifier <- 1
			stopNotifier <- 2
			checkCountChan <- checkCount
		}()
		waitForSchedulerStart(scheduler)
		var idleCount uint
		var firstIdleTime time.Time
		for {
			if scheduler.Idle() {
				idleCount++
				if idleCount == 1 {
					firstIdleTime = time.Now()
				}
				if idleCount >= maxIdleCount {
					msg := fmt.Sprintf(msgReachMaxIdleCount,
						time.Since(firstIdleTime).String())
					record(SCHEDULER_MONITOR_RECORD_COMMON, msg)
					//	check more once
					if scheduler.Idle() {
						if autoStop {
							var result string = "failing"
							if scheduler.Stop() {
								result = "success"
							}
							msg = fmt.Sprintf(msgStopScheduler, result)
							record(SCHEDULER_MONITOR_RECORD_COMMON, msg)
						}
						break
					} else {
						if idleCount > 0 {
							idleCount = 0
						}
					}
				}
			} else {
				if idleCount > 0 {
					idleCount = 0
				}
			}
			checkCount++
			time.Sleep(intervalNs)
		}
	}()
}

func reportError(
	scheduler Scheduler,
	record Record,
	stopNotifier <-chan byte) {
	go func() {
		waitForSchedulerStart(scheduler)
		for {
			select {
			case <-stopNotifier:
				return
			default:
			}
			errorChan := scheduler.ErrorChan()
			if errorChan == nil {
				return
			}
			err := <-errorChan
			if err != nil {
				errMsg := fmt.Sprintf("Error (received from error channel): %s", err)
				record(SCHEDULER_MONITOR_RECORD_ERROR, errMsg)
			}
			time.Sleep(time.Microsecond)
		}
	}()
}

func waitForSchedulerStart(scheduler Scheduler) {
	for !scheduler.Running() {
		time.Sleep(time.Microsecond)
	}
}

func recordSummary(
	scheduler Scheduler,
	detailSummary bool,
	record Record,
	stopNotifier <-chan byte) {
	go func() {
		waitForSchedulerStart(scheduler)
		var recordCount uint64 = 1
		startTime := time.Now()
		var prevSchedSummary SchedSummary
		var prevNumGoroutine int
		for {
			select {
			case <-stopNotifier:
				return
			default:
			}

			//	runtime goroutine number
			currNumGoroutine := runtime.NumGoroutine()
			currSchedSummary := scheduler.Summary("     ")
			//	比对前后两份摘要信息的一致性。只有不一致时才会记录
			if currNumGoroutine != prevNumGoroutine ||
				!currSchedSummary.Same(prevSchedSummary) {
				schedSummaryStr := func() string {
					if detailSummary {
						return currSchedSummary.Detail()
					} else {
						return currSchedSummary.String()
					}
				}()
				info := fmt.Sprintf(summaryForMonitoring,
					recordCount,
					currNumGoroutine,
					schedSummaryStr,
					time.Since(startTime).String())
				record(SCHEDULER_MONITOR_RECORD_COMMON, info)
				prevNumGoroutine = currNumGoroutine
				prevSchedSummary = currSchedSummary
				recordCount++
			}
			time.Sleep(time.Microsecond)
		}
	}()
}
