package frame

import (
	"bufio"
	"errors"
	"fmt"
	"io"
	"server-transport-go-usage/lib/message"
	"server-transport-go-usage/lib/utils"
	"time"
)

const (
	bufferSize = 1024 // frames default size is 1k.
)

// 1st January 2015 GMT
var signatureReferenceDate = time.Date(2015, 0o1, 0o1, 0, 0, 0, 0, time.UTC)

// ReaderConf is the configuration of a Reader.
type ReaderConf struct {
	// the underlying bytes reader.
	Reader io.Reader // 新的连接

	// (optional) the dialect which contains the messages that will be read.
	// If not provided, messages are decoded into the MessageRaw struct.
	DialectRW *message.ReadWriter
}

// Reader is a Frame reader.
type Reader struct {
	conf                 ReaderConf
	br                   *bufio.Reader // 新连接
	curReadSignatureTime uint64
	scer                 *bufio.Scanner // 用于读取，切分协议
}

// NewReader allocates a Reader.
func NewReader(conf ReaderConf) (*Reader, error) {
	if conf.Reader == nil {
		return nil, fmt.Errorf("Reader not provided")
	}

	ret := &Reader{
		conf: conf,
		br:   bufio.NewReaderSize(conf.Reader, bufferSize),
		scer: bufio.NewScanner(conf.Reader),
	}

	ret.scer.Buffer(make([]byte, bufferSize), bufferSize*2)
	ret.scer.Split(message.ProtocolSplitter) //TODO:
	return ret, nil
}

func ResetScanner(r *Reader) {
	r.scer = bufio.NewScanner(r.conf.Reader)
	r.scer.Buffer(make([]byte, bufferSize), bufferSize*2)
	r.scer.Split(message.ProtocolSplitter)
}

// ReadPkg reads a Frame from the reader.
// It must not be called by multiple routines in parallell.
// ReadPkg 从底层网络中读取数据
func (r *Reader) ReadPkg() (MsgFrame, error) {
	item := &message.DecodedMessage{}
	///
	ioRW, ok := r.conf.Reader.(utils.BizIoWRWrapper)
	if !ok {
		utils.LogPrintf("io reader is not implement BizIoWRWrapper interface.")
		//
	}
	err := item.UnPackageMessage(r.scer, r.conf.DialectRW, ioRW)
	if errors.As(err, &message.IoTimeoutError{}) {
		ResetScanner(r)
		ioRW.SetReadDeadline()
	}
	return item, err
}
