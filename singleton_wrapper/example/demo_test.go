package example

import (
	"fmt"
	"testing"

	singletonimpl "github.com/achilsh/go-dev-tools-collection/singleton_wrapper/singleton_impl"
)

type typeTest1 struct {
	x int
	y float32
}

func (ty *typeTest1) F1(x int) {
	ty.x = x
	fmt.Println("v: ", ty.x)
}
func (ty *typeTest1) F2() {
	fmt.Println(ty.x)
}

func TestSingletonObjCreateFuncFactory(t *testing.T) {
	//创建一个获取单例的对象,后续只需要调用该对象函数就可以。
	//如果想创建另外的单例，需要再次调用，获得两次是不同的单例对象。
	var GetTypeTest1I = singletonimpl.SingletonObjCreateFuncFactory[typeTest1](nil)
	GetTypeTest1I().F1(100)
	GetTypeTest1I().F2()

	GetTypeTest1I().F1(200)
	GetTypeTest1I().F2()

	var GEtTypeTest1II = singletonimpl.SingletonObjCreateFuncFactory(func(v *typeTest1) {
		v.x = 10000
	})
	GEtTypeTest1II().F2()
}
