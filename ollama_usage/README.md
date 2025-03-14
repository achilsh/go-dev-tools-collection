## ollama 介绍
### ollama 作用 和 类似框架
* 本地部署和运行大模型，方便个人开发。类似的框架有： langchain；langchaingo；llama.cpp 等。
*  

### ollama 主要功能和内部逻辑
* ollama 对外提供的通信方式， 对内和模型推理模块通信方式。
* ollama 本地加载模型流程，调度和运行模型机制。3



### 如何使用 ollama 
* 本地编译， go build  生成 ollama 二进制文件
* ollama 运行命令:
``` shell
Large language model runner

Usage:
  ollama [flags]
  ollama [command]

Available Commands:
  serve       Start ollama
  create      Create a model from a Modelfile
  show        Show information for a model
  run         Run a model
  stop        Stop a running model
  pull        Pull a model from a registry
  push        Push a model to a registry
  list        List models
  ps          List running models
  cp          Copy a model
  rm          Remove a model
  help        Help about any command

```