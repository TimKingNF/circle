package base

import (
	"errors"
	"fmt"
)

type DatabaseArgs struct {
	dividerAddr string // 远程分配器地址
	sqlAddr     string // 数据库地址
	dbname      string // 数据库名称
	description string
}

const (
	databaseCtrlMArgsTemplate string = "{ dividerAddr: %s\n, sqlAddr: %s\n, dbname: %s }"
)

func NewDatabaseArgs(
	dividerAddr,
	sqlAddr,
	dbname string) DatabaseArgs {
	return DatabaseArgs{
		dividerAddr: dividerAddr,
		sqlAddr:     sqlAddr,
		dbname:      dbname,
	}
}

func (args *DatabaseArgs) Check() error {
	if len(args.dividerAddr) == 0 {
		return errors.New("The divider addr can not be nil.(database)")
	}
	if len(args.sqlAddr) == 0 {
		return errors.New("The sql addr can not be nil.(database)")
	}
	if len(args.dbname) == 0 {
		return errors.New("The sql db name can not be nil.(database)")
	}
	return nil
}

func (args *DatabaseArgs) String() string {
	if len(args.description) == 0 {
		args.genDescription()
	}
	return args.description
}

func (args *DatabaseArgs) DividerAddr() string {
	return args.dividerAddr
}

func (args *DatabaseArgs) SqlAddr() string {
	return args.sqlAddr
}

func (args *DatabaseArgs) Dbname() string {
	return args.dbname
}

func (args *DatabaseArgs) genDescription() {
	args.description = fmt.Sprintf(databaseCtrlMArgsTemplate,
		args.dividerAddr,
		args.sqlAddr,
		args.dbname)
}
