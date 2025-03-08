package lib

import (
	"errors"
	"fmt"
	"io"
	"net"
	"time"

	"github.com/pion/transport/v2/udp"

	"server-transport-go-usage/lib/timednetconn"
	. "server-transport-go-usage/lib/utils"
)

// 非面向连接 udp 服务
type EndPointUDPNoDirectServer struct {
	Address string
}

func(EndPointUDPNoDirectServer)isUDP() bool {
	return true
}
func(e EndPointUDPNoDirectServer) getAddress() string {
	return e.Address
}

func (conf EndPointUDPNoDirectServer) init (node *Node) (Endpoint , error) {
	return initEndpointUDPNoDirectServer(node, conf)
}

func initEndpointUDPNoDirectServer(node *Node, conf endpointServerConf) (Endpoint, error) {
	conn, err := net.ListenPacket("udp4", conf.getAddress())
    if err != nil {
        return nil, fmt.Errorf("invalid address")
    }

	t := &endPointUDPNoDirectServer{
		conf:         conf,
		writeTimeout: node.conf.WriteTimeout,
		idleTimeout:  node.conf.IdleTimeout,
		listener:     conn,
		terminate:    make(chan struct{}),
	}
	return t, nil
}

type endPointUDPNoDirectServer struct {
	conf endpointServerConf
	listener     net.PacketConn
	writeTimeout time.Duration
	idleTimeout  time.Duration

	// in
	terminate chan struct{}
}
func(t *endPointUDPNoDirectServer) getConn() net.PacketConn{
	return t.listener
}

func (t *endPointUDPNoDirectServer) isEndpoint() {}

func (t *endPointUDPNoDirectServer) Conf() EndpointConf {
	return t.conf
}
func (t *endPointUDPNoDirectServer) close() {
	close(t.terminate)
	//
	LogPrintf("close end udp-no-direct-connect server.")
	t.listener.Close()
}

func (t *endPointUDPNoDirectServer) oneChannelAtAtime() bool {
	return false
}
func (t *endPointUDPNoDirectServer) label() string {
	return fmt.Sprintf("udp:%s", t.listener)
}



// udp 单播服务
type EndPointUDPSingleServer struct {
	Address string
}


// 实现 一个 conn 对象； 使用该对象 直接进行 ReadFrom 数据.
func (conf EndPointUDPSingleServer) init(node *Node) (Endpoint, error) {
	pc, err := net.ListenPacket("udp4", conf.Address)
	if err != nil {
		return nil, err
	}

	t := &endPointUDPSingleServer{
		conf: conf,
		pc:   pc,
		writeTimeout: node.conf.WriteTimeout,
		idleTimeout:  node.conf.IdleTimeout,
		terminate:    make(chan struct{}),
	}
	return t, nil
}

type endPointUDPSingleServer struct {
	conf         EndPointUDPSingleServer
	pc           net.PacketConn // listen addr.
	writeTimeout time.Duration
	idleTimeout  time.Duration

	// in
	terminate chan struct{}
}

func (t *endPointUDPSingleServer) isEndpoint() {}

func (t *endPointUDPSingleServer) Conf() EndpointConf {
	return t.conf
}

func (t *endPointUDPSingleServer) label() string {
	return fmt.Sprintf("udp:%s", t.pc)
}

func (t *endPointUDPSingleServer) Close() error {
	close(t.terminate)
	t.pc.Close()
	return nil
}
func (t *endPointUDPSingleServer) SetReadDeadline() error {
	err := t.pc.SetWriteDeadline(time.Now().Add(t.idleTimeout))
	return err
}
func (t *endPointUDPSingleServer) SetWriteDeadline() error {
	err := t.pc.SetWriteDeadline(time.Now().Add(t.writeTimeout))
	return err
}

