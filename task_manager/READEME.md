## 任务及管理的封装

* 主要使用场景：
* 异步运行一个任务，调用方需要等待异步任务结果。此时需要任务来跟踪异步执行，包括任务投递，任务运行，任务结果通知等。每个异步都是一个任务。


## 主要接口：
* 业务线程： 新建一个任务：NewAsyncTAskWrapper(....) 
* 业务线程：把新建任务添加到任务管理器中， 并等待任务运行结果，获取返回值和错误信息： GetAsyncTaskMngInstance().SyncWait(key, task)
* 任务处理线程：接收任务（可以底层网络事件通知，或者从其他通道收到任务），处理业务逻辑，处理完 发起任务结果的通知： GetAsyncTaskMngInstance.NotifyDone(key, tmpResultTask)


## 使用该库的基础和场景：
* 该库主要是解决 业务发起一些任务，这些任务都在一个 线程池中统一被处理。业务不关心由多个线程来处理，
* 只要把数据发给任务处理集，等待任务运行结果通知即可。
* 要求业务 发起一个任务并得到结果后再继续 发起下一个任务，如果用户并行发起多个任务 比如
```
var wg sync.Waitgroup
wg.Add(2)
 go func(){
    defer wg.Done()
    GetAsyncTaskMngInstance().SyncWait(key1, task) 
 }

  go func(){
    defer wg.Done()
    GetAsyncTaskMngInstance().SyncWait(key2, task) 
 }

 wg.Wait()
```
需要保证 不同的 key1, key2 值和 key.IsAppend() 是 false。因为对相同key 多次处理是采用fifo模式。