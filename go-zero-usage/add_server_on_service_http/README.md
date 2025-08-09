## 给 http server api 文件的每个 service 增加

* 定义 user.api 文件
* 运行 命令 生成 http server:   goctl api go  -api user.api -dir ./

## 生成代码，服务端接收客户端请求处理逻辑

* 首先在 router.go 注册每个 url 的handler, 包括注册一些中间件，把url 进行分组等。 每个handler 内部处理是调用 logic 中的具体逻辑函数。

* 客户端请求路由需要加上 /v1/** 具体的路由信息

## 重新运行命令，文件覆盖情况

* 如果更新了 api文件，重新运行  goctl api go  -api user.api -dir ./  命令， 会覆盖  types/types.go 文件，router.go 文件；
