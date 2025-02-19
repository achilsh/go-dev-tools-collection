// / dig 库使用的实例调用封装
package common

import (
	"fmt"

	"go.uber.org/dig"
)

type FObj struct {
	A int
}

// NewFObj 创建一个独立对象
func NewFObj() *FObj {
	fmt.Println("New F obj.")
	return &FObj{
		A: 100,
	}
}
func (f *FObj) Call() {
	fmt.Println("call f: ", f.A)
}

/////////////////////////////

type EObj struct {
	FHandle *FObj
	E       float32
}

// NewEObj 创建 E 对象
func NewEObj(f *FObj) *EObj {
	fmt.Println("new E obj")
	return &EObj{
		FHandle: f,
		E:       1.123,
	}
}
func (e *EObj) CallE() {
	if e.FHandle != nil {
		e.FHandle.Call()
	}
	fmt.Println("call E: ", e.E)
}

// /////////////////////////////
type BObj struct {
	EHandle *EObj
	B       string
}

func (b *BObj) CallB() {
	if b.EHandle != nil {
		b.EHandle.CallE()
	}

	fmt.Println("B: ", b.B)
}

func NewBObj(e *EObj) *BObj {
	fmt.Println("new B obj")
	return &BObj{
		EHandle: e,
		B:       "BObj",
	}
}
func (b *BObj) Call() {
	if b.EHandle != nil {
		b.EHandle.CallE()
	}
	fmt.Println("call B: ", b.B)
}

// /
type GObj struct {
	FHandle *FObj
	G       bool
}

func NewGObj(f *FObj) *GObj {
	fmt.Println("new G obj")
	return &GObj{
		FHandle: f,
		G:       true,
	}
}
func (g *GObj) Call() {
	if g.FHandle != nil {
		g.FHandle.Call()
	}
	fmt.Println("G: ", g.G)
}

////

type CObj struct {
	GHandle *GObj
	C       float64
}

func NewCObj(g *GObj) *CObj {
	fmt.Println("new C obj")
	return &CObj{
		GHandle: g,
		C:       123.456,
	}
}
func (c *CObj) Call() {
	if c.GHandle != nil {
		c.GHandle.Call()
	}
	fmt.Println("Call C: ", c.C)
}

//////

type HObj struct {
	H string
}

func NewHObj() *HObj {
	fmt.Println("new H obj")
	return &HObj{
		H: "HObj",
	}
}
func (h *HObj) Call() {
	fmt.Println("H: ", h.H)
}

// ////
type DObj struct {
	HHandle *HObj
	D       int
}

func NewDObj(h *HObj) *DObj {
	fmt.Println("new D obj")
	return &DObj{
		HHandle: h,
		D:       1000,
	}
}
func (d *DObj) Call() {
	if d.HHandle != nil {
		d.HHandle.Call()
	}
	fmt.Println("D: ", d.D)

}

// ///
type AObj struct {
	BHandle *BObj
	CHandle *CObj
	DHandle *DObj
}

func NewAObj(b *BObj, c *CObj, d *DObj) *AObj {
	fmt.Println("new A obj")
	return &AObj{
		BHandle: b,
		CHandle: c,
		DHandle: d,
	}
}
func (a *AObj) CallDemo() {
	if a.BHandle != nil {
		a.BHandle.CallB()
		fmt.Println()
	}
	if a.CHandle != nil {
		a.CHandle.Call()
		fmt.Println()
	}
	if a.DHandle != nil {
		a.DHandle.Call()
	}
	fmt.Println("call AObj demo.\n")

}

type InParamOne struct {
	A int
}

func NewInParamOne() *InParamOne {
	return &InParamOne{
		A: 1,
	}
}

type InParamTwo struct {
	B int
}

func NewInParamTwo() *InParamTwo {
	return &InParamTwo{
		B: 2,
	}
}

type InParamThree struct {
	C int
}

func NewInParamThree() *InParamThree {
	return &InParamThree{
		C: 3,
	}
}

type NODepend struct {
	NoItem int
}

// DigInData 将多个散装的依赖项（创建对象的函数入参），逻辑上组合在一起；
// 实际上每种类型是可以独立作为依赖项。
type DigInData struct {
	dig.In
	// 说明下面这些 是 可依赖的项； 作为参数的入参
	//等效与下面三个同时作为构造另外一个对象的函数入参。
	BIn    *InParamTwo
	AIn    *InParamOne
	CIn    *InParamThree
	NoItem *NODepend `optional:"true"` //如果没有通过 Provide(func)中的func 提供 *NoDepend 类型创建函数，也可以。
}

