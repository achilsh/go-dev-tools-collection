package main

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
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

var wrapperedFunc reflect.Value

func init() {
	wrapperedFunc = wrapperFunc(toWrapperFunc)
}
func dynamicCallFunc() {

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

//  使用 reflect.FuncOf() 动态生成客户端函数。

type UserImpl struct {
	ID   int    `json:"id"`
	Name string `json:"name"`
}
type UserService interface {
	GetUser(id int) (UserImpl, error)
	CreateUser(u UserImpl) error
}

type RPCClient struct {
	url string
}

func NewRPCClient(url string) *RPCClient {
	return &RPCClient{
		url: url,
	}
}

// getMethodNameFromFunc 从函数指针获取方法名
func getMethodNameFromFunc(funcVal interface{}) string {
	// 简化实现，实际项目中可通过更复杂的反射获取方法名
	// 此处假设函数指针的字符串表示包含方法名
	funcStr := fmt.Sprintf("%T", funcVal)
	parts := strings.Split(funcStr, ".")
	return parts[len(parts)-1]
}

// getMethodTypes 获取方法的输入输出类型
func getMethodTypes(methodType reflect.Type) ([]reflect.Type, []reflect.Type) {
	inTypes := make([]reflect.Type, methodType.NumIn())
	for j := 0; j < methodType.NumIn(); j++ {
		inTypes[j] = methodType.In(j)
	}

	outTypes := make([]reflect.Type, methodType.NumOut())
	for j := 0; j < methodType.NumOut(); j++ {
		outTypes[j] = methodType.Out(j)
	}

	return inTypes, outTypes
}

// processRPCResult 处理RPC调用结果
func processRPCResult(err error, result []interface{}, outTypes []reflect.Type) []reflect.Value {
	if err != nil {
		errType := reflect.TypeOf((*error)(nil)).Elem()
		results := make([]reflect.Value, len(outTypes))
		for i, t := range outTypes {
			if t == errType {
				results[i] = reflect.ValueOf(err)
			} else {
				results[i] = reflect.Zero(t)
			}
		}
		return results
	}

	results := make([]reflect.Value, len(outTypes))
	for i, t := range outTypes {
		if i < len(result) {
			results[i] = reflect.ValueOf(result[i])
		} else {
			results[i] = reflect.Zero(t)
		}
	}
	return results
}

// createServiceStruct 创建含函数字段的结构体
func createServiceStruct(interfaceType reflect.Type) (reflect.Type, reflect.Value) {
	methodCount := interfaceType.NumMethod()
	serviceMethods := make([]reflect.StructField, methodCount)

	for i := 0; i < methodCount; i++ {
		method := interfaceType.Method(i)
		methodType := method.Type
		inTypes, outTypes := getMethodTypes(methodType)

		serviceMethods[i] = reflect.StructField{
			Name: method.Name,
			Type: reflect.FuncOf(inTypes, outTypes, false),
		}
	}

	structType := reflect.StructOf(serviceMethods)
	structValue := reflect.New(structType).Elem()
	return structType, structValue
}

// generateMethodFuncs 生成方法调用函数并赋值
func generateMethodFuncs(c *RPCClient, interfaceType reflect.Type, serviceStruct reflect.Value) error {
	for i := 0; i < interfaceType.NumMethod(); i++ {
		method := interfaceType.Method(i)
		methodType := method.Type
		inTypes, outTypes := getMethodTypes(methodType)

		funcType := reflect.FuncOf(inTypes, outTypes, false)
		if funcType == nil {
			return fmt.Errorf("无法生成函数类型 for method: %s", method.Name)
		}

		methodFunc := reflect.MakeFunc(funcType, func(args []reflect.Value) []reflect.Value {
			result, err := c.InvokeRPC(method.Name, args)
			return processRPCResult(err, result, outTypes)
		})

		serviceStruct.FieldByName(method.Name).Set(methodFunc)
	}
	return nil
}

// createInterfaceImplementation 通过闭包实现接口
func createInterfaceImplementation(interfaceType reflect.Type, serviceStruct reflect.Value) reflect.Value {
	interfaceImpl := reflect.New(interfaceType).Elem()

	// 缓存方法映射，避免重复查找
	methodMap := make(map[string]reflect.Value)
	for i := 0; i < interfaceType.NumMethod(); i++ {
		methodName := interfaceType.Method(i).Name
		methodMap[methodName] = serviceStruct.FieldByName(methodName)
	}

	// 为每个接口方法生成函数类型
	methodFuncs := make([]reflect.Value, interfaceType.NumMethod())
	for i := 0; i < interfaceType.NumMethod(); i++ {
		method := interfaceType.Method(i)
		methodType := method.Type
		inTypes, outTypes := getMethodTypes(methodType)
		funcType := reflect.FuncOf(inTypes, outTypes, false)

		methodFuncs[i] = reflect.MakeFunc(funcType, func(args []reflect.Value) []reflect.Value {
			// 从方法名获取函数字段
			methodName := method.Name
			methodField, ok := methodMap[methodName]
			if !ok {
				panic(fmt.Sprintf("方法 %s 未找到", methodName))
			}
			// 调用函数字段
			return methodField.Interface().(func([]reflect.Value) []reflect.Value)(args)
		})
	}
	interfaceImpl.Set(reflect.ValueOf(methodFuncs[0].Interface()).Convert(interfaceType))
	return interfaceImpl

	// interfaceImpl.Set(reflect.MakeFunc(interfaceType, func(args []reflect.Value) []reflect.Value {
	// 	// 获取方法名（第一个参数是函数指针，通过类型断言获取方法名）
	// 	methodPtr := args[0].Interface()
	// 	methodName := getMethodNameFromFunc(methodPtr)

	// 	// 从缓存中获取函数字段
	// 	methodField, ok := methodMap[methodName]
	// 	if !ok {
	// 		panic(fmt.Sprintf("方法 %s 未找到", methodName))
	// 	}

	// 	// 调用函数字段（去除self参数，传递实际参数）
	// 	return methodField.Interface().(func([]reflect.Value) []reflect.Value)(args[1:])
	// }))

	// return interfaceImpl
}

// 动态生成客户端实例代码， 入参是客户端访问服务端的抽象接口
func (c *RPCClient) CreateService(serviceType reflect.Type) any {
	// 检查是否为接口类型
	if serviceType.Kind() != reflect.Interface {
		return errors.New("服务类型必须是接口")
	}

	// 动态创建结构体类型，用于承载接口方法的实现
	_, serviceStruct := createServiceStruct(serviceType)

	// 为每个接口方法生成调用函数
	if err := generateMethodFuncs(c, serviceType, serviceStruct); err != nil {
		return err
	}

	// 通过闭包实现接口方法，创建接口实例
	interfaceImpl := createInterfaceImplementation(serviceType, serviceStruct)

	return interfaceImpl.Interface()

	// if serverInterfaceType.Kind() != reflect.Interface {
	// 	return errors.New("input serverType is not interface.")
	// }

	// // interfaceValue := reflect.New(serverInterfaceType).Elem()

	// // 动态创建结构体类型，用于承载接口方法的实现
	// serviceMethods := make([]reflect.StructField, serverInterfaceType.NumMethod())
	// for i := 0; i < serverInterfaceType.NumMethod(); i++ {
	// 	method := serverInterfaceType.Method(i)
	// 	methodType := method.Type // 获取方法的反射类型

	// 	inTypes := make([]reflect.Type, methodType.NumIn())
	// 	outTypes := make([]reflect.Type, methodType.NumOut())

	// 	for j := 0; j < methodType.NumIn(); j++ {
	// 		inTypes[j] = methodType.In(j) // 填充具体类型
	// 	}
	// 	for j := 0; j < methodType.NumOut(); j++ {
	// 		outTypes[j] = methodType.Out(j) // 填充具体类型
	// 	}

	// 	serviceMethods[i] = reflect.StructField{
	// 		Name: method.Name,
	// 		Type: reflect.FuncOf(inTypes, outTypes, false),
	// 	}
	// }

	// // 创建结构体类型
	// serviceStructType := reflect.StructOf(serviceMethods)
	// serviceStruct := reflect.New(serviceStructType).Elem()

	// for i := 0; i < serverInterfaceType.NumMethod(); i++ {
	// 	method := serverInterfaceType.Method(i)
	// 	methodType := method.Type

	// 	inTypes := make([]reflect.Type, methodType.NumIn())
	// 	for j := 0; j < len(inTypes); j++ {
	// 		inTypes[j] = methodType.In(j)
	// 	}

	// 	outTypes := make([]reflect.Type, methodType.NumOut())
	// 	// fmt.Println("outTypes len: ", len(outTypes))
	// 	// fmt.Println("methodType len: ", methodType.NumOut())

	// 	for j := 0; j < methodType.NumOut(); j++ {
	// 		// fmt.Println("j: ", j)
	// 		outTypes[j] = methodType.Out(j)
	// 	}
	// 	// fmt.Println("------")

	// 	funcType := reflect.FuncOf(inTypes, outTypes, false)
	// 	methodFunc := reflect.MakeFunc(funcType, func(args []reflect.Value) []reflect.Value {

	// 		result, err := c.InvokeRPC(method.Name, args)

	// 		if err != nil {
	// 			errType := reflect.TypeOf((*error)(nil)).Elem()
	// 			results := make([]reflect.Value, len(outTypes))

	// 			for i, t := range outTypes {
	// 				if t == errType {
	// 					results[i] = reflect.ValueOf(err)
	// 				} else {
	// 					results[i] = reflect.Zero(t)
	// 				}
	// 			}
	// 			return results
	// 		}

	// 		// ok
	// 		results := make([]reflect.Value, len(outTypes))
	// 		for i, t := range outTypes {
	// 			if i < len(result) {
	// 				results[i] = reflect.ValueOf(result[i])
	// 			} else {
	// 				results[i] = reflect.Zero(t)
	// 			}
	// 		}
	// 		return results
	// 	})

	// 	serviceStruct.FieldByName(method.Name).Set(methodFunc)
	// }

	// //
	// interfaceImpl := reflect.New(serverInterfaceType).Elem()

	// // 缓存方法映射，避免重复查找
	// methodMap := make(map[string]reflect.Value)
	// for i := 0; i < serverInterfaceType.NumMethod(); i++ {
	// 	methodName := serverInterfaceType.Method(i).Name
	// 	methodMap[methodName] = serviceStruct.FieldByName(methodName)
	// }

	// interfaceImpl.Set(reflect.MakeFunc(serverInterfaceType, func(args []reflect.Value) []reflect.Value {
	// 	// 获取方法名（第一个参数是函数指针，通过类型断言获取方法名）
	// 	methodPtr := args[0].Interface()
	// 	methodName := getMethodNameFromFunc(methodPtr)

	// 	// 从缓存中获取函数字段
	// 	methodField, ok := methodMap[methodName]
	// 	if !ok {
	// 		panic(fmt.Sprintf("方法 %s 未找到", methodName))
	// 	}

	// 	// 调用函数字段（去除self参数，传递实际参数）
	// 	return methodField.Interface().(func([]reflect.Value) []reflect.Value)(args[1:])
	// }))

	// return interfaceImpl.Interface()
	// 将结构体转换为接口实例
	// interfaceValue.Set(reflect.ValueOf(serviceStruct.Interface()))

	// return interfaceValue.Interface()
}

// 输入参数是 方法名，业务数据。
// 返回参数： 数据列表，错误码
func (c *RPCClient) InvokeRPC(methodName string, args []reflect.Value) ([]any, error) {
	url := fmt.Sprintf("%v/%v", c.url, methodName)

	body, err := json.Marshal(args)
	if err != nil {
		return nil, err
	}

	resp, err := http.Post(url, "application/json", bytes.NewBuffer(body))
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	// 解析响应
	var response struct {
		Success bool   `json:"success"`
		Data    []any  `json:"data"`
		Error   string `json:"error"`
	}
	if err := json.NewDecoder(resp.Body).Decode(&response); err != nil {
		return nil, err
	}

	if !response.Success {
		return nil, errors.New(response.Error)
	}

	return response.Data, nil
}

func clientCallServer() {
	client := NewRPCClient("http://127.0.0.1:8080/rpc")

	userService := client.CreateService(reflect.TypeOf((*UserService)(nil)).Elem())

	userServiceImpl := (userService).(UserService)
	err := userServiceImpl.CreateUser(UserImpl{
		ID:   1000,
		Name: "achilsh client call service",
	})
	if err != nil {
		fmt.Println("call create user fail, err: ", err)
		return
	}

	userInfo, err := userServiceImpl.GetUser(1000)
	if err != nil {
		fmt.Println("call get user fail, err: ", err)
		return
	}
	fmt.Println("get user info: ", userInfo)
}
