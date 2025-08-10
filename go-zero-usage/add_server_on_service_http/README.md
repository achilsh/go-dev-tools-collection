## 给 http server api 文件的每个 service 增加

* 定义 user.api 文件
* 运行 命令 生成 http server:   goctl api go  -api user.api -dir ./

## 生成代码，服务端接收客户端请求处理逻辑

* 首先在 router.go 注册每个 url 的handler, 包括注册一些中间件，把url 进行分组等。 每个handler 内部处理是调用 logic 中的具体逻辑函数。

* 客户端请求路由需要加上 /v1/** 具体的路由信息

## 重新运行命令，文件覆盖情况

* 如果更新了 api文件，重新运行  goctl api go  -api user.api -dir ./  命令， 会覆盖  types/types.go 文件，router.go 文件；

## api 定义http协议的描述

* api 文件主要包括的内容：api 领域特性语言包含 “语法版本”，“info 块”，“结构体声明”，“服务描述” 等几大块语法组成，其中结构体和 Golang 结构体 语法几乎一样，只是移除了 struct 关键字。

* api 文件中的注释： 单行注释 //， 多行注释 /*
    this is demo test.
**/

* api 文件 版本声明语法：

```
syntax = "v1"
```

* info 语句，仅对api文件进行描述，不参与代码生成，在生成 的types.go文件中不会有info 语句。 格式：

```
info()
或者
info(
   k1: "value" 
   k2: "value"
)
```

* import 语句，在 api 中引入其他 api 文件的语法块，其支持相对/绝对路径。
