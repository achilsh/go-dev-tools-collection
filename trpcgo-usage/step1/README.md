## trpc 工具的安装

* trpc-cmdline ，eg: <https://github.com/trpc-group/trpc-cmdline>

```
First, add the following into your ~/.gitconfig:

[url "ssh://git@github.com/"]
    insteadOf = https://github.com/
Then run the following to install trpc-cmdline:

go install trpc.group/trpc-go/trpc-cmdline/trpc@latest
```

工具选项：
![alt text](image.png)

* 安装依赖：

```
trpc setup
```

* 使用 protobuffer 文件来生成项目，比如项目proto文件：

```
syntax = "proto3";
package helloworld;

option go_package = "github.com/some-repo/examples/helloworld";

// HelloRequest is hello request.
message HelloRequest {
  string msg = 1;
}

// HelloResponse is hello response.
message HelloResponse {
  string msg = 1;
}

// HelloWorldService handles hello request and echo message.
service HelloWorldService {
  // Hello says hello.
  rpc Hello(HelloRequest) returns(HelloResponse);
}
```

* 使用上面proto文件产生一个完整的项目：

```
其中 -p 指定 proto 文件， -o 指示项目生成的目录，默认是和 pb文件时同一个级别的目录。
trpc create -p pb/helloworld.proto -o helloworld 
或者使用下面命令， -d 选项指定 proto文件的路径， -f 强制覆盖输出
trpc create -p helloworld.proto -d pb -o helloworld
```

运行服务器端服务：

```
go run .
```

运行客户端：

```
go run ./cmd/client/main.go
```

* 上面是产生一个完整的项目，现在只产生一个rpc 源码，还是之前根据helloworld.proto 文件。

```
trpc create -d pb -o rpc --rpconly
```
