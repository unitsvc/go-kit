package hiface

// ITcpPkgHandle 消息管理抽象层
type ITcpPkgHandle interface {
	// DoPkgHandler 马上以非阻塞方式处理数据包
	DoPkgHandler(request ITcpRequest)
	// AddRouter 为数据包添加具体的处理逻辑
	AddRouter(handlerRouter string, router ITcpRouter)
	// StartWorkerPool 启动worker工作池
	StartWorkerPool()
	// SendPkgToTaskQueue 将数据包交给TaskQueue，由worker进行处理
	SendPkgToTaskQueue(request ITcpRequest)
}
