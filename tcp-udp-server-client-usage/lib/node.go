package lib

import (
	"fmt"
	"sync"
	"time"

	"server-transport-go-usage/lib/frame"
	"server-transport-go-usage/lib/message"
)

var errTerminated = fmt.Errorf("terminated")

type writeToReq struct {
	ch   *Channel
	what interface{}
}

type writeExceptReq struct {
	except *Channel
	what   interface{}
}

// Node: 是请求，相应的节点抽象 
type Node struct {
	// conf 是 配置 node的信息，包括 节点连接信息等。
	conf             NodeConf
	dialectRW        *message.ReadWriter
	wg               sync.WaitGroup
	channelProviders map[*channelProvider]struct{}
	channels         map[*Channel]struct{}

	// in 
	// 用于通知新连接的通道
	chNewChannel   chan *Channel
	chCloseChannel chan *Channel
	// 用于通知 写数据的通道chWriteTo
	chWriteTo      chan writeToReq
	chWriteAll     chan interface{}
	chWriteExcept  chan writeExceptReq
	terminate      chan struct{}

	// out
	chEvent chan Event
	done    chan struct{}
}

// NewNode allocates a Node. See NodeConf for the options.
func NewNode(conf NodeConf) (*Node, error) {
	if len(conf.Endpoints) == 0 {
		return nil, fmt.Errorf("at least one endpoint must be provided")
	}
	if conf.ReadTimeout == 0 {
		conf.ReadTimeout = 10 * time.Second
	}
	if conf.WriteTimeout == 0 {
		conf.WriteTimeout = 10 * time.Second
	}
	if conf.IdleTimeout == 0 {
		conf.IdleTimeout = 60 * time.Second
	}

	dialectRW, err := func() (*message.ReadWriter, error) {
		if conf.Dialect == nil {
			return nil, nil
		}
		return message.NewReadWriter(conf.Dialect)
	}()
	if err != nil {
		return nil, err
	}

	n := &Node{
		conf:             conf,
		dialectRW:        dialectRW,
		channelProviders: make(map[*channelProvider]struct{}),
		channels:         make(map[*Channel]struct{}),
		chNewChannel:     make(chan *Channel),
		chCloseChannel:   make(chan *Channel),
		chWriteTo:        make(chan writeToReq),
		chWriteAll:       make(chan interface{}),
		chWriteExcept:    make(chan writeExceptReq),
		terminate:        make(chan struct{}),
		chEvent:          make(chan Event),
		done:             make(chan struct{}),
	}

	closeExisting := func() {
		for ch := range n.channels {
			ch.close()
		}
		for ca := range n.channelProviders {
			ca.close()
		}
	}

	// endpoints
	for _, tconf := range conf.Endpoints {
		tp, err := tconf.init(n)
		if err != nil {
			closeExisting()
			return nil, err
		}

		switch ttp := tp.(type) {
		case endpointChannelProvider:
			ca, err := newChannelProvider(n, ttp) //ttp is endpoint;
			if err != nil {
				closeExisting()
				return nil, err
			}

			n.channelProviders[ca] = struct{}{}

		// case endpointChannelSingle:
		// 	ch, err := newChannel(n, ttp, ttp.label(), ttp)
		// 	if err != nil {
		// 		closeExisting()
		// 		return nil, err
		// 	}

		// 	n.channels[ch] = struct{}{}

		default:
			panic(fmt.Errorf("endpoint %T does not implement any interface", tp))
		}
	}

	// for ch := range n.channels {
	// 	ch.start()
	// }

	for ca := range n.channelProviders {
		ca.start()
	}

	go n.run()

	return n, nil
}

func (n *Node) pushEvent(evt Event) {
	select {
	case n.chEvent <- evt:
	case <-n.terminate:
	}
}

func (n *Node) newChannel(ch *Channel) {
	select {
	case n.chNewChannel <- ch:
	case <-n.terminate:
		ch.close()
	}
}

func (n *Node) closeChannel(ch *Channel) {
	select {
	case n.chCloseChannel <- ch:
	case <-n.terminate:
	}
}

// Close halts node operations and waits for all routines to return.
func (n *Node) Close() {
	close(n.terminate)
	<-n.done
}

