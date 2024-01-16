package hiface

import "net"

type ITcpConnection interface {
	// Start 启动连接，让当前连接开始工作
	Start()
	// Stop 停止连接，结束当前连接状态
	Stop()
	// GetTCPConnection 从当前连接获取原始的Socket TCPConn
	GetTCPConnection() *net.TCPConn
	// GetConnID 获取当前连接ID（会话ID）
	GetConnID() int64
	// RemoteAddr 获取远程客户端地址信息
	RemoteAddr() net.Addr
	// SetProperty 设置连接属性
	SetProperty(key string, value interface{})
	// GetProperty 获取连接属性
	GetProperty(key string) (interface{}, error)
	// RemoveProperty 移除连接属性
	RemoveProperty(key string)
	// SendTcpPkg 发送tcp数据包
	//  @pkgBodyType 数据类型
	//  @pkg         数据内容（自定义结构体对象）
	//  @userBuf     是否启动缓冲（默认关闭）
	SendTcpPkg(pkgBodyType byte, pkg interface{}, userBuf ...bool) error
}
