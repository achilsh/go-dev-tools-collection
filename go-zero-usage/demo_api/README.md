## 生成 http server: api_server 服务

* goctl api new  server_name 开发者需要在终端指定 "服务名称参数"；输出目录为当前工作目录；快速生成 go http server.
* 自动生成 http 服务，运行命令：  goctl  api  new   api_server 就会生成 api_server 服务，目录列表：

```
├── etc
│   └── userapi.yaml     # 配置文件
├── internal
│   ├── config
│   │   └── config.go    # 配置定义
│   ├── handler
│   │   └── handler.go   # 处理器
│   ├── logic
│   │   └── logic.go     # 业务逻辑
│   ├── svc
│   │   └── service.go   # 服务上下文
│   └── types
│       └── types.go     # 类型定义
├── userapi.go           # 入口文件
└── userapi.api          # API 定义文件
```

* 根据 api 接口文档 生成 服务,参考其他项目， 首先要定义api 文件，然后使用命令生成：
**
