## 演示如何创建 http-trpc 的后端服务

* 运行命令， 产生项目：  trpc create -p pb/helloworld.proto -m=http_trpc_demo -f  --nogomod  --alias -o ./

* 修改服务端的配置文件，支持 http ：

```
    - name: http.helloworld.Greeter  # Route name for the service.
      ip: 127.0.0.1  # Service listening IP address, can use placeholder ${ip}. Use either ip or nic, ip takes priority.
      # nic: eth0
      port: 8001  # Service listening port, can use placeholder ${port}.
      network: tcp  # Network listening type: tcp or udp.
      protocol: http  # Application layer protocol: trpc or http.
      timeout: 1000  # Maximum processing time for requests in milliseconds.

```

* 搭建 http-rpc 协议的文档参考： <https://github.com/trpc-group/trpc-go/blob/main/http/README.zh_CN.md>
