package lib

import (
	"context"
	"encoding/json"
	"fmt"
	"runtime"
	"strings"
	"time"

	"github.com/IBM/sarama"
)

// MessagePublisher 消息队列的发布接口
type MessagePublisher interface {
	Publish(ctx context.Context, topic string, partition string, event any) error
	Close()
}

// KafkaPublisher 消息队列发布-kafka的实现
type KafkaPublisher struct {
	producer sarama.SyncProducer
}

// ProducerConfig 定义设置 Config 参数
type ProducerConfig func(*sarama.Config)

// WithProducerRequireAcks 设置 producer ack type
func WithProducerRequireAcks(ack int) ProducerConfig {
	return func(cfg *sarama.Config) {
		if cfg == nil {
			return
		}
		// ack:  sarama.WaitForAll; sarama.WaitForLocal; sarama.NoResponse
		cfg.Producer.RequiredAcks = sarama.RequiredAcks(ack)
	}
}

// WithProducerPartitions 设置 partition; 比如： sarama.NewHashPartitioner / sarama.NewRandomPartitioner
func WithProducerPartitions(fn func(topic string) sarama.Partitioner) ProducerConfig {
	return func(cfg *sarama.Config) {
		if cfg == nil {
			return
		}
		cfg.Producer.Partitioner = fn
	}
}

// WithProducerReturn 在 succ chan 返回 连续投递的消息； 默认为true.
func WithProducerReturn(Successes bool) ProducerConfig {
	return func(cfg *sarama.Config) {
		if cfg == nil {
			return
		}
		cfg.Producer.Return.Successes = Successes
	}
}

// NewKafkaPublisher 创建一个kafka 消息队列发布者, broker是节点列表;可选参数就是上面的选项参数。
func NewKafkaPublisher(broker []string, opts ...ProducerConfig) MessagePublisher {
	cfg := sarama.NewConfig()

	for _, opt := range opts {
		opt(cfg)
	}

	if len(broker) == 0 {
		fmt.Printf("producer broker node is empty")
		return nil
	}

	producer, err := sarama.NewSyncProducer(broker, cfg)
	if err != nil || producer == nil {
		fmt.Printf("create producer fail, err: %v", err)
		return nil
	}

	return &KafkaPublisher{
		producer: producer,
	}
}

// Close 关闭生产者
func (p *KafkaPublisher) Close() {
	if p == nil {
		return
	}
	if p.producer == nil {
		return
	}
	p.producer.Close()
}

// Publish 产生者发布消息
func (p *KafkaPublisher) Publish(ctx context.Context, topic string, partitionKey string, event any /*eg: *AlarmRecord */) error {
	body, err := json.Marshal(event)
	if err != nil {
		return err
	}

	message := &sarama.ProducerMessage{
		Key:   sarama.StringEncoder(partitionKey),
		Topic: topic,
		Value: sarama.StringEncoder(body),
	}

	partition, offset, err := p.producer.SendMessage(message)
	if err != nil {
		fmt.Printf("topic: %v, send kafka msg fail, err: %v", topic, err)
		return err
	}
	_, _ = partition, offset

	// fmt.Printf("topic: %v, Message sent to partition: %d at offset: %d", topic, partition, offset)
	return nil
}

// MessageHandler  定义消息处理器接口
type MessageHandler interface {
	Handle(ctx context.Context, message interface{}) error
	HandleBatch(ctx context.Context, message []*sarama.ConsumerMessage) error
}

// MessageConsumer 消费者定义接口
type MessageConsumer interface {
	Consume(ctx context.Context, topic string, handler MessageHandler) error
	Close()
}

// ConsumerConfig 定义设置 Config 参数
type ConsumerConfig func(*sarama.Config)

// WithConsumerReturnError 设置消费时发生错误返回
func WithConsumerReturnError(flag bool) ConsumerConfig {
	return func(cfg *sarama.Config) {
		if cfg == nil {
			return
		}
		cfg.Consumer.Return.Errors = flag // default: true
		cfg.Consumer.Offsets.Initial = sarama.OffsetOldest
	}
}

// WithConsumerOffset 设置消费者的消费偏移量 sarama.OffsetNewest sarama.OffsetOldest
func WithConsumerOffset(offset int64) ConsumerConfig {
	return func(cfg *sarama.Config) {
		if cfg == nil {
			return
		}
		cfg.Consumer.Offsets.Initial = offset
	}
}

// KafkaConsumer kafka实现消费者
type KafkaConsumer struct {
	client sarama.ConsumerGroup
}

