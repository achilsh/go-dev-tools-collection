package tools

import (
	"fmt"
	"sync"
)


// 泛型接口
type BizProcesser[IN any, OUT any] interface {
	Handle(in IN) (OUT, error)
}

// 注册中心（使用 sync.Map 防止并发冲突）
type registryKey struct {
	inType  string
	outType string
	tag     int
}

var (
	processerRegistry = sync.Map{} // map[registryKey]any
)

// 注册方法
func RegisterBizProcesser[IN any, OUT any](tag int, p BizProcesser[IN, OUT]) {
	key := registryKey{
		inType:  fmt.Sprintf("%T", *new(IN)),
		outType: fmt.Sprintf("%T", *new(OUT)),
		tag:     tag,
	}
	fmt.Printf("key:%+v\n ", key)
	processerRegistry.Store(key, p)
}

// 获取方法
func GetBizProcesser[IN any, OUT any](tag int) (BizProcesser[IN, OUT], bool) {
	key := registryKey{
		inType:  fmt.Sprintf("%T", *new(IN)),
		outType: fmt.Sprintf("%T", *new(OUT)),
		tag:     tag,
	}
	if val, ok := processerRegistry.Load(key); ok {
		return val.(BizProcesser[IN, OUT]), true
	}
	return nil, false
}
