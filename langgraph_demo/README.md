* langgraphgo demo 关键流程：
* 1 源码分析

* 1.1 结构体描述

```
包括： 节点集合； 节点之间的边； 图入口的节点名称。

type MessageGraph struct {
 // nodes is a map of node names to their corresponding Node objects.
 nodes map[string]Node

 // edges is a slice of Edge objects representing the connections between nodes.
 edges []Edge

 // entryPoint is the name of the entry point node in the graph.
 entryPoint string
}

```

* 1.2 关键接口：使用给定名字和函数的节点添加到图中；其中函数参数消息数据，输出的是消息数据。

```
func (g *MessageGraph) AddNode(name string, fn func(ctx context.Context, state []llms.MessageContent) ([]llms.MessageContent, error)) 

```

* 1.3 在节点 from, to 节点之间增加一条边

```
func (g *MessageGraph) AddEdge(from, to string)
```

* 1.4 设置图 入口点节点的 名称

```
func (g *MessageGraph) SetEntryPoint(name string)
```

* 1.5 编译消息图，返回可运行的实例

```
func (g *MessageGraph) Compile() (*Runnable, error) 
```

* 1.6 使用给定输入消息 执行编译后的消息图；返回结果或中间的第一个出现的错误

```
func (r *Runnable) Invoke(ctx context.Context, messages []llms.MessageContent) ([]llms.MessageContent, error)
```
