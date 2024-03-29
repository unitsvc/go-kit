package ztest

import (
	"github.com/unitsvc/go-kit/net/ztcp/ziface"
	"github.com/unitsvc/go-kit/net/ztcp/zlog"
	"github.com/unitsvc/go-kit/net/ztcp/znet"
)

// PingRouter 自定义路由
type PingRouter struct {
	znet.BaseRouter
}

// Handle Ping处理器
func (this *PingRouter) Handle(request ziface.IRequest) {

	zlog.Debug("执行ping处理器")
	// 先读取客户端的数据，再回写ping...ping...ping
	zlog.Debug("接收客户端消息：msgId=", request.GetMsgID(), "，data=", string(request.GetData()))

	err := request.GetConnection().SendBuffMsg(0, []byte("ping"))
	if err != nil {
		zlog.Error(err)
	}
}
