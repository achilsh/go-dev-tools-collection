## 使用 go-zero 代理 http--> http

* client 使用 get 方法访问： <http://172.31.60.55:8888/users/you> 其中 header 字段和value：
* x-abc: 1123123 实现 gateway 修改header字段的value 并转发给 upstream

## gateway url映射关系

```
Name: gateway-example # gateway name
Host: 0.0.0.0 # gateway host
Port: 8888 # gateway port
Upstreams: # upstreams
  - Name: userservice
    Http: # grpc upstream
        Target: 0.0.0.0:8080 # grpc target,the direct grpc server address,for only one node
        Prefix: /api/v1  # 可选 后端服务的路由前缀，后端服务的url 为该前缀 + 代理对外提供的url
        Timeout: 3000    # 单位为毫秒，默认值 3000
    Mappings: # routes mapping： map request url to  /api/v1
      - Method: GET  ## 代理对外提供的方法,也是后端服务访问的方法
        Path: /ping   ## 代理对外提供的url, 后端路由在该url 前面加上 prefix的字段值
      - Method: GET
        Path: /users_only
      - Method: GET
        Path: /users/:id
      - Method: POST
        Path: /users/post
```

## gateway 只需在配置文件配置 后端服务的路由信息，至于如何转发 不需要 编写 lb 和 路由重写 中间件，由go-zero特性分支提供

* 直接有 go-zero 特性分支提供能力， 启动 gateway 时的配置文件参考：

```
 go run ./ -f etc/gateway-http2http-rewrite-inner.yaml 
```

## gateway 通过 etcd 连接 后端 http server, 不需要在gateway的配置文件上指定 后端 http ip地址

* 需要在upstream的http 上加上 服务注册能力，比如在 upstream 系统 studentservice 服务上，增加
* 1. 在配置文件 upstreams-server/studentservice/etc 目录下的配置文件
   student-api-ectd.yaml  student-api-etcd_2.yaml
   分别增加向etcd注册的etcd配置：

  ```

## 向etcd 注册自己的身份 key

Etcd:
  Hosts:
    - 127.0.0.1:2379
  Key: student_info.http.1  ## 向etcd 注册自己的身份 key

  ```

  * 2. 修改 upstream 源码， 在upstreams-server/studentservice 目录下文件：
  student.go 中增加向etcd注册：
  ```

  import "github.com/zeromicro/zero-contrib/rest/registry/etcd"

  ```
  在main 函数中增加：
  ```

 // http 自身向 etcd 注册动作
 logx.Must(etcd.RegisterRest(c.Etcd, c.RestConf))

  ```

  * 3. 修改配置源码， 在upstreams-server/studentservice/internal/config/config.go 中修改源码：
  ```

type Config struct {
    rest.RestConf

    // 新增 http server 自身节点向etcd上报的etcd节点配置：
    Etcd discov.EtcdConf
}

  ```

  * 4. 修改 gateway 的配置文件：http_gateway_2_http_with_etcd/etc/gateway-http2http-rewrite-inner.yaml 
  增加如下配置：
  ```

      Etcd:  ## 访问 client target时，使用的etcd,从中获取服务端的节点信息
        Hosts:  ## 这是etcd 的节点
        - 127.0.0.1:2379
        Key: student_info.http.1   ## 对端服务在etcd上注册key，根目录

  ```
  注意，目前gateway 同时 支持 通过配置后端ip地址，和etcd 服务发现后端地址，如果两者都配置了，那么优先使用 配置ip地址。


## 本项目支持 gateway->http with etcd and lb on ip:port.
