package memory

import (
	"sync"

	"github.com/google/uuid"
	"go-micro.dev/v4/broker"
	"go-micro.dev/v4/logger"
)

var (
	Broker = NewMemoryBroker()
)

// NewBroker memory消息队列
func NewBroker() broker.Broker {
	return Broker
}

type memoryBroker struct {
	opts *broker.Options

	sync.RWMutex
	Subscribers map[string][]*memorySubscriber
}

func (m *memoryBroker) Options() broker.Options {
	return *m.opts
}

func (m *memoryBroker) Address() string {
	return ""
}

func (m *memoryBroker) Connect() error {
	return nil
}

func (m *memoryBroker) Disconnect() error {
	return nil
}

func (m *memoryBroker) Init(opts ...broker.Option) error {
	for _, o := range opts {
		o(m.opts)
	}
	return nil
}

// Publish 发布消息
//
// topic: 消息主题, msg: 消息体, opts: 用不到
func (m *memoryBroker) Publish(topic string, msg *broker.Message, opts ...broker.PublishOption) error {
	m.RLock()
	subs, ok := m.Subscribers[topic]
	m.RUnlock()
	if !ok {
		return nil
	}

	p := &memoryEvent{
		topic:   topic,
		message: msg,
		opts:    m.opts,
	}

	for _, sub := range subs {
		if err := sub.handler(p); err != nil {
			logger.Error(err)
			continue
		}
	}
	return nil
}

// Subscribe 消息订阅
//
// topic: 消息主题, handler: 消息处理函数, opts: 用不到
func (m *memoryBroker) Subscribe(topic string, handler broker.Handler, opts ...broker.SubscribeOption) (broker.Subscriber, error) {
	var options broker.SubscribeOptions
	for _, o := range opts {
		o(&options)
	}

	sub := &memorySubscriber{
		exit:    make(chan bool, 1),
		id:      uuid.New().String(),
		topic:   topic,
		handler: handler,
		opts:    options,
	}

	m.Lock()
	m.Subscribers[topic] = append(m.Subscribers[topic], sub)
	m.Unlock()

	go func() {
		<-sub.exit
		m.Lock()
		var newSubscribers []*memorySubscriber
		for _, sb := range m.Subscribers[topic] {
			if sb.id == sub.id {
				continue
			}
			newSubscribers = append(newSubscribers, sb)
		}
		m.Subscribers[topic] = newSubscribers
		m.Unlock()
	}()

	return sub, nil
}

func (m *memoryBroker) String() string {
	return "memory"
}

type memoryEvent struct {
	opts    *broker.Options
	topic   string
	message interface{}
	err     error
}

// Topic 返回消息的主题
func (m *memoryEvent) Topic() string {
	return m.topic
}

// Message 返回消息
func (m *memoryEvent) Message() *broker.Message {
	return m.message.(*broker.Message)
}

func (m *memoryEvent) Ack() error {
	return nil
}

func (m *memoryEvent) Error() error {
	return m.err
}

type memorySubscriber struct {
	id      string
	topic   string
	exit    chan bool
	handler broker.Handler
	opts    broker.SubscribeOptions
}

func (m *memorySubscriber) Options() broker.SubscribeOptions {
	return m.opts
}

// Topic 返回订阅消息的主题
func (m *memorySubscriber) Topic() string {
	return m.topic
}

// Unsubscribe 退订消息
func (m *memorySubscriber) Unsubscribe() error {
	m.exit <- true
	return nil
}

func NewMemoryBroker(opts ...broker.Option) broker.Broker {
	options := broker.NewOptions(opts...)

	return &memoryBroker{
		opts:        options,
		Subscribers: make(map[string][]*memorySubscriber),
	}
}
