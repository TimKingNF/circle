/* scheduler pool args */
package base

import (
	"errors"
	"fmt"
)

type SchePoolArgs struct {
	downLoaderPoolSize uint32 // 网页下载器池的尺寸
	analyzerPoolSize   uint32 // 分析器池的尺寸
	description        string // 描述
}

const (
	SchePoolArgsTemplate string = "{ DownloaderPoolSize: %d," +
		" analyzerPoolSize: %d }"
)

func NewSchePoolArgs(downLoaderPoolSize, analyzerPoolSize uint32) SchePoolArgs {
	return SchePoolArgs{
		downLoaderPoolSize: downLoaderPoolSize,
		analyzerPoolSize:   analyzerPoolSize,
	}
}

func (args *SchePoolArgs) Check() error {
	if args.downLoaderPoolSize == 0 {
		return errors.New("The downloader pool size can not be 0!\n")
	}
	if args.analyzerPoolSize == 0 {
		return errors.New("The analyzer pool size can not be 0!\n")
	}
	return nil
}

func (args *SchePoolArgs) String() string {
	if len(args.description) == 0 {
		args.genDescription()
	}
	return args.description
}

/*获得网页下载器池的尺寸*/
func (args *SchePoolArgs) DownloaderPoolSize() uint32 {
	return args.downLoaderPoolSize
}

/*获得分析器池的尺寸*/
func (args *SchePoolArgs) AnalyzerPoolSize() uint32 {
	return args.analyzerPoolSize
}

func (args *SchePoolArgs) genDescription() {
	args.description = fmt.Sprintf(SchePoolArgsTemplate,
		args.downLoaderPoolSize,
		args.analyzerPoolSize)
}
