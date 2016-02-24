/* crawler scheduler analyzerPool */
package analyzer

import (
	base "circle/crawler/base"
	mdw "circle/crawler/middleware"
	"errors"
	"fmt"
	"reflect"
)

type GenAnalyzer func() Analyzer

type AnalyzerPool interface {
	Total() uint32
	Used() uint32
	Take() (Analyzer, error)
	Return(analyzer Analyzer) error
}

type myAnalyzerPool struct {
	pool  mdw.Pool
	etype reflect.Type
}

func NewAnalyzerPool(
	total uint32,
	gen GenAnalyzer) (AnalyzerPool, error) {
	etype := reflect.TypeOf(gen())
	genEntity := func() base.Entity {
		return gen()
	}
	pool, err := mdw.NewPool(total, etype, genEntity)
	if err != nil {
		return nil, err
	}
	alzpool := &myAnalyzerPool{pool: pool, etype: etype}
	return alzpool, nil
}

func (alzpool *myAnalyzerPool) Take() (Analyzer, error) {
	entity, err := alzpool.pool.Take()
	if err != nil {
		return nil, err
	}
	alz, ok := entity.(Analyzer)
	if !ok {
		errMsg := fmt.Sprintf("The type of entity is NOT %s!\n", alzpool.etype)
		panic(errors.New(errMsg))
	}
	return alz, nil
}

func (alzpool *myAnalyzerPool) Return(alz Analyzer) error {
	return alzpool.pool.Return(alz)
}

func (alzpool *myAnalyzerPool) Total() uint32 {
	return alzpool.pool.Total()
}

func (alzpool *myAnalyzerPool) Used() uint32 {
	return alzpool.pool.Used()
}
