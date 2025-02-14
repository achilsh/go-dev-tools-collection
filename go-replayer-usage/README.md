* http 在客户端-服务端通信的录制和回放方式：
* https://github.com/buger/goreplay/wiki/Getting-Started 介绍如何使用工具来录制 http client发过来的请求。
* a：只要在运行服务端安装上面工具，下载源码，进入目录 goreplay/cmd/gor  运行命令 go build . 生成 gor 二进制文件。
* b：  运行工具录制命令，将结果存放到文件中，比如： sudo ./gor --input-raw :6666 --output-file=request.gor 其中 6666是服务器监听请求端口
* c：  如果要将结果输出到终端，运行命令: sudo ./gor --input-raw :6666 --output-stdout
* d：  将文件保存文件 移动到其他服务器上（比如测试环境中）进行流量重放:  sudo ./gor --input-file ./request_0.gor --output-http="http://localhost:6666" 其中 --output-http就是其他服务的访问节点。