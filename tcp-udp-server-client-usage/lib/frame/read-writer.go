package frame

import (
	"io"

	"server-transport-go-usage/lib/message"
	. "server-transport-go-usage/lib/utils"
)

// ReadWriterConf is the configuration of a ReadWriter.
type ReadWriterConf struct {
	// the underlying bytes ReadWriter.
	ReadWriter io.ReadWriter // 新的连接

	// (optional) the dialect which contains the messages that will be read.
	// If not provided, messages are decoded into the MessageRaw struct.
	DialectRW *message.ReadWriter //所有消息的集合
}

// ReadWriter is a Frame Reader and Writer.
type ReadWriter struct {
	*Reader
	*Writer
}

// NewReadWriter allocates a ReadWriter.
func NewReadWriter(conf ReadWriterConf) (*ReadWriter, error) {
	// 创建一个读器: 包含读取数据，解析协议。
	r, err := NewReader(ReaderConf{
		Reader:    conf.ReadWriter, // 新的连接
		DialectRW: conf.DialectRW,
	})
	if err != nil {
		LogPrintf("new reader for created-connect fail, err: %v\n", err)
		return nil, err
	}

	w, err := NewWriter(WriterConf{
		Writer:             conf.ReadWriter,
		DialectRW:          conf.DialectRW,
	})
	if err != nil {
		return nil, err
	}

	return &ReadWriter{
		Reader: r,
		Writer: w,
	}, nil
}
