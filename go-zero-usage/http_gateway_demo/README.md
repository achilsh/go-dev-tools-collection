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
