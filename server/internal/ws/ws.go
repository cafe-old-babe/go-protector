package ws

import (
	"go-protector/server/internal/base"
	"io"
)

type IWsWriter interface {
	Write(msg *base.WsMsg) error
}

type IWsWriteCloser interface {
	IWsWriter
	io.Closer
}

type ICtxWsWriter interface {
	base.IService
	IWsWriter
}

type ICtxWsWriteCloser interface {
	base.IService
	IWsWriteCloser
}
