/* scheduler channelManager args */
package base

import (
	"errors"
	"fmt"
)

type ChannelArgs struct {
	reqChanLen   uint   // 请求通道的长度
	respChanLen  uint   // 响应通道的长度
	itemChanLen  uint   // 条目通道的长度
	errorChanLen uint   // 错误通道的长度
	description  string // 描述
}

const (
	channelArgsTemplate string = "{ reqChanLen: %d, respChanLen: %d," +
		" itemChanLen: %d, errorChanLen: %d }"
)

func NewChannelArgs(reqChanLen, respChanLen, itemChanLen, errorChanLen uint) ChannelArgs {
	return ChannelArgs{
		reqChanLen:   reqChanLen,
		respChanLen:  respChanLen,
		itemChanLen:  itemChanLen,
		errorChanLen: errorChanLen,
	}
}

func (args *ChannelArgs) Check() error {
	if args.reqChanLen == 0 {
		return errors.New("The request channel max length (capacity) can not be 0!\n")
	}
	if args.respChanLen == 0 {
		return errors.New("The response channel max length (capacity) can not be 0!\n")
	}
	if args.itemChanLen == 0 {
		return errors.New("The item channel max length (capacity) can not be 0!\n")
	}
	if args.errorChanLen == 0 {
		return errors.New("The error channel max length (capacity) can not be 0!\n")
	}
	return nil
}

func (args *ChannelArgs) String() string {
	if len(args.description) == 0 {
		args.genDescription()
	}
	return args.description
}

/*获得请求通道的长度*/
func (args *ChannelArgs) ReqChanLen() uint {
	return args.reqChanLen
}

/*获得响应通道的长度*/
func (args *ChannelArgs) RespChanLen() uint {
	return args.respChanLen
}

/*获得条目通道的长度*/
func (args *ChannelArgs) ItemChanLen() uint {
	return args.itemChanLen
}

/*获得错误通道的长度*/
func (args *ChannelArgs) ErrorChanLen() uint {
	return args.errorChanLen
}

func (args *ChannelArgs) genDescription() {
	args.description = fmt.Sprintf(channelArgsTemplate,
		args.reqChanLen,
		args.respChanLen,
		args.itemChanLen,
		args.errorChanLen)
}
