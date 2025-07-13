package main

import (
	"encoding/json"
	"fmt"
	"reflect"
)

func MapOfDemo() {
	// key: string, value: int
	mType := reflect.MapOf(reflect.TypeOf(""), reflect.TypeOf(0))
	fmt.Println(mType)
	//
	//  创建实例：
	mValue := reflect.MakeMap(mType)
	fmt.Println("new map obj type: ", mValue.Type(), ", is nil: ", mValue.IsNil()) //返回的是 false.

	// 向实例中添加元素
	mValue.SetMapIndex(reflect.ValueOf("age"), reflect.ValueOf(100))

	//  从实例中获取元素:
	itemValue := mValue.MapIndex(reflect.ValueOf("age"))
	fmt.Println("item value: ", itemValue.Int(), ", type: ", itemValue.Type())

	// 检查记录是否存在
	exist := mValue.MapIndex(reflect.ValueOf("abc")).IsValid()
	fmt.Println("abc exist: ", exist)

	// 遍历记录：
	iter := mValue.MapRange()
	for iter.Next() {
		fmt.Println("map key: ", iter.Key(), " , map value: ", iter.Value())
	}

	// map key: string, value: any
	anyMap := reflect.MapOf(
		reflect.TypeOf(""),
		reflect.TypeOf((*any)(nil)).Elem(),
	)
	// 创建map 实例
	anyMapItem := reflect.MakeMap(anyMap)

	//设置不同类型的 value.
	anyMapItem.SetMapIndex(reflect.ValueOf("age"), reflect.ValueOf(12))
	anyMapItem.SetMapIndex(reflect.ValueOf("name"), reflect.ValueOf("achilsh"))

	// 获取map中的值
	ageV := anyMapItem.MapIndex(reflect.ValueOf("age")).Interface()
	fmt.Println("age: ", ageV)
	nameV := anyMapItem.MapIndex(reflect.ValueOf("name")).Interface()
	fmt.Println("name: ", nameV)

	// 测试 toMap 接口
	type Address struct {
		City  string `json:"city"`
		State string `json:"state"`
	}

	type Person struct {
		Name string  `json:"name"`
		Age  int     `json:"age"`
		Addr Address `json:"address"`
	}

	person := Person{
		Name: "Alice",
		Age:  30,
		Addr: Address{
			City:  "Beijing",
			State: "CN",
		},
	}
	personJson, _ := json.Marshal(person)
	fmt.Println("person json: ", string(personJson))

	personMap, err := ToMap(person)
	if err != nil {
		fmt.Println("person struct to map fail, err: ", err)
		return
	}
	fmt.Println("person struct to map, value: ", personMap)

	jsonMap, err := json.Marshal(personMap)
	if err != nil {
		fmt.Println("json marshal fail, err: ", err)
		return
	}
	fmt.Println("json marsh map ret: ", string(jsonMap))
}

// 把任何对象转化为 map
func ToMap(obj any) (map[string]any, error) {
	val := reflect.ValueOf(obj)
	if val.Kind() == reflect.Pointer {
		val = val.Elem()
	}
	if val.Kind() != reflect.Struct {
		return nil, fmt.Errorf("is not struct for obj.")
	}

	// 构造 map 类型
	mType := reflect.MapOf(
		reflect.TypeOf(""),
		reflect.TypeOf((*any)(nil)).Elem(),
	)

	// 创建map实例
	mItem := reflect.MakeMap(mType)
	for i := 0; i < val.NumField(); i++ {
		field := val.Field(i)
		fieldType := val.Type().Field(i)

		if fieldType.IsExported() == false {
			continue
		}

		// 只取json tag的记录
		jsonTag := fieldType.Tag.Get("json")
		if jsonTag == "" || jsonTag == "-" {
			continue
		}

		if field.Kind() == reflect.Struct {
			nestMap, err := ToMap(field.Interface())
			if err != nil {
				return nil, fmt.Errorf("to map fail")
			}

			mItem.SetMapIndex(reflect.ValueOf(jsonTag), reflect.ValueOf(nestMap))
			continue
		}

		mItem.SetMapIndex(reflect.ValueOf(jsonTag), field)
	}
	return mItem.Interface().(map[string]any), nil
}
