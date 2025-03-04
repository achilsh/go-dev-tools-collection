package frame

import (
	"fmt"
	"io"

	"server-transport-go-usage/lib/dialect"
	"server-transport-go-usage/lib/message"
)

// WriterConf is the configuration of a Writer.
type WriterConf struct {
	// the underlying bytes writer.
	Writer io.Writer

	// (optional) the dialect which contains the messages that will be written.
	DialectRW *dialect.ReadWriter

	// Mavlink version used to encode messages.
	//OutVersion WriterOutVersion
	// the system id, added to every outgoing frame and used to identify this
	// node in the network.
	OutSystemID byte
	// (optional) the component id, added to every outgoing frame, defaults to 1.
	OutComponentID byte
	// (optional) the value to insert into the signature link id.
	// This feature requires v2 frames.
	OutSignatureLinkID byte
	// (optional) the secret key used to sign outgoing frames.
	// This feature requires v2 frames.
	//OutKey *V2Key
}

// Writer is a Frame writer.
type Writer struct {
	conf                   WriterConf
	bw                     []byte
	curWriteSequenceNumber byte
}

// NewWriter allocates a Writer.
func NewWriter(conf WriterConf) (*Writer, error) {
	if conf.Writer == nil {
		return nil, fmt.Errorf("Writer not provided")
	}

	return &Writer{
		conf: conf,
		bw:   make([]byte, bufferSize),
	}, nil
}

// WriteMessage writes a Message.
// The Message is wrapped into a Frame whose fields are filled automatically.
// It must not be called by multiple routines in parallel.
func (w *Writer) WriteMessage(m message.Message) error {
	//if w.conf.OutVersion == V1 {
	//	return w.writeFrameAndFill(&V1Frame{Message: m})
	//}
	//return w.writeFrameAndFill(&V2Frame{Message: m})
	//TODO:
	return nil
}

func (w *Writer) writeFrameAndFill(fr MsgFrame) error {
	if fr.GetMessage() == nil {
		return fmt.Errorf("message is nil")
	}

	w.curWriteSequenceNumber++

	if w.conf.DialectRW == nil {
		return fmt.Errorf("dialect is nil")
	}
	return w.writeFrameInner(fr)
}

// WriteFrame writes a Frame.
// It must not be called by multiple routines in parallel.
// This function is intended only for routing pre-existing frames to other nodes,
// since all frame fields must be filled manually.
func (w *Writer) WriteFrame(fr MsgFrame) error {
	if fr.GetMessage() == nil {
		return fmt.Errorf("message is nil")
	}
	return w.writeFrameInner(fr)
}

func (w *Writer) writeFrameInner(fr MsgFrame) error {
	n, err := fr.PackageMessage(w.bw)
	if err != nil {
		return err
	}

	// do not check n, since io.Writer is not allowed to return n < len(buf)
	// without throwing an error
	_, err = w.conf.Writer.Write(w.bw[:n])
	return err
}



func (w *Writer) WriteMsgFrame(fr MsgFrame) error {
	if fr.GetMessage() == nil {
		return fmt.Errorf("message is nil")
	}

	return w.writeFrameInner(fr)
}
