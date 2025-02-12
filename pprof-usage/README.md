*  使用 github.com/google/gops/agent 检测go实现进程性能参数,比如cpu使用，内存使用，栈空间，trace能力
*  包括本地运行的和远程运行的。
*  如何使用？ 如果在源码中没有集成 gops的agent，运行gops 只能得到一些基础信息，不如不带任何参数运行 gops只得到当前运行有哪些go进程。而像 memstats， pprof-cpu， pprof-heap，trace ，stack等。
*  gops option pid
*  远程诊断方式： gops option ip:port 要求在远程的服务上集成agent,比如：
*   if err := agent.Listen(agent.Options{ \
		Addr:"0.0.0.0:6667", \
	 }); err != nil { \
        fmt.Println(err)\
    }

* 数据的比较： 参考 github.com/google/go-cmp 它可替换 reflect.DeepEqual
* google uuid的生成： https://github.com/google/uuid