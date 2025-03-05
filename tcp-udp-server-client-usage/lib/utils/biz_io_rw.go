package utils

import "io"

type BizIoWRDeadliner interface {
	SetReadDeadline() error
	SetWriteDeadline() error
}

type BizIoWRWrapper interface {
	BizIoWRDeadliner
	io.ReadWriteCloser
}