func (n *Node) run() {
	defer close(n.done)

outer:
	for {
		select {
		case ch := <-n.chNewChannel:
			n.channels[ch] = struct{}{}
			ch.start()

		case ch := <-n.chCloseChannel:
			delete(n.channels, ch)

		case req := <-n.chWriteTo: // 框架统一接收写数据通知（数据和具体某个连接通道）
			if _, ok := n.channels[req.ch]; !ok {
				continue
			}
			req.ch.write(req.what)

		case what := <-n.chWriteAll:
			for ch := range n.channels {
				ch.write(what)
			}

		case req := <-n.chWriteExcept:
			for ch := range n.channels {
				if ch != req.except {
					ch.write(req.what)
				}
			}

		case <-n.terminate:
			break outer
		}
	}

	for ca := range n.channelProviders {
		ca.close()
	}

	for ch := range n.channels {
		ch.close()
	}

	n.wg.Wait()

	close(n.chEvent)
}

func (n *Node) encodeMessage(msg message.Message) (message.Message, error) {
	if _, ok := msg.(*message.MessageRaw); !ok {
		if n.dialectRW == nil {
			return nil, fmt.Errorf("dialect is nil")
		}

		mp := n.dialectRW.GetMessage(msg.GetID())
		if mp == nil {
			return nil, fmt.Errorf("message is not in the dialect")
		}

		msgRaw := mp.Write(msg, n.conf.OutVersion == V2)
		return msgRaw, nil
	}

	return msg, nil
}

// Events returns a channel from which receiving events. Possible events are:
//
// * EventChannelOpen
// * EventChannelClose
// * EventFrame
// * EventParseError
// * EventStreamRequested
//
// See individual events for details.
func (n *Node) Events() chan Event {
	return n.chEvent
}

// WriteMessageTo writes a message to given channel.
func (n *Node) WriteMessageTo(channel *Channel, m message.Message) error {
	m, err := n.encodeMessage(m)
	if err != nil {
		return err
	}

	select {
	case n.chWriteTo <- writeToReq{channel, m}:
	case <-n.terminate:
	}

	return nil
}

// WriteMessageAll writes a message to all channels.
func (n *Node) WriteMessageAll(m message.Message) error {
	m, err := n.encodeMessage(m)
	if err != nil {
		return err
	}

	select {
	case n.chWriteAll <- m:
	case <-n.terminate:
	}

	return nil
}

// WriteMessageExcept writes a message to all channels except specified channel.
func (n *Node) WriteMessageExcept(exceptChannel *Channel, m message.Message) error {
	m, err := n.encodeMessage(m)
	if err != nil {
		return err
	}

	select {
	case n.chWriteExcept <- writeExceptReq{exceptChannel, m}:
	case <-n.terminate:
	}

	return nil
}

// WriteFrameTo writes a frame to given channel.
// This function is intended only for routing pre-existing frames to other nodes,
// since all frame fields must be filled manually.
// 其中 MsgFrame 填充 payload 字段 或者 nil 值。其他的一些关键字段
// 在框架上 向指定 chan 上发送业务数据；可以在多线程上调用。
func (n *Node) WriteFrameTo(channel *Channel, fr frame.MsgFrame) error {
	err := n.encodeFrame(fr) // 消息体序列化
	if err != nil {
		return err
	}

	select {
	case n.chWriteTo <- writeToReq{channel, fr}: //发送到框架的统一写队列上。
	case <-n.terminate:
	}

	return nil
}
// 把数据帧序列化成二进制数据
func (n *Node) encodeFrame(fr frame.MsgFrame) error {
	payLoad := fr.GetMessage()
	if payLoad == nil {
		return nil
	}
	
	binPayLoad, err := n.dialectRW.GetCodecs().Encode(payLoad)
	if err != nil {
		fmt.Println("encode msg fail, err: ", err)
		return err
	}

	fr.SetMessage(binPayLoad)
	//设置 序列化 payLoad 的数据长度。
	fr.SetPayLoadLen(uint16(len(binPayLoad)))
	return nil
}

// WriteFrameAll writes a frame to all channels.
// This function is intended only for routing pre-existing frames to other nodes,
// since all frame fields must be filled manually.
func (n *Node) WriteFrameAll(fr frame.MsgFrame) error {
	err := n.encodeFrame(fr)
	if err != nil {
		return err
	}

	select {
	case n.chWriteAll <- fr:
	case <-n.terminate:
	}

	return nil
}

// WriteFrameExcept writes a frame to all channels except specified channel.
// This function is intended only for routing pre-existing frames to other nodes,
// since all frame fields must be filled manually.
func (n *Node) WriteFrameExcept(exceptChannel *Channel, fr frame.MsgFrame) error {
	err := n.encodeFrame(fr)
	if err != nil {
		return err
	}

	select {
	case n.chWriteExcept <- writeExceptReq{exceptChannel, fr}:
	case <-n.terminate:
	}

	return nil
}
