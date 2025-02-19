* 工具作用：静态检查 nil 对象；
* 工具链接： https://github.com/uber-go/nilaway
* 工具安装： go install go.uber.org/nilaway/cmd/nilaway@latest
* 工具使用：  nilaway -pretty-print=true  -fix  -include-pkgs="nilaway-usage,nilaway-usage"   ./...   