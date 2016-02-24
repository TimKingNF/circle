/* crawler middleware channel manager */
package middleware

import (
	base "circle/crawler/base"
	"errors"
	"fmt"
	"sync"
)

type ChannelManagerStatus uint8

const (
	defaultChanLen uint = 20

	chanmanSummaryTemplate string = "status：%s，" +
		"requestChannel：%d/%d，" +
		"responseChannel：%d/%d，" +
		"itemChannel：%d/%d，" +
		"errorChannel：%d/%d"

	CHANNEL_MANAGER_STATUS_UNINITALIZED ChannelManagerStatus = 0 // 未初始化状态
	CHANNEL_MANAGER_STATUS_INITALIZED   ChannelManagerStatus = 1 // 已初始化状态
	CHANNEL_MANAGER_STATUS_CLOSED       ChannelManagerStatus = 2 // 已关闭状态
)

var statusNameMap = map[ChannelManagerStatus]string{
	CHANNEL_MANAGER_STATUS_UNINITALIZED: "uninitialized",
	CHANNEL_MANAGER_STATUS_INITALIZED:   "initialized",
	CHANNEL_MANAGER_STATUS_CLOSED:       "closed",
}

type ChannelManager interface {
	Init(
		channelArgs base.ChannelArgs,
		reset bool,
	) bool

	Close() bool
	ReqChan() (chan base.Request, error)
	RespChan() (chan base.Response, error)
	ItemChan() (chan base.Item, error)
	ErrorChan() (chan error, error)
	Status() ChannelManagerStatus
	Summary() string
}

type myChannelManager struct {
	channelArgs base.ChannelArgs
	reqCh       chan base.Request
	respCh      chan base.Response
	itemCh      chan base.Item
	errorCh     chan error
	status      ChannelManagerStatus
	rwmutex     sync.Mutex
}

func NewChannelManager(channelArgs base.ChannelArgs) ChannelManager {
	chanman := &myChannelManager{}
	chanman.Init(channelArgs, true)
	return chanman
}

func (chanman *myChannelManager) Init(
	channelArgs base.ChannelArgs,
	reset bool) bool {
	if err := channelArgs.Check(); err != nil {
		panic(err)
	}
	chanman.rwmutex.Lock()
	defer chanman.rwmutex.Unlock()
	//	chanman initalized
	if chanman.status == CHANNEL_MANAGER_STATUS_INITALIZED && !reset {
		return false
	}
	chanman.channelArgs = channelArgs
	chanman.reqCh = make(chan base.Request, channelArgs.ReqChanLen())
	chanman.respCh = make(chan base.Response, channelArgs.RespChanLen())
	chanman.itemCh = make(chan base.Item, channelArgs.ItemChanLen())
	chanman.errorCh = make(chan error, channelArgs.ErrorChanLen())
	chanman.status = CHANNEL_MANAGER_STATUS_INITALIZED
	return true
}

func (chanman *myChannelManager) Close() bool {
	chanman.rwmutex.Lock()
	defer chanman.rwmutex.Unlock()
	if chanman.status != CHANNEL_MANAGER_STATUS_INITALIZED {
		return false
	}
	close(chanman.reqCh)
	close(chanman.respCh)
	close(chanman.itemCh)
	close(chanman.errorCh)
	chanman.status = CHANNEL_MANAGER_STATUS_CLOSED
	return true
}

func (chanman *myChannelManager) CheckStatus() error {
	if chanman.status == CHANNEL_MANAGER_STATUS_INITALIZED {
		return nil
	}
	statusName, ok := statusNameMap[chanman.status]
	if !ok {
		statusName = fmt.Sprintf("%d", chanman.status)
	}
	errMsg := fmt.Sprintf("The undesirable status of channel manager：%s!\n", statusName)
	return errors.New(errMsg)
}

func (chanman *myChannelManager) ReqChan() (chan base.Request, error) {
	chanman.rwmutex.Lock()
	defer chanman.rwmutex.Unlock()
	if err := chanman.CheckStatus(); err != nil {
		return nil, err
	}
	return chanman.reqCh, nil
}

func (chanman *myChannelManager) RespChan() (chan base.Response, error) {
	chanman.rwmutex.Lock()
	defer chanman.rwmutex.Unlock()
	if err := chanman.CheckStatus(); err != nil {
		return nil, err
	}
	return chanman.respCh, nil
}

func (chanman *myChannelManager) ItemChan() (chan base.Item, error) {
	chanman.rwmutex.Lock()
	defer chanman.rwmutex.Unlock()
	if err := chanman.CheckStatus(); err != nil {
		return nil, err
	}
	return chanman.itemCh, nil
}

func (chanman *myChannelManager) ErrorChan() (chan error, error) {
	chanman.rwmutex.Lock()
	defer chanman.rwmutex.Unlock()
	if err := chanman.CheckStatus(); err != nil {
		return nil, err
	}
	return chanman.errorCh, nil
}

func (chanman *myChannelManager) Summary() string {
	return fmt.Sprintf(chanmanSummaryTemplate,
		statusNameMap[chanman.status],
		len(chanman.reqCh), cap(chanman.reqCh),
		len(chanman.respCh), cap(chanman.respCh),
		len(chanman.itemCh), cap(chanman.itemCh),
		len(chanman.errorCh), cap(chanman.errorCh))
}

func (chanman *myChannelManager) Status() ChannelManagerStatus {
	return chanman.status
}
