# 在go-0上使用 grpc 服务

## 手动快速构建

* 快速生成一个 proto 文件 命令：

```
goctl  rpc --o pb/demo.proto
```

上面命令生成了一个 demo.proto 文件。

* 根据上面一步的proto文件生成 pb 源码：

```
 protoc pb/demo.proto --go_out=./pb --go-grpc_out=./pb
```

* 根据 上面 proto 文件生成 zrpc 后端服务命令（单个rpc服务）：

```
goctl rpc protoc pb/demo.proto --go_out=./pb --go-grpc_out=./pb --zrpc_out=. --client=true 
```

* 根据上面proto 文件生成zrpc 后端服务命令（多个rpc服务）

```
goctl rpc protoc pb/demo.proto --go_out=./pb --go-grpc_out=./pb --zrpc_out=. --client=true  -m
```

*注意事项：
goctl rpc protoc 指令生成 rpc 服务对 proto 有一些事项须知：

proto 文件中如果有 import 语句，goctl 不会对 import 的 proto 文件进行处理，需要自行手动处理。
rpc service 中的请求体和响应体必须是当前 proto 文件中的 message，不能是 import 的 proto 文件中的 message。

--------------------------------------------

* 上面手动生成pb文件，生成源码不太方便，下面使用 命令： goctl rpc new 生成 rpc 服务。

## 自动构建 rpc服务

* 快速生成一个 rpc 服务； 接收一个终端参数来指定服务名称，比如 Helloworld：

```
 goctl rpc new helloworld
```

生成一个 Helloworld 目录，里面包含  rpc服务, 格式和 使用 proto文件生成的服务类似。

## 开启 grpc 调试开关

* gRPC 提供了调试功能，以便于我们可以通过 grpcurl 等工具进行调试， 在 go-zero，建议在开发环境和测试环境开启，预生产环境和正式环境建议关闭，因此我们在静态配置文件中将环境模式配置为 dev 或者 test 时才会开启（默认为 dev 环境）。在 config.yaml文件增加选项： Mode 和 值为 dev

```
Name: greet.rpc
ListenOn: 0.0.0.0:8080
Mode: dev

```

## rpc 服务端设置 中间件

* 框架提供的 中间件注册 是有配置决定，比如配置项：

``
先定义服务中间配置类型：
type (
    ServerMiddlewaresConf = internal.ServerMiddlewaresConf
)
ServerMiddlewaresConf struct {
  Trace      bool     `json:",default=true"`
  Recover    bool     `json:",default=true"`
  Stat       bool     `json:",default=true"`
  StatConf   StatConf `json:",optional"`
  Prometheus bool     `json:",default=true"`
  Breaker    bool     `json:",default=true"`
 }

在服务端配置定义使用上面配置类型：
type RpcServerConf struct {
    Middlewares ServerMiddlewaresConf
}

```

服务端依据上面配置定义，在服务端初始化的时候做判断是否需要注册中间件。比如初始化服务流程：
```

// NewServer returns a RpcServer.
func NewServer(c RpcServerConf, register internal.RegisterFn) (*RpcServer, error) {
 var err error
 if err = c.Validate(); err != nil {
  return nil, err
 }

 var server internal.Server
 metrics := stat.NewMetrics(c.ListenOn)
 serverOptions := []internal.ServerOption{
  internal.WithRpcHealth(c.Health),
 }

 if c.HasEtcd() {
  server, err = internal.NewRpcPubServer(c.Etcd, c.ListenOn, serverOptions...)
  if err != nil {
   return nil, err
  }
 } else {
  server = internal.NewRpcServer(c.ListenOn, serverOptions...)
 }

 server.SetName(c.Name)
 metrics.SetName(c.Name)
 setupStreamInterceptors(server, c)
 setupUnaryInterceptors(server, c, metrics)
 if err = setupAuthInterceptors(server, c); err != nil {
  return nil, err
 }

 //具体注册流程：

 func setupStreamInterceptors(svr internal.Server, c RpcServerConf) {
 if c.Middlewares.Trace {
  svr.AddStreamInterceptors(serverinterceptors.StreamTracingInterceptor)
 }
 if c.Middlewares.Recover {
  svr.AddStreamInterceptors(serverinterceptors.StreamRecoverInterceptor)
 }
 if c.Middlewares.Breaker {
  svr.AddStreamInterceptors(serverinterceptors.StreamBreakerInterceptor)
 }
}

func setupUnaryInterceptors(svr internal.Server, c RpcServerConf, metrics *stat.Metrics) {
 if c.Middlewares.Trace {
  svr.AddUnaryInterceptors(serverinterceptors.UnaryTracingInterceptor)
 }
 if c.Middlewares.Recover {
  svr.AddUnaryInterceptors(serverinterceptors.UnaryRecoverInterceptor)
 }
 if c.Middlewares.Stat {
  svr.AddUnaryInterceptors(serverinterceptors.UnaryStatInterceptor(metrics, c.Middlewares.StatConf))
 }
 if c.Middlewares.Prometheus {
  svr.AddUnaryInterceptors(serverinterceptors.UnaryPrometheusInterceptor)
 }
 if c.Middlewares.Breaker {
  svr.AddUnaryInterceptors(serverinterceptors.UnaryBreakerInterceptor)
 }
 if c.CpuThreshold > 0 {
  shedder := load.NewAdaptiveShedder(load.WithCpuThreshold(c.CpuThreshold))
  svr.AddUnaryInterceptors(serverinterceptors.UnarySheddingInterceptor(shedder, metrics))
 }
 if c.Timeout > 0 {
  svr.AddUnaryInterceptors(serverinterceptors.UnaryTimeoutInterceptor(
   time.Duration(c.Timeout)*time.Millisecond, c.MethodTimeouts...))
 }
}

func setupAuthInterceptors(svr internal.Server, c RpcServerConf) error {
 if !c.Auth {
  return nil
 }
 rds, err := redis.NewRedis(c.Redis.RedisConf)
 if err != nil {
  return err
 }

 authenticator, err := auth.NewAuthenticator(rds, c.Redis.Key, c.StrictControl)
 if err != nil {
  return err
 }

 svr.AddStreamInterceptors(serverinterceptors.StreamAuthorizeInterceptor(authenticator))
 svr.AddUnaryInterceptors(serverinterceptors.UnaryAuthorizeInterceptor(authenticator))

 return nil
}

关键流程是：定义 中间件 f， 调用 Server.AddUnaryInterceptors(f)
Server.AddStreamInterceptors()
如果业务注册自定义的中间件，需要的步骤：

1. 定义
2. 注册

```

