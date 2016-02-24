/* crawler scheduler pool */
package middleware

import (
	base "circle/crawler/base"
	"errors"
	"fmt"
	"reflect"
	"sync"
)

type Pool interface {
	Take() (base.Entity, error)
	Return(entity base.Entity) error
	Total() uint32
	Used() uint32
}

type myPool struct {
	total       uint32
	etype       reflect.Type       // 池内的实例类型
	genEntity   func() base.Entity // 池内的实例的生成函数
	container   chan base.Entity   // 使用通道用于容纳实例的容器
	idContainer map[uint32]bool    // 池内实例的id列表
	mutex       sync.Mutex         // 互斥锁
}

func NewPool(
	total uint32,
	entityType reflect.Type,
	genEntity func() base.Entity) (Pool, error) {
	if total == 0 {
		errMsg := fmt.Sprintf("The pool can not be initialized! (total=%d)\n", total)
		return nil, errors.New(errMsg)
	}
	size := int(total)
	container := make(chan base.Entity, size)
	idContainer := make(map[uint32]bool)
	for i := 0; i < size; i++ {
		newEntity := genEntity()
		if entityType != reflect.TypeOf(newEntity) {
			errMsg := fmt.Sprintf("The type of result of function genEntity() is NOT %s!\n",
				entityType)
			return nil, errors.New(errMsg)
		}
		container <- newEntity
		idContainer[newEntity.Id()] = true
	}
	return &myPool{
		total:       total,
		etype:       entityType,
		genEntity:   genEntity,
		container:   container,
		idContainer: idContainer}, nil
}

func (pool *myPool) Total() uint32 {
	return pool.total
}

func (pool *myPool) Used() uint32 {
	return pool.total - uint32(len(pool.container))
}

func (pool *myPool) Take() (base.Entity, error) {
	entity, ok := <-pool.container
	if !ok {
		return nil, errors.New("The inner container is invalid!")
	}
	pool.mutex.Lock()
	defer pool.mutex.Unlock()
	//	此处添加互斥锁 保证map的原子操作
	pool.idContainer[entity.Id()] = false
	return entity, nil
}

func (pool *myPool) Return(entity base.Entity) error {
	if entity == nil {
		return errors.New("The returning entity is invalid!")
	}
	if pool.etype != reflect.TypeOf(entity) {
		errMsg := fmt.Sprintf("The type of returning entity is NOT %s!\n", pool.etype)
		return errors.New(errMsg)
	}
	//	判断该实例是否属于该池
	entityId := entity.Id()
	casRet := pool.compareAndSetForIdContainer(entityId, false, true)
	switch casRet {
	case 1:
		pool.container <- entity
		return nil
	case 0:
		errMsg := fmt.Sprintf("The entity (id=%d) is already in the pool!\n", entityId)
		return errors.New(errMsg)
	case -1:
		errMsg := fmt.Sprintf("The entity (id=%d) is illegal!\n", entityId)
		return errors.New(errMsg)
	}
	return nil
}

/*比较并设置实例ID容器中与给定实例ID对应的键值对的元素值*/
//	结果值:
//			-1	表示键值对不存在
//			 0	表示操作失败
//			 1	表示操作成功
func (pool *myPool) compareAndSetForIdContainer(
	entityId uint32,
	oldValue bool,
	newValue bool) int8 {
	pool.mutex.Lock()
	defer pool.mutex.Unlock()
	// 在单独的方法操作中用对map进行原子读取
	v, ok := pool.idContainer[entityId]
	if !ok {
		return -1
	}
	if v != oldValue {
		return 0
	}
	pool.idContainer[entityId] = newValue
	return 1
}
