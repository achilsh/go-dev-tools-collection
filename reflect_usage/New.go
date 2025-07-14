package main

import (
	"fmt"
	"reflect"
)

func reflectNew() {
	// 返回一个指针类型的对象
	intPtr := reflect.New(reflect.TypeOf(0))
	fmt.Println(intPtr.Kind())

	// 获取指针所指的内容
	intType := intPtr.Elem()
	fmt.Println(intType.Kind())

	// 向指针指向的内容设置值：
	intType.Set(reflect.ValueOf(1000))
	fmt.Println(intType.Interface())
	//
	checkIsNilZero()
}

func checkIsNilZero() {
	var xx chan int
	vRf := reflect.ValueOf(xx)
	fmt.Println("is nil: ", vRf.IsNil())
	//
	fmt.Println("is zero: ", vRf.IsZero())
	fmt.Println("is valid: ", vRf.IsValid())
}
