package hiface

type ITcpRequest interface {
	// GetConnection 获取请求连接信息
	GetConnection() ITcpConnection
	// GetPkg 获取请求数据包
	GetPkg() []byte
	// GetPkgBody 获取请求数据（场景：1、解密后的数据包内容 2、真正的消息体）
	GetPkgBody() []byte
	// GetHandlerRouter 获取请求路由
	GetHandlerRouter() string
}
