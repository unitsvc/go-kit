package hiface

import "io"

type ITcpPkg interface {

	// GetPkgHeadLen 获取数据包消息头长度
	GetPkgHeadLen() int16
	// GetPkg 获取请求数据包
	GetPkg() []byte
	// GetPkgBody 获取数据包消息体内容
	GetPkgBody() []byte
	// GetHandlerRouter 获取路由（即消息处理器）
	GetHandlerRouter() string
	// Pack 封包方法
	Pack() ([]byte, error)
	// Unpack 拆包方法
	Unpack(conn io.Reader) (ITcpPkg, error)
}
