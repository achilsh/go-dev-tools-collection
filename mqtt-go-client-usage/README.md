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
        _ = t.Wait() // 这个动作可以异步来操作，放在独立 goroutinue 中；Can also use '<-t.Done()' in releases > 1.2.0
        if t.Error() != nil {
			fmt.Printf("ERROR PUBLISHING: %s\n", err)
		}
    ```

* 15 ）关闭连接：
    ```
     client.Disconnect(1000)
    ```

* #
#
* sub client 实例关键流程：
* 1） 
* 2）
* 3）
* 4）
* 5）