type OutRetOne struct {
	OA int
}
type OutRetTwo struct {
	OB int
}

type OutRetThree struct {
	OC int
}

type OutNoDepend struct {
	Item *NODepend // is option ret output.
}

// DigOutData 将散装的多个结果作为逻辑的整体，作为对象创建函数的返回值 一同返回。
// 实际上每个返回值对象都是独立的，
type DigOutData struct {
	dig.Out
	//说明下面几个类型是 函数返回的对象。下面等效于 同时返回 下面三个类型
	A1          *OutRetOne
	B1          *OutRetTwo
	C1          *OutRetThree
	SpecialItem *OutNoDepend
}

// 演示使用命名依赖，解决同一类型多个实例问题：
type NamedDepend struct {
	NameItem string
}

// 通过 provide 返回的多个依赖项；这些值是通过 provide 创建，
type OUtUseTwoItemInSameType struct {
	dig.Out
	Item1 *NamedDepend `name:"one"`
	Item2 *NamedDepend `name:"two"`
}

// NewMoreItemInSameType 是给provide调用的，生成多个相同类型的实例。
// 使用这些相同类型的不同对象。如果一些其他的依赖项，想通过 已被创建的这些实例来创建，那么需要把这些创建实例
// 作为 provide(func)的 func()入参传入。或者把这些相同类型的不同对象加入到 dig.In的 结构体中，使用 tag： `name:"xxx"`
func NewMoreItemInSameType() OUtUseTwoItemInSameType {
	ret := OUtUseTwoItemInSameType{}
	ret.Item1 = &NamedDepend{
		NameItem: "one",
	}
	fmt.Println("new instance one for OUtUseTwoItemInOneType， value: ", ret.Item1.NameItem)

	ret.Item2 = &NamedDepend{
		NameItem: "two",
	}
	fmt.Println("new instance two for OUtUseTwoItemInOneType， value: ", ret.Item1.NameItem)
	return ret
}

// 同时使用同种类型的多个实例，作为依赖项
type InUseTwoItemInSameType struct {
	dig.In

	ItemOne *NamedDepend `name:"one"` //
	ItemTwo *NamedDepend `name:"two"` //
}
type BizUseTwoItemInSameType struct {
	Item1 *NamedDepend
	Item2 *NamedDepend
}

// 仅使用同中类型的多个实例的一个例子。
type InUseOneTypeItem struct {
	dig.In
	ItemOne *NamedDepend `name:"one"`
}
type BizUseOneItemInSameType struct {
	Item1 *NamedDepend
}

// 定义数据结构实现 provide(func)中func返回值是接口和func的入参是接口场景：
type CallServer interface {
	SetAge(age int)
	GetAge()int 
}
//
type XiaoMingCallServer struct {
	Age int
}
func(s *XiaoMingCallServer) SetAge(age int) {
	s.Age = age
}
func (s *XiaoMingCallServer)GetAge()int {
	return s.Age
}
func NewXiaoMingCallServer() CallServer  {
	return &XiaoMingCallServer{}
}
// 定义一个包含抽象接口的结构体类型.
type WrapCallServer struct {
	item CallServer
}
// 定义结构体创建函数,该函数用于注册到 provoide()入参数中.
func SetCallServer(s CallServer) *WrapCallServer {
	return &WrapCallServer{
		item:s,
	}
}

