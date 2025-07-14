package main

import (
	"fmt"
	"reflect"
	"time"
)

type Column struct {
	Name   string //列名
	Type   string //数据类型： INT, VARCHAR DATATIME 等
	GoType reflect.Type
}

func getTableColumns() []Column {
	return []Column{
		{Name: "id", Type: "INT", GoType: reflect.TypeOf(int(0))},
		{Name: "username", Type: "VARCHAR", GoType: reflect.TypeOf("")},
		{Name: "age", Type: "INT", GoType: reflect.TypeOf(int(0))},
		{Name: "create_at", Type: "DATETIME", GoType: reflect.TypeOf(time.Time{})},
	}
}

func dbTypeToGoType(dbType string) (reflect.Type, error) {
	switch dbType {
	case "INT":
		return reflect.TypeOf(int(0)), nil
	case "VARCHAR", "TXT":
		return reflect.TypeOf(""), nil
	case "DATETIME", "TIMESTAMP":
		return reflect.TypeOf(time.Time{}), nil
	default:
		return nil, fmt.Errorf("not support type: %v", dbType)
	}
}

// 根据输入的列信息创建一个表示表结构的 struct 类型
func generateStructFromColumns(tabName string, cols []Column) (reflect.Type, error) {
	//创建 struct 中的多个 fields. 其中reflect.StructField 是struct 的一个field 结构体。
	structFields := make([]reflect.StructField, 0, len(cols))

	for _, col := range cols {
		fieldName := fmt.Sprintf("%c%s", col.Name[0]-32, col.Name[1:])
		fmt.Println("fieldName: ", fieldName)

		field := reflect.StructField{
			Name: fieldName,
			Type: col.GoType,
			Tag:  reflect.StructTag(fmt.Sprintf(`db:"%s"`, col.Name)),
		}

		structFields = append(structFields, field)
	}

	// 创建一个struct 的类型
	structType := reflect.StructOf(structFields)
	return structType, nil
}

// 把map中的数据 转到 struct 类型的中，返回 struct 示例对象
func mapToStruct(data map[string]any, structType reflect.Type) (any, error) {
	//  创建 struct 对象 示例
	structVal := reflect.New(structType).Elem()

	// 遍历struct 中的每个一 field. 填充他们
	for i := 0; i < structType.NumField(); i++ {
		fieldType := structType.Field(i)
		//  获取tag name
		dbColName := fieldType.Tag.Get("db")
		if dbColName == "" {
			continue
		}

		// 根据tag name从 data 中获取数据
		v, ok := data[dbColName]
		if !ok {
			continue
		}

		vReflect := reflect.ValueOf(v)
		// 类型不匹配
		if !vReflect.Type().AssignableTo(fieldType.Type) {
			continue
		}
		// 设置 struct field 的 value
		structVal.Field(i).Set(vReflect)
	}
	return structVal.Interface(), nil
}

func structOfDemo() {
	colums := []Column{
		{Name: "id", Type: "INT", GoType: reflect.TypeOf(int(0))},
		{Name: "username", Type: "VARCHAR", GoType: reflect.TypeOf("")},
		{Name: "age", Type: "INT", GoType: reflect.TypeOf(int(0))},
		{Name: "created_at", Type: "DATETIME", GoType: reflect.TypeOf(time.Time{})},
	}

	userStructType, err := generateStructFromColumns("user", colums)
	if err != nil {
		panic(err)
	}
	fmt.Printf("动态生成的结构体类型: %v\n", userStructType)

	// 3. 模拟数据库查询结果（实际中从 SQL 查询获取）
	queryResult := map[string]interface{}{
		"id":         1001,
		"username":   "alice",
		"age":        25,
		"created_at": time.Date(2023, time.January, 1, 0, 0, 0, 0, time.UTC),
	}

	// 4. 将查询结果转换为动态结构体实例
	user, err := mapToStruct(queryResult, userStructType)
	if err != nil {
		panic(err)
	}
	fmt.Printf("结构体实例: %+v\n", user)
}
