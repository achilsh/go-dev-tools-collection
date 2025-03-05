package lib

import (
	"errors"

	. "server-transport-go-usage/lib/utils"
)

type channelProvider struct {
	n   *Node
	eca endpointChannelProvider //是一个服务器的节点，用于接受外面连接请求。

	terminate chan struct{}
}

func newChannelProvider(n *Node, eca endpointChannelProvider) (*channelProvider, error) {
	return &channelProvider{
		n:         n,
		eca:       eca,
		terminate: make(chan struct{}),
	}, nil
}

func (cp *channelProvider) close() {
	close(cp.terminate)
	cp.eca.close()
}

func (cp *channelProvider) start() {
	cp.n.wg.Add(1)
	//启动一个新线程 循环等待新连接，如果有新连接请求，则创建新连接 chan，并将 chan 发送到 框架node 统一处理。
	go cp.run()
}

func (cp *channelProvider) run() {
	defer cp.n.wg.Done()

	for {
		label, rwc, err := cp.eca.provide()
		if err != nil {
			if !errors.Is(err, errTerminated) {
				LogPrintln("errTerminated is the only error allowed here")
			}
			break
		}

		ch, err := newChannel(cp.n, cp.eca, label, rwc)
		if err != nil {
			LogPrintf("newChannel unexpected error: %v\n", err)
		}

		cp.n.newChannel(ch) //把创建的新连接 发送到框架Node的统一处理；

		if cp.eca.oneChannelAtAtime() {
			select {
			case <-ch.done:
			case <-cp.terminate:
				return
			}
		}
	}
}
