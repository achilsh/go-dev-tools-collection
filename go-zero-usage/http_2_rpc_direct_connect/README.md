## 演示 一个http 和 一个个 rpc 服务之间直连场景

### 快速新建一个 Http server

* 命令：goctl api  new http_demo_server

### 快速新建多个 rpc servers

* 命令： goctl rpc new  rpc_demo_server1

### 编辑 http server， 在逻辑处理中增加对rpc 接口调用

#### http 内部采用直连接 rpc 服务

* 在http 配置文件增加连接rpc 的client 配置：

```
#  直接连接
DirectConnetClientCfg:
  Target: dns:///127.0.0.1:8080
```

* 在http 源码配置增加 定义：
config.go 增加

```
type Config struct {
 //  http 服务端配置
 rest.RestConf

 //  rpc client 的配置定义
 DirectConnetClientCfg zrpc.RpcClientConf
}

```

在http server ctx 中增加 rpc client创建，在 servicecontext.go 文件中增加定义, 并引用 rpc client 的协议包：

```

import (
    pb "rpc_demo_server1/rpc_demo_server1"
)


type ServiceContext struct {
 Config config.Config
 // 增加客户端句柄信息
 DirectConnClient pb.RpcDemoServer1Client
}

func NewServiceContext(c config.Config) *ServiceContext {
 retCtx := &ServiceContext{
  Config: c,
 }

 //创建 rpc client
 retCtx.DirectConnClient = pb.NewRpcDemoServer1Client(zrpc.MustNewClient(c.DirectConnetClientCfg).Conn())
 return retCtx
}

```

* 在 http server logic 层增加访问rpc 调用，在 logic/httpdemoserverlogic.go 修改

```

 rpcClientReq := &pb.Request{
  Ping: fmt.Sprintf("rpc req msg: %v", req.Name),
 }

 rsp, err := l.svcCtx.DirectConnClient.Ping(context.Background(), rpcClientReq)
 if err != nil {
  l.Logger.Errorf("rpc client call fail, err: %v", err)
  return nil, err
 }

 if rsp == nil {
  l.Logger.Errorf("rpc client call response is nil ")
  return nil, fmt.Errorf("rpc client call response is nil")
 }
 resp = &types.Response{
  Message: fmt.Sprintf("rpc response data: %v", rsp.Pong),
 }

 return resp, nil

```

* 修改 rpc 的逻辑处理, 在源文件  pinglogic.go 增加逻辑：

```

func (l *PingLogic) Ping(in *rpc_demo_server1.Request) (*rpc_demo_server1.Response, error) {
 // todo: add your logic here and delete this line

 return &rpc_demo_server1.Response{
  Pong: "111111111111111",
 }, nil
}
```

* 在浏览器上输入命令：<http://172.31.60.55:8888/from/you>

## 演示一个http 和多个 rpc 服务之间直连

* 服用上面创建的http server
* 创建另外两个rpc server, 命令：

```
创建服务2： 
goctl rpc new  rpc_demo_server2

创建服务3，和服务2业务处理一致，只是启动端口不一样：
goctl rpc new  rpc_demo_server3
```

* 在http 服务的配置文件中增加 调用 rpc 的 client 配置：

```
# 新增配置：
# 直连多个节点：

MoreNodeConnClientCfg:
  # 使用 endpoint，只需要配置ip:port 即可
  Endpoints:
    - 127.0.0.1:8082
    - 127.0.0.1:8083

```

* 在conf 源码上新增多节点定义,修改文件 config/config.go：

```
type Config struct {
 //  http 服务端配置
 rest.RestConf

 //  rpc client 的配置定义
 DirectConnetClientCfg zrpc.RpcClientConf

 // 多个rpc 服务 client 配置定义
 MoreNodeConnClientCfg zrpc.RpcClientConf
}

```

* 创建和多个rpc 连接客户端, 修改文件： svc/servicecontext.go

```
type ServiceContext struct {
 Config config.Config
 // 增加客户端句柄信息
 DirectConnClient pb.RpcDemoServer1Client
 // 多服务端的节点直连
 MoreCliConnClients mp2.RpcDemoServer2Client
}

func NewServiceContext(c config.Config) *ServiceContext {
 retCtx := &ServiceContext{
  Config: c,
 }

 //创建 单一rpc client并初始化
 retCtx.DirectConnClient = pb.NewRpcDemoServer1Client(zrpc.MustNewClient(c.DirectConnetClientCfg).Conn())
 // 初始化 多个 rpc client
 retCtx.MoreCliConnClients = mp2.NewRpcDemoServer2Client(zrpc.MustNewClient(c.MoreNodeConnClientCfg).Conn())

 return retCtx
}
```

* 在 http 服务业务逻辑层调用多个 rpc后端服务,修改文件：httpdemoserverlogic.go

```

func (l *Http_demo_serverLogic) Http_demo_server(req *types.Request) (resp *types.Response, err error) {
 // todo: add your logic here and delete this line
 // call rpc clients.

 rpcClientReq := &pb.Request{
  Ping: fmt.Sprintf("rpc req msg: %v", req.Name),
 }

 rsp, err := l.svcCtx.DirectConnClient.Ping(context.Background(), rpcClientReq)
 if err != nil {
  l.Logger.Errorf("rpc client call fail, err: %v", err)
  return nil, err
 }

 if rsp == nil {
  l.Logger.Errorf("rpc client call response is nil ")
  return nil, fmt.Errorf("rpc client call response is nil")
 }
 resp = &types.Response{
  Message: fmt.Sprintf("rpc response data: %v", rsp.Pong),
 }

 r, e := l.call_more_rpc_services(req)
 if e != nil {
  return nil, e
 }

 resp.Message += ", " + r.Message
 return resp, nil
}

func (l *Http_demo_serverLogic) call_more_rpc_services(req *types.Request) (resp *types.Response, err error) {

 rpcClientReq := &mp2.Request{
  Ping: fmt.Sprintf("rpc req msg: %v", req.Name),
 }

 rsp, err := l.svcCtx.MoreCliConnClients.Ping(context.Background(), rpcClientReq)
 if err != nil {
  l.Logger.Errorf("rpc client call fail, err: %v", err)
  return nil, err
 }

 if rsp == nil {
  l.Logger.Errorf("rpc client call response is nil ")
  return nil, fmt.Errorf("rpc client call response is nil")
 }
 resp = &types.Response{
  Message: fmt.Sprintf("more rpc client: %v", rsp.Pong),
 }

 return resp, nil
}


```

* 修改rpc 端口，修改 rpc_demo_server2, rpc_demo_server3 etc/rpcdemoserver2.yaml 文件， 分别 为 8082， 8083
* 启动各个rpc 服务，http server

* 调用http 服务： <http://172.31.60.55:8888/from/you>
