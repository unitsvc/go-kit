package hnet

import (
	"github.com/unitsvc/go-kit/net/htcp/hiface"
)

type TcpRequest struct {
	conn hiface.ITcpConnection // 已经和客户端建立的连接
	pkg  hiface.ITcpPkg        // 客户端请求的数据包
}

// GetConnection 获取请求连接信息
func (r *TcpRequest) GetConnection() hiface.ITcpConnection {
	return r.conn
}

// GetPkg 获取请求消息的数据包
func (r *TcpRequest) GetPkg() []byte {
	return r.pkg.GetPkg()
}

// GetPkgBody 获取请求消息的数据包
func (r *TcpRequest) GetPkgBody() []byte {
	return r.pkg.GetPkgBody()
}

// GetHandlerRouter 获取请求的路由
func (r *TcpRequest) GetHandlerRouter() string {
	return r.pkg.GetHandlerRouter()
}
