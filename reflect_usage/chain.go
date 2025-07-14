package main

import (
	"fmt"
	"reflect"
)

//解释如何使用reflect.ChanOf()创建 chan 类型， 使用 reflect.MakeChan()创建 chan 对象

func createChanType() {
	// 创建一个双向， int type 的 chan 类型
	rwChanType := reflect.ChanOf(reflect.BothDir, reflect.TypeOf(0))
	fmt.Println("rw chan TYpe: ", rwChanType)
	// 创建写chan, int type 的 chan 类型
	wChanType := reflect.ChanOf(reflect.SendDir, reflect.TypeOf(1))
	fmt.Println("w chan type: ", wChanType)
	// 创建读 chan, int type 的 chan 类型
	rChanType := reflect.ChanOf(reflect.RecvDir, reflect.TypeOf(1))
	fmt.Println("r chan type: ", rChanType)

	//
	creatChanItem()
}

func creatChanItem() {
	// 创建一个 双向的， int type 的 chan 类型
	rwChType := reflect.ChanOf(reflect.BothDir, reflect.TypeOf(0))
	//  创建 上面chan 类型对象：
	rwChValue := reflect.MakeChan(rwChType, 10)
	for i := 0; i < 5; i++ {
		//向里面发送数据
		v := reflect.New(reflect.TypeOf(0)).Elem()
		// 设置 要发送的数据值
		v.SetInt(int64(i))

		// 向通道发送数据
		rwChValue.Send(v)
	}
	rwChValue.Close()

	for {
		if rwChValue.IsNil() {
			break
		}

		// 通过通道来接收数据
		v, ok := rwChValue.Recv()
		if ok {
			fmt.Println("recv data: ", v.Interface())
			continue
		}
		break
	}
}
