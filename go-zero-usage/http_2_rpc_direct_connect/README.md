## 演示 一个http 和 多个 rpc 服务之间直连场景

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
