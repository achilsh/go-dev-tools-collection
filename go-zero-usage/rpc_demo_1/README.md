## 使用 goctl 创建一个简单的服务

* 使用 goctl rpc -o pb/simple_demo.proto 创建一个 pb 文件

* 在这基础上编写 rpc 分组， 通过不同 service 名称来区分 不同的rpc分组

* 不同proto 文件之间的引用实例：
按下面格式来创建不同proto文件，然后运行命令：

```
goctl rpc protoc pb/greet/greet.proto --go_out=../ --go-grpc_out=../ --zrpc_out=. --client=true  -m
```

被import 的proto 源码生成命令(需要单独编译源码)：

```
protoc pb/base/base.proto --go_out=.. --go-grpc_out=..
```

生成 rpc_demo_1 目录下的源码文件：

```
├── base.pb.go
│   │   └── base.proto
│   └── greet
│       ├── greet_grpc.pb.go
│       ├── greet.pb.go
│       └── greet.proto
```

其中 greet.proto 文件内容：

```
// 声明 proto 语法版本，固定值
syntax = "proto3";

// proto 包名
package greet;

// 生成 golang 代码后的包名
option go_package = "rpc_demo_1/pb/greet";

import "pb/base/base.proto";

enum Status{
  UNSPECIFIED = 0;
  SUCCESS = 1;
  FAILED = 2;
}

message SendMessageReq{
  string message = 1;
}

message SendMessage{
  // 使用枚举
  Status status = 1;
  // 数组
  repeated string array = 2;
  // map
  map<string,int32> map = 3;
  // 布尔类型
  bool boolean = 4;
  // 序列号保留
  reserved 5;
}

message SendMessageResp{
//其中base 是包名；就是在 base.proto中定义的package base
  base.Base base = 1;
  SendMessage data = 2;
}

// 定义 Greet 服务
service Greet {
  // 定义客户端流式 rpc
  rpc SendMessage(stream SendMessageReq) returns (SendMessageResp);
}

// 定义 一个 Message service 实现和上面不同的服务分组

service message {
   rpc Pong(SendMessageReq) returns(SendMessageResp);
}
```

其中 base.proto 文件内容：

```
syntax = "proto3";

// proto 包名
package base;

// 生成 golang 代码后的包名
option go_package = "rpc_demo_1/pb/base";

message Base{
  int32 code = 1;
  string msg = 2;
}
```

* 生成类似如下源码 proto, pb 文件, 要求所有的生成源码放在项目 pb/gen目录下：

```
├── pb
│   ├── base
│   │   └── base.proto
│   ├── gen
│   │   ├── base.pb.go
│   │   ├── greet_grpc.pb.go
│   │   └── greet.pb.go
│   └── greet
│       └── greet.proto
```

如上面所示，把所有的pb 源码生成在 gen 目录 场景。
生成上面命令：

```
goctl rpc protoc pb/greet/greet.proto --go_out=../ --go-grpc_out=../ --zrpc_out=. --client=true  -m

生成 import proto的源码：
protoc pb/base/base.proto --go_out=.. --go-grpc_out=..
```

其中 greet.proto 内容：

```
// 声明 proto 语法版本，固定值
syntax = "proto3";

// proto 包名
package greet;

// 生成 golang 代码后的包名
option go_package = "rpc_demo_1/pb/gen";

// 从公共目录的根目录开始
import "pb/base/base.proto";

enum Status{
  UNSPECIFIED = 0;
  SUCCESS = 1;
  FAILED = 2;
}

message SendMessageReq{
  string message = 1;
}

message SendMessage{
  // 使用枚举
  Status status = 1;
  // 数组
  repeated string array = 2;
  // map
  map<string,int32> map = 3;
  // 布尔类型
  bool boolean = 4;
  // 序列号保留
  reserved 5;
}

message SendMessageResp{
//其中base 是包名；就是在 base.proto中定义的package base
  base.Base base = 1;
  SendMessage data = 2;
}

// 定义 Greet 服务
service Greet {
  // 定义客户端流式 rpc
  rpc SendMessage(stream SendMessageReq) returns (SendMessageResp);
}

// 定义 一个 Message service 实现和上面不同的服务分组

service message {
   rpc Pong(SendMessageReq) returns(SendMessageResp);
}
```

其中 base.proto 内容：

```
syntax = "proto3";

// proto 包名
package base;

// 生成 golang 代码后的包名
option go_package = "rpc_demo_1/pb/gen";

message Base{
  int32 code = 1;
  string msg = 2;
}

```
