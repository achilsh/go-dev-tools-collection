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

* import 语句，在 api 中引入其他 api 文件的语法块，其支持相对/绝对路径。比如 user.api 使用了 persion.api 文件，参考：

```
import "person.api"
```

其他示例：

```
// 单行 import
import "foo"
import "/path/to/file"

// import 组
import ()
import (
    "bar"
    "relative/to/file"
)
```

* 数据类型支持：
目前支持： 沿用了 Golang 的数据类型。比如，基础类型：

```

 "bool"    | "uint8"     | "uint16"     | "uint32" | "uint64"  |
 "int8"    | "int16"     | "int32"      | "int64"  | "float32" |
 "float64" | "complex64" | "complex128" | "string" | "int"     |
 "uint"    | "uintptr"   | "byte"       | "rune"   | "any"     | 

 map[基础类型]  基础类型 或者type struct 自定义类型
 []基础类型， 或者 [] type  自定义类型
 [n] 基础类型， 或者 type  自定义类型

参考： https://go-zero.dev/docs/tutorials/api/types
```

结构体分组： type() eg:

```
type (
    Int int
    Integer = int
    Bar {
        Foo int               `json:"foo"`
        Bar bool              `json:"bar"`
        Baz []string          `json:"baz"`
        Qux map[string]string `json:"qux"`
    }
)
```

* server 语句：@server 语句是对一个服务语句的 meta 信息描述，其对应特性包含但不限于：

```
jwt 开关
中间件
路由分组
路由前缀

```

格式是：@server()空内容， @server(//  )有内容，如下：

```
 @server(
    // jwt 声明
    // 如果 key 固定为 “jwt:”，则代表开启 jwt 鉴权声明
    // value 则为配置文件的结构体名称
    jwt: Auth



    // 路由前缀
    // 如果 key 固定为 “prefix:”
    // 则代表路由前缀声明，value 则为具体的路由前缀值，字符串中没让必须以 / 开头
    prefix: /v1



    // 路由分组
    // 如果 key 固定为 “group:”，则代表路由分组声明
    // value 则为具体分组名称，在 goctl生成代码后会根据此值进行文件夹分组
    group: Foo



    // 中间件
    // 如果 key 固定为 middleware:”，则代表中间件声明
    // value 则为具体中间件函数名称，在 goctl生成代码后会根据此值进生成对应的中间件函数；代码目录在解释器目录中。
    middleware: AuthInterceptor



    // 超时控制
    // 如果 key 固定为  timeout:”，则代表超时配置
    // value 则为具体中duration，在 goctl生成代码后会根据此值进生成对应的超时配置
    timeout: 3s
)
```

* 服务条目：服务条目（ServiceItemStmt）是对单个 HTTP 请求的描述，包括 @doc 语句，@handler 语句，路由语句信息，格式：

```
service user (
    @doc "登录"
    @handler login

)

```

其中：@handler 语句是对单个路由的 handler 信息控制，主要用于生成 golang http.HandleFunc 的实现转换方法。

* 请求参数，返回值 struct 的tag定义，头部header 字段定义参考： <https://go-zero.dev/docs/tutorials/api/parameter>

## api 同目录引用，跨目录引用

* 同目录引用：参考：

```
import "person.api"
import "struct_def.api"
然后直接使用struct 类型即可。 
```

*跨目录引用，父子目录：

```
import "./company/company_api.api"
直接使用子目录文件定义的类型
```

* 跨目录引用，同级父目录的不同子目录：

```
import "../student/student_api.api"
```

## 配置的使用

* 修改  etc/foo.yaml 配置文件
* 在internal/config/config.go文件增加配置结构体定义。
* conf 目前已经默认自动支持 key 大小写不敏感，例如对应如下的配置我们都可以解析出来.
* 参考： <https://go-zero.dev/docs/tutorials/go-zero/configuration/overview>
*其中配置规则也如上所示，示例：

```
type Config struct {
    Name string // 没有任何 tag，表示配置必填

    <!--  default 当前参数默认值  -->
    Port int64 `json:",default=8080"` // 如果配置中没有配置，将会初始成 8080

    <!-- 当前字段是可选参数，允许为零值(zero value)  -->
    Path string `json:",optional"` //当前字段是可选参数，允许为零值(zero value)

    <!-- options 当前参数仅可接收的枚举值； default 当前参数默认值 -->
    Mode       string `json:",default=pro,options=dev|test|rt|pre|pro"`  

    <!-- json tag 的后面加上 env=SERVER_NAME 的标签，conf 将会自动去加载对应的环境变量。 -->
    ServerName string `json:",env=SERVER_NAME"`
}
```
