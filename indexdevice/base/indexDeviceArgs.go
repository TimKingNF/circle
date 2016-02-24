package base

import (
	"errors"
	"fmt"
)

type IndexDeviceArgs struct {
	dividerAddr string // 远程分配器地址
	description string
}

const (
	indexDeviceCtrlMArgsTemplate string = "{ dividerAddr: %s }"
)

func NewIndexDeviceArgs(
	dividerAdder string) IndexDeviceArgs {
	return IndexDeviceArgs{
		dividerAddr: dividerAdder,
	}
}

func (args *IndexDeviceArgs) Check() error {
	if len(args.dividerAddr) == 0 {
		return errors.New("The divider addr can not be nil.(index device)")
	}
	return nil
}

func (args *IndexDeviceArgs) String() string {
	if len(args.description) == 0 {
		args.genDescription()
	}
	return args.description
}

func (args *IndexDeviceArgs) DividerAddr() string {
	return args.dividerAddr
}

func (args *IndexDeviceArgs) genDescription() {
	args.description = fmt.Sprintf(indexDeviceCtrlMArgsTemplate,
		args.dividerAddr)
}
