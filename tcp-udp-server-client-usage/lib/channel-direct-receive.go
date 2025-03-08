package lib

import (
	"net"

	"server-transport-go-usage/lib/message"
	. "server-transport-go-usage/lib/utils"
)

// udp 非连接 启动后 收包 step；只有一个实例，内部处理函数：使用独立 协程 循环接收 和读数据。
type channelNoDirectReceiver struct {
	n   *Node
	eca endpointNoDirectReceiver // udp 非连接 启动后的服务

	terminate chan struct{}
}

// 入参 eca 是一个非连接 udp服务节点。
func newChannelDirectReceiver(n *Node, eca endpointNoDirectReceiver) (*channelNoDirectReceiver, error) {
	return &channelNoDirectReceiver{
		n:         n,
		eca:       eca,
		terminate: make(chan struct{}),
	}, nil
}

func (cp *channelNoDirectReceiver) start() {
	cp.n.wg.Add(1)
	go cp.run()
}
func (cp *channelNoDirectReceiver) close() {
	close(cp.terminate)
	cp.eca.close()
}



func (cp *channelNoDirectReceiver) run() {
	defer cp.n.wg.Done()

	// 非连接性请求： 可循环读取数据： 根据数据组装一个对象（数据，源对象）；把数据发给框架通道；
	for {
		select {
		case <-cp.terminate:
			return 
		default:
			var receiveBuf = make([]byte, 1024)
			n, clientAddr, err := cp.eca.getConn().ReadFrom(receiveBuf)
			if err != nil {
				LogPrintf("read data fail, err: %v", err)
				continue
			}
			if n <= 0 {
				LogPrintf("read udp data len is 0")
				continue
			}

			//
			originData := NewUdpChanData(clientAddr,receiveBuf[:n],cp.eca,cp.n)

			select {
			case cp.n.udpOriginDataCh <-originData:
			case <-cp.n.terminate:
				return 
			}
		}
	}
}

//////
type UdpChanData struct {
	clientAddr net.Addr
	recvData []byte 
	udpServerNode endpointNoDirectReceiver //单实例 udp 非连接 启动后的服务
	n *Node
}
func NewUdpChanData(clientAddr net.Addr, recvData []byte, 
	udpServerNode endpointNoDirectReceiver, n *Node)*UdpChanData {
	return &UdpChanData{
		clientAddr: clientAddr,
		recvData :recvData,
		udpServerNode :udpServerNode, //单实例 udp 非连接 启动后的服务
		n :n,
	}
}
func (ud *UdpChanData) Parse() {
	go func() {
		pkg := &message.DecodedMessage{}
		err := pkg.ParseAndValidPkg(ud.recvData)
		if err != nil {
			LogPrintf("pase pkg and check fail, err: %v", err)
			return 
		}
		if err := pkg.Unpackage(ud.n.dialectRW);err !=  nil {
			LogPrintf("unpackage fail, err: %v", err)
			return 
		}
		//发送到框架的 业务处理逻辑中  
		evt := &EventFrame{Frame: pkg, ExtData: ud}
		ud.n.pushEvent(evt)
	}()
}