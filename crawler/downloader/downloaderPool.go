/* crawler scheduler downloaderPool */
package downloader

import (
	base "circle/crawler/base"
	mdw "circle/crawler/middleware"
	"errors"
	"fmt"
	"reflect"
)

type GenDownloader func() Downloader

type DownloaderPool interface {
	Total() uint32
	Used() uint32
	Take() (Downloader, error)
	Return(dl Downloader) error
}

type myDownloaderPool struct {
	pool  mdw.Pool
	etype reflect.Type
}

func NewDownloaderPool(
	total uint32,
	gen GenDownloader) (DownloaderPool, error) {
	etype := reflect.TypeOf(gen())
	genEntity := func() base.Entity {
		return gen()
	}
	pool, err := mdw.NewPool(total, etype, genEntity)
	if err != nil {
		return nil, err
	}
	return &myDownloaderPool{
		pool:  pool,
		etype: etype,
	}, nil
}

func (dlpool *myDownloaderPool) Take() (Downloader, error) {
	entity, err := dlpool.pool.Take()
	if err != nil {
		return nil, err
	}
	dl, ok := entity.(Downloader)
	if !ok {
		errMsg := fmt.Sprintf("The type of entity is NOT %s!\n", dlpool.etype)
		panic(errors.New(errMsg))
	}
	return dl, nil
}

func (dlpool *myDownloaderPool) Return(dl Downloader) error {
	return dlpool.pool.Return(dl)
}

func (dlpool *myDownloaderPool) Total() uint32 {
	return dlpool.pool.Total()
}

func (dlpool *myDownloaderPool) Used() uint32 {
	return dlpool.pool.Used()
}
