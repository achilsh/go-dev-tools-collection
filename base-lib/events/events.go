package events

import (
	"go-micro.dev/v4/events"
)

// Mq 消息队列接口
type Mq interface {
	Publish(msg interface{})
	Subscribe(deal func(event events.Event), opts ...SubscribeOption)
	UnSubscribe()
}

// MqRestrict 消息队列相关参数
//
// 主题，消费者组名，消费传递模式等
type MqRestrict struct {
	Topic         string
	Group         string
	ConsumeModule string
}

// ConnectOptions 所有连接参数
type ConnectOptions struct {
	Mr MqRestrict
}

// ConnectOption 设置连接参数
type ConnectOption func(o *ConnectOptions)

// WithMqRestrict 设置MqRestrict连接参数
func WithMqRestrict(mr MqRestrict) ConnectOption {
	return func(o *ConnectOptions) {
		o.Mr = mr
	}
}

const (
	// DefaultNumInGroup 默认消费者组内数量
	DefaultNumInGroup = 1
)

// SubscribeOptions 所有订阅参数
type SubscribeOptions struct {
	NumInGroup int
}

// SubscribeOption 设置订阅参数
type SubscribeOption func(o *SubscribeOptions)

// WithConsumeNumOfGroup 设置消费者组内数量订阅参数
func WithConsumeNumOfGroup(num int) SubscribeOption {
	return func(o *SubscribeOptions) {
		o.NumInGroup = num
	}
}

// ConsumeModule 消息传递模式
type ConsumeModule int8

const (
	// AtMostOnce 至多一次消息传递模式
	AtMostOnce ConsumeModule = iota
	// AtLeastOnce 至少一次消息传递模式
	AtLeastOnce
	// ExactlyOnce 精确一次消息传递模式
	ExactlyOnce
)

// 消息传递模式对应的字符串
func (o ConsumeModule) String() string {
	switch o {
	case AtMostOnce:
		return "at-most-once"
	case AtLeastOnce:
		return "at-least-once"
	case ExactlyOnce:
		return "exactly-once"
	}
	return "at-most-once"
}
