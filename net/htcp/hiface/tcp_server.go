package hiface

// ITcpServer 定义服务器接口
type ITcpServer interface {
	// Start 启动服务器方法
	Start()
	// Stop 停止服务器方法
	Stop()
	// Serve 开启业务服务方法
	Serve()
	// AddRouter 路由功能：给当前服务注册一个路由业务方法，供客户端连接处理使用
	AddRouter(handlerRouter string, router ITcpRouter)
	// GetConnMgr 获取连接管理
	GetConnMgr() ITcpConnManager
	// SetOnConnStart 设置该Server的连接创建时Hook函数
	SetOnConnStart(func(ITcpConnection))
	// SetOnConnStop 设置该Server的连接断开时的Hook函数
	SetOnConnStop(func(ITcpConnection))
	// CallOnConnStart 调用连接OnConnStart Hook函数
	CallOnConnStart(conn ITcpConnection)
	// CallOnConnStop 调用连接OnConnStop Hook函数
	CallOnConnStop(conn ITcpConnection)
}
