package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"reflect"
	"runtime"
	"strings"
	"time"
)

// 通过反射 解释如何获取一个接口的类型

type DemoAInterface interface {
	DoCaller()
}

type AA1 struct {
}

func (a *AA1) DoCaller() {

}
func getInterfaceTypeByReflect() {
	//获取接口本身类型
	tType := reflect.TypeOf((*DemoAInterface)(nil)).Elem()
	fmt.Println("interface type: ", tType)

	//获取接口指向的变量类型
	var xx DemoAInterface = &AA1{}
	fmt.Println("interface dynamic type: ", reflect.TypeOf(xx))

	//检查某种类型是否实现了某个具体接口
	checkTypeImplementInterface()
}

func checkTypeImplementInterface() {
	xx := new(AA1)
	tType := reflect.TypeOf((*DemoAInterface)(nil)).Elem()
	implement := reflect.TypeOf(xx).Implements(tType)
	fmt.Println("AA1 implement interface DemoAInterface: ", implement)
}

type AI1Interface interface {
	DoCallOne(a int)
	DoCallTwo()
}

type AA2 struct {
}

func (a *AA2) DoCallOne(int) {}
func (a *AA2) DoCallTwo()    {}

func reflectTypeUsage() {
	typeMethod := reflect.TypeOf((*AI1Interface)(nil)).Elem()
	// 获取接口的第n个函数和方法
	fmt.Println(typeMethod.Method(0))
	fmt.Println(typeMethod.Method(1))

	//获取接口内的方法个数
	fmt.Println("method nums: ", typeMethod.NumMethod())
	//通过函数名找到对应的函数：
	findMethod, ok := typeMethod.MethodByName("DoCallOne")
	if ok {
		fmt.Println("find DoCallOne method: ", findMethod.Name)
	}

	// 获取 接口类型的名字
	fmt.Println("types name: ", typeMethod.Name())
	// 获取接口的包名：
	fmt.Println("pkg name for type: ", typeMethod.PkgPath())
	// 获取接口类型的 类型, 比如类型
	fmt.Println("type value: ", typeMethod.Kind())

	// 获取接口指向的类型
	var xx AI1Interface = &AA2{}
	fmt.Println("type elem type: ", reflect.TypeOf(xx).Elem())

	//获取函数类型的入参个数：
	innerFunc := func(int, int) int { return 0 }
	funcInParamNums := reflect.TypeOf(innerFunc).NumIn()
	fmt.Println("func input param nums: ", funcInParamNums)

	// 获取函数类型的出参个数：
	funcOutParamNums := reflect.TypeOf(innerFunc).NumOut()
	fmt.Println("func output param nums: ", funcOutParamNums)

	//  获取函数返回的参数信息
	for i := 0; i < funcOutParamNums; i++ {
		fmt.Println("func output param type: ", reflect.TypeOf(innerFunc).Out(i))
	}

	//获取函数每个输入参数类型：
	for i := 0; i < funcInParamNums; i++ {
		fmt.Println("func in param type: ", reflect.TypeOf(innerFunc).In(i))
	}

	// 通过reflect 操作 struct field 信息：
	opStructFieldByReflect()

	reflectFuncOf()

}

type FieldsOnStruct struct {
	A int
	B float32
	C string
}

func opStructFieldByReflect() {
	var args FieldsOnStruct = FieldsOnStruct{
		A: 12,
		B: 123.12,
		C: "this is struct op",
	}

	structType := reflect.TypeOf(args)
	// 获取结构体的成员个数：
	fmt.Println("struct field nums: ", structType.NumField())

	//  获取结构体每个成员信息：
	for i := 0; i < structType.NumField(); i++ {
		fieldItem := structType.Field(i)
		fmt.Println("structType field: i ", fieldItem)

	}
	fmt.Println("structTYpe field: ", structType.FieldByIndex([]int{1}))

	//通过名字获取 field
	AField, ok := structType.FieldByName("A")
	if ok {
		fmt.Println("exist A field: ", AField)
	}
	fmt.Println("----call func...")
	//通过名字函数 获取 field
	AAField, ok := structType.FieldByNameFunc(func(p string) bool {
		if p == "A" {
			return true
		}
		return false
	})
	if ok {
		fmt.Println("find by call name func: ", AAField)
	} else {
		fmt.Println("not find by name func")
	}
}

