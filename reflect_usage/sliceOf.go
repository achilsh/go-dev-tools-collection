package main

import (
	"fmt"
	"reflect"
)

func sliceOfDemo() {
	//  通过反射创建和操作 slice
	// 创建slice 类型：
	sliceType := reflect.SliceOf(reflect.TypeOf(1))

	//  创建该类型的 slice 对象
	sliceObj := reflect.MakeSlice(sliceType, 10, 10)

	fmt.Println("sliceObj cap: ", sliceObj.Cap(), ", len: ", sliceObj.Len()) //sliceObj.Len()

	// 给创建的slice 对象赋值：
	for i := 0; i < sliceObj.Len(); i++ {
		// 设置每一个元素的内容
		sliceObj.Index(i).Set(reflect.ValueOf(i))
	}
	//答应slice 对象：

	fmt.Println("slice obj: ", sliceObj.Interface())
}