// NewKafkaConsumer 创建 kafka 消费者对象
func NewKafkaConsumer(brokers []string, groupID string, options ...ConsumerConfig) (MessageConsumer, error) {
	cfg := sarama.NewConfig()
	for _, opt := range options {
		opt(cfg)
	}
	version, err := sarama.ParseKafkaVersion("3.5.1")
	if err != nil {
		fmt.Printf("parse kafka version 3.5.1 fail, err: %v", err)
		return nil, err
	}
	cfg.Version = version

	client, err := sarama.NewConsumerGroup(brokers, groupID, cfg)
	if err != nil {
		fmt.Printf("new consumer group fail, err: %v", err)
		return nil, err
	}

	return &KafkaConsumer{
		client: client,
	}, nil
}

// ConsumerGroupHandler 消费者分组 消息处理对象
type ConsumerGroupHandler struct {
	handler MessageHandler
}

func (h *ConsumerGroupHandler) Setup(sarama.ConsumerGroupSession) error   { return nil }
func (h *ConsumerGroupHandler) Cleanup(sarama.ConsumerGroupSession) error { return nil }

// ConsumeClaim 批量消费和确认
func (h *ConsumerGroupHandler) ConsumeClaim(sess sarama.ConsumerGroupSession, claim sarama.ConsumerGroupClaim) error {
	tmrDurTime := 500 * time.Millisecond
	timer := time.NewTimer(tmrDurTime)
	defer timer.Stop()
	const batchSize = 1000
	originMessage := make([]*sarama.ConsumerMessage, 0, batchSize)

	for {
		select {
		case message, ok := <-claim.Messages():
			if !ok {
				fmt.Printf("receive close consumer.")
				return nil
			}
			originMessage = append(originMessage, message)
			if len(originMessage) >= batchSize {
				if err := h.handler.HandleBatch(sess.Context(), originMessage); err != nil {
					fmt.Printf("Message topic:%q partition:%d offset:%d, err: %v", message.Topic, message.Partition, message.Offset, err)
					return err
				}
				h.processAndCommit(originMessage, sess)
				originMessage = originMessage[:0]
				timer.Reset(tmrDurTime)
			}
		case <-timer.C:
			if len(originMessage) > 0 {
				if err := h.handler.HandleBatch(sess.Context(), originMessage); err != nil {
					fmt.Printf("Message topic:%v partition:%d offset:%d, err: %v", originMessage[0].Topic, originMessage[0].Partition, originMessage[0].Offset, err)
					return err
				}
				h.processAndCommit(originMessage, sess)
				originMessage = originMessage[:0]
			}
			timer.Reset(tmrDurTime)

		case <-sess.Context().Done():
			fmt.Printf("consumer receive stop signal, out from receiving.")
			return nil
		}
	}
}

// processAndCommit 手动提交批量数据
func (h *ConsumerGroupHandler) processAndCommit(messages []*sarama.ConsumerMessage, sess sarama.ConsumerGroupSession) {
	//for _, msg := range messages {
	//	fmt.Printf("Consumed message: %s, Topic: %s, Partition: %d, Offset: %d", string(msg.Value), msg.Topic, msg.Partition, msg.Offset)
	//}

	if len(messages) > 0 {
		lastMessage := messages[len(messages)-1]
		sess.MarkMessage(lastMessage, "")
		//fmt.Printf("Batch committed up to offset: %d, Partition: %d", lastMessage.Offset, lastMessage.Partition)
	}
}

// Consume 实现kafka 的消费者 消费接口, 需要传入实际消息处理函数 handler
func (c *KafkaConsumer) Consume(ctx context.Context, topic string, handler MessageHandler /* 消息处理函数 */) error {
	consumerGroupHandler := &ConsumerGroupHandler{handler: handler}

	go func() {
		defer func() {
			if e := recover(); e != nil {
				fmt.Printf("Kafka Consume panic error: %s", ErrTrace(fmt.Sprintf("%+v", e)))
			}
		}()

		for {
			if err := c.client.Consume(ctx, []string{topic}, consumerGroupHandler); err != nil {
				fmt.Printf("Error from consumer: %v, topic: %v", err, topic)
			}
		}
	}()

	<-ctx.Done()
	return ctx.Err()
}

// Close 调用关闭消费者组对象，否则会出现内存泄漏
func (c *KafkaConsumer) Close() {
	if c == nil || c.client == nil {
		return
	}
	c.client.Close()
}


func ErrTrace(message string) string {
	var pcs [32]uintptr
	n := runtime.Callers(3, pcs[:]) // skip first 3 caller

	var str strings.Builder
	str.WriteString(message + "\nTraceback:")
	for _, pc := range pcs[:n] {
		fn := runtime.FuncForPC(pc)
		file, line := fn.FileLine(pc)
		str.WriteString(fmt.Sprintf("\n\t%s:%d", file, line))
	}
	return str.String()
}
