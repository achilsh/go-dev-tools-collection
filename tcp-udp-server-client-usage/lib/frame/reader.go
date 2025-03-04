package frame

import (
	"bufio"
	"fmt"
	"io"
	"time"

	"server-transport-go-usage/lib/message"
)

const (
	bufferSize = 1024 // frames default size is 1k.
)

// 1st January 2015 GMT
var signatureReferenceDate = time.Date(2015, 0o1, 0o1, 0, 0, 0, 0, time.UTC)

// ReadError is the error returned in case of non-fatal parsing errors.
type ReadError struct {
	str string
}

func (e ReadError) Error() string {
	return e.str
}

func newError(format string, args ...interface{}) ReadError {
	return ReadError{
		str: fmt.Sprintf(format, args...),
	}
}

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
	ret.scer.Split(message.ProtocolSplitter)
	return ret, nil
}

// ReadPkg reads a Frame from the reader.
// It must not be called by multiple routines in parallell.
// ReadPkg 从底层网络中读取 多个包
func (r *Reader) ReadPkg() (MsgFrame, error) {
	item := &message.DecodedMessage{}
	err := item.UnPackageMessage(r.scer, r.conf.DialectRW)
	return item, err
}
