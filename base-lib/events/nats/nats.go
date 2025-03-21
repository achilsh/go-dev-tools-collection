package nats

import (
	"context"
	"sync"
	"time"

	mq "github.com/achilsh/go-dev-tools-collection/base-lib/events"
	"github.com/go-micro/plugins/v4/events/natsjs"
	"go-micro.dev/v4/events"
	"go-micro.dev/v4/logger"
)

var (
	// DefaultNatsAddr 默认Nats地址
	DefaultNatsAddr = "BXWLM0bixdze56JftuYUpws94OR2o3qS@127.0.0.1:4222"
	// Consumes 用于记录停止消费者协程的context
	Consumes    = make(map[string]context.CancelFunc)
	ConsumesMux sync.RWMutex
)

// NewStream 创建一个nats-jetStream连接
func NewStream(opts ...mq.ConnectOption) mq.Mq {
	options := mq.ConnectOptions{}
	for _, o := range opts {
		o(&options)
	}

	njs, err := natsjs.NewStream(natsjs.Address(DefaultNatsAddr))
	if err != nil {
		logger.Error(err)
		return nil
	}

	n := &njsStream{
		njs:        njs,
		MqRestrict: options.Mr,
	}
	return n
}

// NjsStream nats-jetStream连接
type njsStream struct {
	mq.MqRestrict
	njs events.Stream
}

// Publish 通过nats-jetStream发布消息
func (n *njsStream) Publish(msg interface{}) {
	err := n.njs.Publish(n.Topic, msg)
	if err != nil {
		logger.Error(err)
	}
}

// Subscribe 通过nats-jetStream订阅消息
func (n *njsStream) Subscribe(deal func(event events.Event), opts ...mq.SubscribeOption) {
	options := mq.SubscribeOptions{
		NumInGroup: mq.DefaultNumInGroup,
	}
	for _, o := range opts {
		o(&options)
	}
	num := options.NumInGroup

	consumes := make([]<-chan events.Event, num)
	for i := 0; i < num; i++ {
		var err error
		consumes[i], err = n.njs.Consume(n.Topic, events.WithGroup(n.Group), GetConsumeModule(n.ConsumeModule))
		if err != nil {
			logger.Errorf("consume topic[%v] error: %v", n.Topic, err)
			return
		}
	}

	ConsumesMux.Lock()
	defer ConsumesMux.Unlock()
	var ctx context.Context
	ctx, Consumes[n.Group] = context.WithCancel(context.Background())
	for i := 0; i < num; i++ {
		go func(idx int) {
			for {
				select {
				case <-ctx.Done():
					return
				case event, ok := <-consumes[idx]:
					if !ok {
						return
					}
					deal(event)
				}
			}
		}(i)
	}
}

// UnSubscribe 停止相应的消费者协程
func (n *njsStream) UnSubscribe() {
	ConsumesMux.Lock()
	defer ConsumesMux.Unlock()
	if Consumes[n.Group] == nil {
		return
	}
	Consumes[n.Group]()
	delete(Consumes, n.Group)
}

// GetConsumeModule 根据消息传递模式对应的字符串，获取events.ConsumeOption
func GetConsumeModule(s string) events.ConsumeOption {
	switch s {
	//case mq.AtMostOnce.String():
	//	return events.WithAutoAck(true, 0)
	case mq.AtLeastOnce.String():
		return events.WithAutoAck(false, 0)
	case mq.ExactlyOnce.String():
		return events.WithAutoAck(false, 30*time.Second)
	}
	return events.WithAutoAck(false, 0)
}
