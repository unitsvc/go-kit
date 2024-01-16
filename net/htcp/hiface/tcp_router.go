package hiface

// ITcpRouter 路由接口，自定义处理业务方法
type ITcpRouter interface {
	// PreHandle 在处理conn业务之前的钩子方法
	PreHandle(request ITcpRequest)
	// Handle 处理conn业务的方法
	Handle(request ITcpRequest)
	// PostHandle 处理conn业务之后的钩子方法
	PostHandle(request ITcpRequest)
}
