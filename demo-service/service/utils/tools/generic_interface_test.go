package tools

import (
	"fmt"
	"strconv"
	"testing"
)

// 实现该接口泛型的实例：

// type demo1GenericInterfaceImpl struct {
// }

// func (d1 *demo1GenericInterfaceImpl) Handle(data int) (string, error) {
// 	return fmt.Sprintf("%v + data", data), nil
// }

// type demo2GenericInterfaceImpl struct {
// }

// 示例实现 1：int => string
type DemoIntToStr struct{}

func (d *DemoIntToStr) Handle(i int) (string, error) {
	return fmt.Sprintf("Received: %d", i), nil
}

// 示例实现 2：string => int
type DemoStrToInt struct{}

func (d *DemoStrToInt) Handle(s string) (int, error) {
	return strconv.Atoi(s)
}

// func (d1 *demo2GenericInterfaceImpl) Handle(data string) (int, error) {
// 	retData, err := strconv.ParseInt(data, 10, 64)
// 	if err != nil {
// 		return -1, err
// 	}
// 	return int(retData), nil
// }

// func TestGerericInterfaceImpl(t *testing.T) {
// 	var i2s = &demo1GenericInterfaceImpl{}
// 	// i2s.Handle(111)
// 	retData, _ := CallBizProcesser[int, string](i2s, 111)
// 	t.Logf("%s", retData)
// }

// 单元测试
func TestGenericRegistry(t *testing.T) {
	// 注册
	RegisterBizProcesser[int, string](1, &DemoIntToStr{})
	RegisterBizProcesser[string, int](2, &DemoStrToInt{})

	// 获取并使用
	if p1, ok := GetBizProcesser[int, string](1); ok {
		res, _ := p1.Handle(42)
		t.Log("int -> string:", res)
	} else {
		t.Error("processor 1 not found")
	}

	if p2, ok := GetBizProcesser[string, int](2); ok {
		res, _ := p2.Handle("123")
		t.Log("string -> int:", res)
	} else {
		t.Error("processor 2 not found")
	}
}