// NewContainer 创建一个容器
func NewContainer() *dig.Container {
	c := dig.New()

	// 生成一个依赖项，类型为接口类型。
	c.Provide(NewXiaoMingCallServer)
	//注册一个依赖于接口的结构体 的创建函数
	c.Provide(SetCallServer)

	////type 1:
	// c.Provide(NewMoreItemInSameType)

	//或者 使用 c.Provide(func,dig.Name("one")); c.Provide(func,dig.Name("two")); 分别创建两个实例。
	//比如下面的：
	// 同个类型多个实例单独创建：
	//typ 2:
	c.Provide(func() *NamedDepend {
		return &NamedDepend{
			NameItem: "one",
		}
	}, dig.Name("one"))
	c.Provide(func() *NamedDepend {
		return &NamedDepend{
			NameItem: "two",
		}
	}, dig.Name("two"))

	////使用同个类型的多个实例变量。或者单独使用入参比如下面：
	c.Provide(func(in InUseTwoItemInSameType) *BizUseTwoItemInSameType {
		ret := &BizUseTwoItemInSameType{
			Item1: in.ItemOne,
			Item2: in.ItemTwo,
		}
		return ret
	})

	//
	c.Provide(func(in InUseOneTypeItem) *BizUseOneItemInSameType {
		ret := &BizUseOneItemInSameType{
			Item1: in.ItemOne,
		}
		return ret
	})

	c.Provide(NewFObj)
	c.Provide(NewEObj)
	c.Provide(NewBObj)

	c.Provide(NewGObj)
	c.Provide(NewCObj)
	//
	c.Provide(NewHObj)
	c.Provide(NewDObj)

	c.Provide(NewHObj)
	c.Provide(NewDObj)
	//
	c.Provide(NewAObj)

	// //////////// 注入各个独立对象的创建函数
	c.Provide(NewInParamOne)
	c.Provide(NewInParamTwo)
	c.Provide(NewInParamThree)
	// 注入 将多个对	//SpecialItem: a.NoItem,象作为参数一起的依赖项 生成 多个返回值的 函数
	c.Provide(func(a DigInData) DigOutData {
		var ret DigOutData = DigOutData{
			A1: &OutRetOne{
				OA: a.AIn.A * 10,
			},
			B1: &OutRetTwo{
				OB: a.BIn.B * 100,
			},
			C1: &OutRetThree{
				OC: a.CIn.C * 1000,
			},

			SpecialItem: &OutNoDepend{
				Item: a.NoItem,
			},
		}
		return ret
	})
	return c
}

// Run: 触发依赖项(Invokde(func)中func函数的入参)递归创建. 然后在运行 func()的函数体.
func Run(d *dig.Container) {
	// 先实例化 AObj 入参，然后再运行 函数体，说明在运行函数之前，函数入参AObj 已经创建好了。
	// AObj 的对象创建过程是 根据他自身的依赖想来创建，依赖关系 是通过  Provide() 提供注册的。
	// 依赖项的创建 不保证按代码 Provide()的调用顺序。
	// 只有运行了 Invoke() 才会触发 依赖的实例化，不信就把这样注释掉，观察 上面对象的创建日志是否有输出。
	d.Invoke(func(a *AObj) {
		a.CallDemo()
	})

	// 可以多次调用，以实现不同注入的运行。
	d.Invoke(func(demo *InParamOne) {
		fmt.Println("new input param init data: ", demo.A)
	})
	d.Invoke(func(demo *OutRetThree) {
		fmt.Println("new three last data: ", demo.OC)
	})
	d.Invoke(func(demo *OutRetOne) {
		fmt.Println("new one last data: ", demo.OA)
	})
	d.Invoke(func(demo *OutRetTwo) {
		fmt.Println("new two last data: ", demo.OB)
	})
	d.Invoke(func(demo *OutNoDepend) {
		if demo.Item == nil {
			fmt.Println("not impl NoDepend obj create func.")
		} else {
			fmt.Println("NoDepend obj value: ", demo.Item.NoItem)
		}
	})
	// 使用同个类型的多个实例变量。
	d.Invoke(func(demo *BizUseTwoItemInSameType) {
		fmt.Printf("item1 addr: %p, value: %v\n", demo.Item1, demo.Item1.NameItem)
		fmt.Printf("item2 addr: %p, value: %v\n", demo.Item2, demo.Item2.NameItem)
	})
	//单独使用一种类型多实例下中一个实例：
	d.Invoke(func(demo *BizUseOneItemInSameType) {
		fmt.Printf(" use one item1 addr: %p, value: %v\n", demo.Item1, demo.Item1.NameItem)
	})

	// //使用接口作为依赖项,如果想替换不同的接口实例化对象，只需要修改 Provoide(func)中的func
	// d.Invoke(func(demo CallServer) {
	// 	demo.SetAge(100)
	// 	fmt.Println("call server age: ", demo.GetAge())
	// })

	// // 触发依赖项的递归创建.
	// d.Invoke(func(demo *WrapCallServer){
	// 	demo.item.SetAge(300)
	// 	fmt.Println("wrapper interface data: ", demo.item.GetAge())
	// })
	d.Invoke(func(d1 CallServer, d2 *WrapCallServer) {
		d1.SetAge(400)
		fmt.Println("call server data: ", d1.GetAge())
		d2.item.SetAge(500)
		fmt.Println("wrapper interface data: ", d2.item.GetAge())
		//
		fmt.Println("again call server data: ", d1.GetAge())
	})
}
