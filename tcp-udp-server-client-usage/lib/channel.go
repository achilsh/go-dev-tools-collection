package lib

import (
	"context"
	"crypto/rand"
	"errors"
	"io"
	"server-transport-go-usage/lib/frame"
)

const (
	writeBufferSize = 128 //TODO: need check on biz.
)

func randomByte() (byte, error) {
	var buf [1]byte
	_, err := rand.Read(buf[:])
	return buf[0], err
}

// Channel is a communication channel created by an Endpoint.
// An Endpoint can create channels.
// For instance, a TCP client endpoint creates a single channel, while a TCP
// server endpoint creates a channel for each incoming connection.
type Channel struct {
	n     *Node
	e     Endpoint
	label string
	rwc   io.Closer // 新连接

	ctx       context.Context
	ctxCancel func()
	frw       *frame.ReadWriter //在新连接上 做读写能力，封装了 scanner 读协议能力
	running   bool

	// in
	chWrite chan interface{}

	// out
	done chan struct{}
}

// 一个连接包括：基于该连接的读写接口，所有消息的编解码接口实现
func newChannel(
	n *Node,
	e Endpoint,
	label string,
	rwc io.ReadWriteCloser, // 一个新的连接
) (*Channel, error) {
	//每个连接都分配一个读写器；读用于接收和解析协议；写用于编码协议和发送
	frw, err := frame.NewReadWriter(frame.ReadWriterConf{
		ReadWriter: rwc,
		DialectRW:  n.dialectRW,
	})
	if err != nil {
		return nil, err
	}

	ctx, ctxCancel := context.WithCancel(context.Background())

	return &Channel{
		n:         n,
		e:         e,
		label:     label,
		rwc:       rwc,
		ctx:       ctx,
		ctxCancel: ctxCancel,
		frw:       frw,
		chWrite:   make(chan interface{}, writeBufferSize),
		done:      make(chan struct{}),
	}, nil
}

func (ch *Channel) close() {
	ch.ctxCancel()
	if !ch.running {
		ch.rwc.Close()
	}
}

func (ch *Channel) start() {
	ch.running = true
	ch.n.wg.Add(1)
	go ch.run()
}

func (ch *Channel) run() {
	defer close(ch.done)
	defer ch.n.wg.Done()

	// 连接的通道读数据协程。
	readerDone := make(chan struct{})
	go ch.runReader(readerDone)

	/// 连接的通道写数据 协程
	writerTerminate := make(chan struct{})
	writerDone := make(chan struct{})
	go ch.runWriter(writerTerminate, writerDone)

	select {
	case <-readerDone:
		ch.rwc.Close()

		close(writerTerminate)
		<-writerDone

	case <-ch.ctx.Done():
		close(writerTerminate)
		<-writerDone

		ch.rwc.Close()
		<-readerDone
	}

	ch.ctxCancel()

	ch.n.pushEvent(&EventChannelClose{ch})
	ch.n.closeChannel(ch)
}

func (ch *Channel) runReader(readerDone chan struct{}) {
	defer close(readerDone)

	// wait client here, in order to allow the writer goroutine to start
	// and allow clients to write messages before starting listening to events
	ch.n.pushEvent(&EventChannelOpen{ch})

	for {
		fr, err := ch.frw.ReadPkg()
		if err != nil {
			var errRead frame.ReadError
			if errors.As(err, &errRead) {
				ch.n.pushEvent(&EventParseError{err, ch})
				continue
			}
			return
		}

		evt := &EventFrame{Frame: fr, Channel: ch}
		ch.n.pushEvent(evt)
	}
}

// 每个连接的独立写 线程
func (ch *Channel) runWriter(writerTerminate chan struct{}, writerDone chan struct{}) {
	defer close(writerDone)

	for {
		select {
		case what := <-ch.chWrite:
			switch wh := what.(type) {
			case frame.MsgFrame:
				ch.frw.WriteMsgFrame(wh) //nolint:errcheck
			}

		case <-writerTerminate:
			return
		}
	}
}

// String implements fmt.Stringer.
func (ch *Channel) String() string {
	return ch.label
}

// Endpoint returns the channel Endpoint.
func (ch *Channel) Endpoint() Endpoint {
	return ch.e
}

// 向具体指定的 channel上发送写数据（因为每个channel都一个独立的写协程在处理。）
func (ch *Channel) write(what interface{}) {
	select {
	case ch.chWrite <- what:
	case <-ch.ctx.Done():
	default: // buffer is full
	}
}
