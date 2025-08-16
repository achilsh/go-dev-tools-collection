## 使用 goctl --style 来支持对 “文件和文件夹的命名风格进行格式化。”

```
--style gozero   表示生成的文件名和文件夹命名风格是小写
--style goZero   表示生成的文件名和文件夹命名风格是驼峰
--style go_zero  表示生成的文件名和文件夹命名风格是蛇形
```

## 实例演示

* 使用  goctl api new snake_case_style  --style go_zero  生成 蛇形 的文件名，文件夹

```
└── snake_case_style
    ├── etc
    │   └── snake_case_style-api.yaml
    ├── go.mod
    ├── internal
    │   ├── config
    │   │   └── config.go
    │   ├── handler
    │   │   ├── routes.go
    │   │   └── snake_case_style_handler.go
    │   ├── logic
    │   │   └── snake_case_style_logic.go
    │   ├── svc
    │   │   └── service_context.go
    │   └── types
    │       └── types.go
    ├── snake_case_style.api
    └── snake_case_style.go
```

* 使用  goctl api new low_case_style  --style gozero  生成 小写 的文件名，文件夹:

```
├── etc
│   └── lowcasestyle-api.yaml
├── go.mod
├── internal
│   ├── config
│   │   └── config.go
│   ├── handler
│   │   ├── lowcasestylehandler.go
│   │   └── routes.go
│   ├── logic
│   │   └── lowcasestylelogic.go
│   ├── svc
│   │   └── servicecontext.go
│   └── types
│       └── types.go
├── low_case_style.api
└── lowcasestyle.go
```

* 使用  goctl api new hump_style  --style goZero  生成 驼峰 的文件名，文件夹:

```
├── etc
│   └── humpStyle-Api.yaml
├── go.mod
├── hump_style.api
├── humpStyle-Api.go
└── internal
    ├── config
    │   └── config.go
    ├── handler
    │   ├── humpStyleHandler.go
    │   └── routes.go
    ├── logic
    │   └── humpStyleLogic.go
    ├── svc
    │   └── serviceContext.go
    └── types
        └── types.go
```
