package broker

import (
	jsoniter "github.com/json-iterator/go"
	"go-micro.dev/v4/broker"
)

func NewMessage(message interface{}, opts ...MessageOption) *broker.Message {
	m := &broker.Message{}

	if len(opts) != 0 {
		opt := MessageOptions{}
		for _, o := range opts {
			o(&opt)
		}
		m.Header = opt.Header
	}

	switch message.(type) {
	case []byte:
		m.Body = message.([]byte)
	default:
		b, _ := jsoniter.Marshal(message)
		m.Body = b
	}
	return m
}

type MessageOptions struct {
	Header map[string]string
}

type MessageOption func(*MessageOptions)

func WithHeader(key, value string) MessageOption {
	return func(o *MessageOptions) {
		o.Header[key] = value
	}
}
