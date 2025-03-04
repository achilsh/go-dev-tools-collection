package codecs

import (
	"google.golang.org/protobuf/proto"
)

type ProtoCodecs struct {
}

// Encode 将 pb 结构体序列化成二进制。
func (pc *ProtoCodecs) Encode(in any) ([]byte, error) {
	pbIn := in.(proto.Message)
	ret, err := proto.Marshal(pbIn)
	return ret, err
}
func (pc *ProtoCodecs) Decode(data []byte, dst any) error {
	msg, _ := dst.(proto.Message)
	return proto.Unmarshal(data, msg)
}
