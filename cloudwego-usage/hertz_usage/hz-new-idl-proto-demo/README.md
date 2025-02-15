* 指定 proto 协议来构建http 服务
* 定义proto文件，其中 proto/api.proto 是hz 提供的注解文件； 请在使用了注解的 proto 文件中import 该文件。
*  hz new --service hz-new-proto --mod hz-new-proto --idl proto/user/* -I proto 
* 如果更新了协议proto中的文件，使用命令来更新项目：  hz update -I proto --idl proto/user/*

#
* 参考： 
  1）https://github.com/cloudwego/hertz-examples/tree/main/hz/protobuf

  2）https://github.com/cloudwego/hertz-examples/tree/main/bizdemo/hertz_gorm_gen

* 修改逻辑， 增加接收请求的处理逻辑，可在目录 biz/handler/user/*.go中添加接收消息处理：