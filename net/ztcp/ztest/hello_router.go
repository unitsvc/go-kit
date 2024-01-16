package ztest

import (
	"github.com/unitsvc/go-kit/net/ztcp/ziface"
	"github.com/unitsvc/go-kit/net/ztcp/zlog"
	"github.com/unitsvc/go-kit/net/ztcp/znet"
)

type HelloRouter struct {
	znet.BaseRouter
}

// Handle HelloZinxRouter 路由
func (this *HelloRouter) Handle(request ziface.IRequest) {
	zlog.Debug("执行路由")

	// 先读取客户端的数据，再回写ping...ping...ping
	zlog.Debug("接收客户端消息：msgId=", request.GetMsgID(), "，data=", string(request.GetData()))

	err := request.GetConnection().SendBuffMsg(1, []byte("Hello gf-plus"))
	if err != nil {
		zlog.Error(err)
	}
}