func reflectFuncOf() {
	// 允许在运行时根据参数和返回值类型定义函数。
	// 动态创建一个函数类型，并带有是输入参数类型和返回值类型。 返回的是类型，不能拿这类型调用.
	fType := reflect.FuncOf([]reflect.Type{reflect.TypeOf(0)}, []reflect.Type{reflect.TypeOf("")}, false)
	fmt.Println(fType.Kind())
	// 返回一个 fType 的 新函数，包装了 func。就是创建函数实例。
	funcDefine := reflect.MakeFunc(fType, func(args []reflect.Value) (results []reflect.Value) {
		x := args[0].Int()
		r := fmt.Sprintf("this is demo: %d", x)
		return []reflect.Value{
			reflect.ValueOf(r),
		}
	})

	ret := funcDefine.Call([]reflect.Value{reflect.ValueOf(1111)})
	for i := 0; i < len(ret); i++ {
		fmt.Println("ret: ", ret[i].String())
	}

	//  使用reflect.FuncOf() 实现动态json 序列化函数：把入参 struct 序列化成字符串
	type User struct {
		Name string `json:"name"`
		Age  int    `json:"age"`
	}
	// 入参为任何类型的元素
	inTypes := []reflect.Type{reflect.TypeOf((*any)(nil)).Elem()}
	//  出参为字符串和错误
	outTypes := []reflect.Type{reflect.TypeOf(""), reflect.TypeOf((*any)(nil)).Elem()}
	// 根据入参和出参定义函数
	funcType := reflect.FuncOf(inTypes, outTypes, false)
	// 创建定义的函数
	defineFunc := reflect.MakeFunc(funcType, func(args []reflect.Value) (results []reflect.Value) {
		obj := args[0].Interface()
		data, err := json.Marshal(obj)

		if err != nil {
			return []reflect.Value{
				reflect.ValueOf(string(data)),
				reflect.ValueOf(err),
			}

		} else {
			return []reflect.Value{
				reflect.ValueOf(string(data)),
				reflect.Zero(reflect.TypeOf((*error)(nil)).Elem()),
			}
		}

	})

	//  调用函数
	u := User{
		Name: "achilsh",
		Age:  120,
	}

	retCall := defineFunc.Call([]reflect.Value{reflect.ValueOf(u)})
	fmt.Println("ret str: ", retCall[0].String())
	fmt.Println("ret err: ", retCall[1].Interface())

	// 动态代理函数
	dynamicCallFunc()
}

func dynamicCallFunc() {
	wrapperedFunc := wrapperFunc(toWrapperFunc)
	{
		wrapperFuncRet := wrapperedFunc.Call([]reflect.Value{reflect.ValueOf(10), reflect.ValueOf(100)})
		fmt.Println(wrapperFuncRet[0], wrapperFuncRet[1])
	}

	// 调用另外数据
	{
		wrapperFuncRet1 := wrapperedFunc.Call([]reflect.Value{reflect.ValueOf(200), reflect.ValueOf(300)})
		fmt.Println(wrapperFuncRet1[0], wrapperFuncRet1[1])
	}
}

func wrapperFunc(inFunc any) reflect.Value {
	targetFunc := reflect.ValueOf(inFunc)
	targetFuncType := targetFunc.Type()

	// 复制原始函数的输入参数类型
	inTypes := make([]reflect.Type, targetFuncType.NumIn())
	for i := 0; i < len(inTypes); i++ {
		inTypes[i] = targetFuncType.In(i)
	}

	outTypes := make([]reflect.Type, targetFuncType.NumOut())
	for i := 0; i < len(outTypes); i++ {
		outTypes[i] = targetFuncType.Out(i)
	}

	callFuncName := ""
	lastFuncNames := strings.Split(runtime.FuncForPC(reflect.ValueOf(inFunc).Pointer()).Name(), ".")
	if len(lastFuncNames) > 0 {
		callFuncName = lastFuncNames[len(lastFuncNames)-1]
	}

	funcType := reflect.FuncOf(inTypes, outTypes, false)
	defineFunc := reflect.MakeFunc(funcType, func(args []reflect.Value) []reflect.Value {
		inputBuf := bytes.NewBuffer([]byte("input parameters: "))
		outBuf := bytes.NewBuffer([]byte(", out parameters: "))

		// 获取输入参数并格式化到buf,用于打印
		for i := 0; i < len(args); i++ {
			inputBuf.Write([]byte(fmt.Sprintf("%v ", args[i].Interface())))
		}

		startTm := time.Now()
		defer func() {
			fmt.Println("func: ", callFuncName, ", cost: ", time.Since(startTm), ",", inputBuf.String(), ", ", outBuf.String())
		}()

		result := targetFunc.Call(args)
		time.Sleep(1 * time.Second)

		//  获取输出参数并格式化到buf中，用于打印
		for i := 0; i < len(result); i++ {
			outBuf.Write([]byte(fmt.Sprintf("%v ", result[i].Interface())))
		}
		return result
	})

	return defineFunc
}
func toWrapperFunc(a, b int) (int, error) {
	return a + b, nil
}
