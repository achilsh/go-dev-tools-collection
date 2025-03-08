package main

import (
	"fmt"
	"net"
	// "time"

	demo "server-transport-go-usage/gen/go/proto"
	"server-transport-go-usage/lib/codecs"
	"server-transport-go-usage/lib/message"
)

func main() {

	var testDemoSeq uint16 = 0
	serverAddr, err := net.ResolveUDPAddr("udp4", "0.0.0.0:5603")
	if err != nil {
		panic(err)
	}

	// 创建 UDP 连接
	conn, err := net.DialUDP("udp4", nil, serverAddr)
	if err != nil {
		panic(err)
	}
	defer conn.Close()

	for i := 0; i < 100; i++ {
		testDemoSeq++
		toSendMsg := &message.DecodedMessage{
			HeaderMessage: &message.HeaderMessage{
				StartFlag: message.PKG_START_FLAG,
				PkgSeq:    testDemoSeq,
				DevType:   1,
				PkgType:   1,
			},
			DecodedMsg: &demo.BizReqMsg{
				Name: "achilsh",
				Age:  120,
			},
		}

		payload := toSendMsg.GetMessage()
		binPayLoad, err := (&codecs.ProtoCodecs{}).Encode(payload)
		if err != nil {
			fmt.Println("encode fail, err: ", err)
			continue
		}
		toSendMsg.SetMessage(binPayLoad)
		//设置 序列化 payLoad 的数据长度。
		toSendMsg.SetPayLoadLen(uint16(len(binPayLoad)))

		var buf = make([]byte, 1024)
		n, err := toSendMsg.PackageMessage(buf)
		if err != nil {
			fmt.Println("package to buf fail, err: ", err)
			continue
		}
		m, err := conn.Write(buf[:n])
		if err != nil {
			fmt.Println("send buf len fail, err: ", err)
			continue
		}
		fmt.Println("send msg len: ", m, n)
		if m != n {
			fmt.Println("send len less than")
		}
		// time.Sleep(5 * time.Second)
	}

	// // 2. 发送消息到服务端
	// message := []byte("Hello UDP Server!")
	// _, err = conn.Write(message)
	// if err != nil {
	// 	panic(err)
	// }
	// fmt.Printf("已发送消息: %s\n", message)

	// // 3. 设置接收超时（可选）
	// conn.SetReadDeadline(time.Now().Add(5 * time.Second))

	// // 4. 读取服务端响应
	// buffer := make([]byte, 1024)
	// n, _, err := conn.ReadFromUDP(buffer)
	// if err != nil {
	// 	panic(err)
	// }

	// fmt.Printf("收到服务端响应: %s\n", string(buffer[:n]))
}
