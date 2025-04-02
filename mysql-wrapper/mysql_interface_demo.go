package main

import (
	"gorm.io/gorm"

	. "github.com/achilsh/go-dev-tools-collection/mysql-wrapper/lib"
)

type DemoDataModel struct {
	ID   int64  `json:"id" redis:"id" gorm:"column:id"` //业务id
	Age  int8   `json:"age" redis:"age" gorm:"column:age"`
	Name string `json:"name" redis:"name" gorm:"column:name"`
}

func (DemoDataModel) TableName() string {
	return "t_demo_data_info"
}

// 定义某个model的所有操作的总类型
type DemoDataModelOpType = DbOpsTplInterface[DemoDataModel]

// 声明一个全局具体的对象 引用
var demoDataModelOpHandle DemoDataModelOpType = nil
var gDb *gorm.DB = nil

func GetDB() *gorm.DB {
	return gDb
}

func init() {
	InitDemoDataOpHandle()
}

func main() {
	GetDemoDataModelHandle().Insert([]*DemoDataModel{
		&DemoDataModel{
			//
		},
	})
}

// InitDemoDataOpHandle； 创建一个具体model的操作 实例对象； 就可以在项目中使用该对象； 在主流程中注册该函数
func InitDemoDataOpHandle() DemoDataModelOpType {
	demoDataModelOpHandle = NewDBModelTplWrapper[DemoDataModel](GetDB()).SetTabName(DemoDataModel{}.TableName())
	return demoDataModelOpHandle
}
func GetDemoDataModelHandle() DemoDataModelOpType {
	return demoDataModelOpHandle
}
