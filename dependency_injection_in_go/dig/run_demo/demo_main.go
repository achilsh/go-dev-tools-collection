package main

import (
	"dig_usage_demo/common"
)

// 运行 一个dig 使用例子，依赖项的关系如 "实现依赖关系图.jpeg" 所示。
func main() {
	d := common.NewContainer()
	common.Run(d)
}
