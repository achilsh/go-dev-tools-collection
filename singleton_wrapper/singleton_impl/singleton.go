package singletonimpl

import (
	"reflect"
	"sync"
)

// SingletonObjCreateFuncFactory 该函数返回一个创建单例的函数；该单例类型是 模板参数T.
// 函数入参 fn 可以为 nil对象，表示不初始化 内部创建的 T 对象。
func SingletonObjCreateFuncFactory[T any](fn func(*T)) func() *T {
	o := new(sync.Once)
	var data *T = nil

	//多次调用该返回值，可以服用相同的上面变量。
	return func() *T {
		o.Do(func() {
			data = singletonConstructor(fn)
		})
		return data
	}
}

// SingletonConstructor 包装入参函数，创建一个对象并通过入参函数初始化该对象，返回该对象
func singletonConstructor[T any](fn func(*T)) *T {
	v := new(T)
	if fn == nil || reflect.ValueOf(fn).IsNil() {
		return v
	}
	fn(v)
	return v
}