// 读数据
func (t *endPointUDPSingleServer) Read(buf []byte) (int, error) {
	// read WITHOUT deadline. Long periods without packets are normal since
	// we're not directly connected to someone.
	n, _, err := t.pc.ReadFrom(buf)
	// wait termination, do not report errors
	if err != nil {
		<-t.terminate
		return 0, errTerminated
	}

	return n, nil
}

// 写数据
func (t *endPointUDPSingleServer) Write(buf []byte) (int, error) {
	err := t.pc.SetWriteDeadline(time.Now().Add(t.writeTimeout))
	if err != nil {
		return 0, err
	}
	return t.pc.WriteTo(buf, t.pc.LocalAddr())
}

// EndpointUDPServer sets up a endpoint that works with an UDP server.
// if they are connected to the same network.
type EndpointUDPServer struct {
	// listen address, example: 0.0.0.0:5600
	Address string
}

func (EndpointUDPServer) isUDP() bool {
	return true
}

func (conf EndpointUDPServer) getAddress() string {
	return conf.Address
}


type endpointServerConf interface {
	isUDP() bool
	getAddress() string
	init(*Node) (Endpoint, error)
}

// EndpointTCPServer sets up a endpoint that works with a TCP server.
// TCP is fit for routing frames through the internet
type EndpointTCPServer struct {
	// listen address, example: 0.0.0.0:5600
	Address string
}

func (EndpointTCPServer) isUDP() bool {
	return false
}

func (conf EndpointTCPServer) getAddress() string {
	return conf.Address
}



func (conf EndpointTCPServer) init(node *Node) (Endpoint, error) {
	return initEndpointServer(node, conf)
}

// 创建一个 监听服务
func (conf EndpointUDPServer) init(node *Node) (Endpoint, error) {
	return initEndpointServer(node, conf)
}

func initEndpointServer(node *Node, conf endpointServerConf) (Endpoint, error) {
	_, _, err := net.SplitHostPort(conf.getAddress())
	if err != nil {
		return nil, fmt.Errorf("invalid address")
	}

	var ln net.Listener
	if conf.isUDP() {
		var addr *net.UDPAddr
		addr, err = net.ResolveUDPAddr("udp4", conf.getAddress())
		if err != nil {
			return nil, err
		}

		ln, err = udp.Listen("udp4", addr)
		if err != nil {
			return nil, err
		}
	} else {
		ln, err = net.Listen("tcp4", conf.getAddress())
		if err != nil {
			return nil, err
		}
	}

	t := &endpointServer{
		conf:         conf,
		writeTimeout: node.conf.WriteTimeout,
		idleTimeout:  node.conf.IdleTimeout,
		listener:     ln,
		terminate:    make(chan struct{}),
	}
	return t, nil
}

// endpointServer 是一个服务端的实现
type endpointServer struct {
	conf         endpointServerConf
	listener     net.Listener
	writeTimeout time.Duration
	idleTimeout  time.Duration

	// in
	terminate chan struct{}
}

func (t *endpointServer) isEndpoint() {}

func (t *endpointServer) Conf() EndpointConf {
	return t.conf
}

func (t *endpointServer) close() {
	close(t.terminate)
	//
	LogPrintf("close end pointer server.")
	t.listener.Close()
}

func (t *endpointServer) oneChannelAtAtime() bool {
	return false
}

func (t *endpointServer) provide() (string, io.ReadWriteCloser, error) {
	nconn, err := t.listener.Accept()
	// wait termination, do not report errors
	if err != nil {
		if errors.Is(err, net.ErrClosed) {
			LogPrintln("Listener closed, exiting Accept loop")
		}
		LogPrintf("accept fail, err: %v", err)
		<-t.terminate
		return "", nil, errTerminated
	}
	LogPrintf("receive new connection, %v", nconn)

	label := fmt.Sprintf("%s:%s", func() string {
		if t.conf.isUDP() {
			return "udp"
		}
		return "tcp"
	}(), nconn.RemoteAddr())

	conn := timednetconn.New(
		t.idleTimeout,
		t.writeTimeout,
		nconn)

	return label, conn, nil
}
