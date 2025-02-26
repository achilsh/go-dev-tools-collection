* myqtt 详情描述：https://evxulacfa8.feishu.cn/docx/YSCUd2jhVoXF4qxgrUpcZi9Sn8e?from=from_copylink
  
* pub client 实例关键流程：
* 1 ）创建client option, 比如： opts := mqtt.NewClientOptions()
* 2 ）增加 broker 地址： opts.AddBroker(SERVERADDRESS)
* 3 ) 设置 client的ID:  opts.SetClientID(CLIENTID)
* 4 ) 设置属性： opts.SetOrderMatters(false) // 除非是要按有序投递，否则将该值设置为false,以便允许乱序消息.
* 5 ) 设置属性： opts.ConnectTimeout = time.Second  // 设置连接超时时间
* 6 ) 设置属性： opts.WriteTimeout = time.Second  // 设置写数据超时时间
* 7 ) 设置属性： opts.KeepAlive = 10    // 设置网络断联检测时间
* 8 ) 设置属性: opts.PingTimeout = time.Second // 设置本地心跳
* 9 ) 设置属性： opts.ConnectRetry = true // 设置主动重连
* 10 )设置属性： opts.AutoReconnect = true // 设置自动重连
* 11 ) 连接日志打印 callback 
    ```
        // Log events
        opts.OnConnectionLost = func(cl mqtt.Client, err error) {
            fmt.Println("connection lost")
        }
        opts.OnConnect = func(mqtt.Client) {
            fmt.Println("connection established")
        }
        opts.OnReconnecting = func(mqtt.Client, *mqtt.ClientOptions) {
            fmt.Println("attempting to reconnect")
        }

    ```
* 12 ) 使用上面的option 创建 pub client:
  ```
    client := mqtt.NewClient(opts)
  ```
  
* 13 ) client 连接 mqtt broker:
    ```
     if token := client.Connect(); token.Wait() && token.Error() != nil {
		panic(token.Error())
	}
    ```

* 14 ) 使用client 发布消息：
    ```
        t := client.Publish(TOPIC, QOS, false, msg)

        go func() {
            _ = t.Wait() // Can also use '<-t.Done()' in releases > 1.2.0
            if t.Error() != nil {
                fmt.Printf("ERROR PUBLISHING: %s\n", err)
            }
		}()
    ```

* 15 ）关闭连接：
    ```
     client.Disconnect(1000)
    ```

* #
#
* sub client 实例关键流程：
* 创建流程大部分和前面保持一致。比如:
* 1）客户端options对象创建：	opts := mqtt.NewClientOptions()
* 2）增加连接broker节点：opts := mqtt.NewClientOptions()
* 3）设置客户端id: opts.SetClientID(CLIENTID)
* 4）设置属性：opts.SetOrderMatters(false)  // 允许乱序消息
* 5）设置连接超时时间：opts.ConnectTimeout = time.Second
* 6）设置写超时时间： opts.WriteTimeout = time.Second
* 7）设置keepalive时间： opts.KeepAlive = 10
* 8）设置 本代理 心跳ping/pong时间： opts.PingTimeout = time.Second
* 9）设置自动重连：opts.ConnectRetry = true
* 10）设置自动重连：opts.AutoReconnect = true
* 11 设置 收到没有订阅topic的消息处理函数：
  ```
    opts.DefaultPublishHandler = func(_ mqtt.Client, msg mqtt.Message) {
        fmt.Printf("UNEXPECTED MESSAGE: %s\n", msg)
    }
  ```
* 12）打印连接断开函数：
  ```
    opts.OnConnectionLost = func(_ mqtt.Client, msg mqtt.Message) {
        fmt.Println("connection lost")
    }
  ```
* 13） 设置连接成功时的处理逻辑：主要是注册topic订阅函数
    ```
        opts.OnConnect = func(cli mqtt.Client){
            fmt.Println("clientid is connected.")

            //client 连接上后，首先是注册 所有 topic和对应的处理函数，比如：
            t := c.Subscribe(TOPIC, QOS, func(client paho.Client, msg paho.Message) {
                // 收到订阅数据的处理逻辑； 包括处理完后对mqtt 发送回包消息（往另外 reply_topic 上发送msg）
            })
        }
    ```
    或者设置连接上的函数处理， 比如
    ```
    type SubscribleRoute struct {  // 定义消息订阅topic和具体的业务逻辑处理
        TopicName string
        BizLogicHandle  func() error
    }
    var logicHandlers = []SubscribleRoute {
        {"aaa", func(msg Message) ([]byte, error) },
        {"bbb", func(msg Message) ([]byte, error) },
    }
    // 注册 sub client 连接上后的回调函数; 主要是注册订阅topic 的回调函数
    opts.SetOnConnectHandler(func(c mqtt.Client) {
        ///
        fmt.Println("clientid is connected.")
        // client 连接上后，首先注册 所有 topic和对应的处理函数，比如：
        for i := 0; i < len(logicHandlers); i++ {
            t := c.Subscribe(logicHandlers[i].TopicName, QOS, func(client paho.Client, msg paho.Message) { 
                // 如果 options.OrderMatters is true； Subscrible内注册函数必须不能阻塞； 就是下面不能阻塞
                // 具体的业务逻辑
                logicHandlers[i].BizLogicHandle()
                // 收到订阅数据的处理逻辑； 包括处理完后对mqtt 发送回包消息（往另外 reply_topic 上发送msg）

            go func() {
                _ = t.Wait() // Can also use '<-t.Done()' in releases > 1.2.0
                if t.Error() != nil {
                    fmt.Printf("ERROR SUBSCRIBING: %s\n", t.Error())
                } else {
                    fmt.Println("subscribed to: ", TOPIC)
                }
		    }()
            })
        }
    })
    ```

* 14 ) 设置连接中的函数处理：
    ```
    opts.OnReconnecting = func(mqtt.Client, *mqtt.ClientOptions) {
		fmt.Println("attempting to reconnect")
	}
    ```
* 15 ) 使用上面的Opt创建客户端: client := mqtt.NewClient(opts)
* 16 ) 使用上面的client 来连接:
```
    if token := client.Connect(); token.Wait() && token.Error() != nil {
	    panic(token.Error())
    }
```

* 17  在进程退出前关闭client 连接： client.Disconnect(1000)

  
* ****
* 如何使用封装好 adaptor.go中：
1） 创建客户端适配器： adp := common.NewAdaptorWithAuth()
2 ) 设置一些属性： adp.SetXXXX()
3 ) 连接客户端： adp.Connect() 
4 ）发布消息 adp.Publish()


* #
* 参考示例：https://github.com/eclipse-paho/paho.mqtt.golang/blob/master/cmd