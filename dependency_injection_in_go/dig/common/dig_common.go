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
}

// NewContainer 创建一个容器
func NewContainer() *dig.Container {
	c := dig.New()

	//
	c.Provide(NewGObj)
	c.Provide(NewCObj)
	//
	c.Provide(NewHObj)
	c.Provide(NewDObj)
	//
	c.Provide(NewAObj)
	//
	c.Provide(NewFObj)
	c.Provide(NewEObj)
	c.Provide(NewBObj)
	//

	return c
}

func Run(d *dig.Container) {
	// 先实例化 AObj 入参，然后再运行 函数体，说明在运行函数之前，函数入参AObj 已经创建好了。
	// AObj 的对象创建过程是 根据他自身的依赖想来创建，依赖关系 是通过  Provide() 提供注册的。
	// 依赖项的创建 不保证按代码 Provide()的调用顺序。
	// 只有运行了 Invoke() 才会触发 依赖的实例化，不信就把这样注释掉，观察 上面对象的创建日志是否有输出。
	d.Invoke(func(a *AObj) {
		a.CallDemo()
	})
}
