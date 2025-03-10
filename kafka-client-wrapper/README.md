## kafka client 封装和使用实例
* 封装开源的kafka client, 实现 消费者 和 生产者的对象创建，消费和发布接口
### 使用实例
1. 生产者创建和发布
```
demoPusher = lib.NewKafkaPublisher([]string{"", ""},
		lib.WithProducerReturn(true),
		lib.WithProducerRequireAcks(int(sarama.WaitForAll)),
		lib.WithProducerPartitions(sarama.NewHashPartitioner))

demoPusher.Publish()
```

2. 消费者创建和订阅：
```
var msgHandle = &MessageHandle{}
c, e := lib.NewKafkaConsumer([]string{"",""}, "consumer_group_id", lib.WithConsumerReturnError(true), lib.WithConsumerOffset(sarama.OffsetOldest))
if e !=nil {
    return 
}
c.Consume(context.Backgroud(), "consumer topic", msgHandle)

type MessageHandle struct {}
type DataMsg struct {} //是一个json 数据

// Handle 只记录一条记录
func (ah *MessageHandle) Handle(ctx context.Context, message any) error {
	//忽略处理逻辑
	return nil
}

// HandleBatch 消费kafka消息处理函数
func (ah *MessageHandle) HandleBatch(ctx context.Context, message []*sarama.ConsumerMessage) error {
	var msgItems []*AlarmRecord
	for _, item := range message {
		if item == nil {
			continue
		}

		newItem := &DataMsg{}
		err := json.Unmarshal(item.Value, newItem)
		if err != nil {
			continue
		}
		msgItems = append(msgItems, newItem)
	}
	return nil
}
```
