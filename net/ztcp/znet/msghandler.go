package znet

import (
	"fmt"
	"strconv"

	"github.com/unitsvc/go-kit/net/ztcp/ziface"
	"github.com/unitsvc/go-kit/net/ztcp/zutils"
)

type MsgHandle struct {
	Apis           map[uint32]ziface.IRouter // 存放每个MsgId 所对应的处理方法的map属性
	WorkerPoolSize uint32                    // 业务工作Worker池的数量
	TaskQueue      []chan ziface.IRequest    // Worker负责取任务的消息队列
}

func NewMsgHandle() *MsgHandle {
	return &MsgHandle{
		Apis:           make(map[uint32]ziface.IRouter),
		WorkerPoolSize: zutils.GlobalObject.WorkerPoolSize,
		// 一个worker对应一个queue
		TaskQueue: make([]chan ziface.IRequest, zutils.GlobalObject.WorkerPoolSize),
	}
}

// SendMsgToTaskQueue 将消息交给TaskQueue，由worker进行处理
func (mh *MsgHandle) SendMsgToTaskQueue(request ziface.IRequest) {
	// 根据ConnID来分配当前的连接应该由哪个worker负责处理
	// 轮询的平均分配法则

	// 得到需要处理此条连接的workerID
	workerID := request.GetConnection().GetConnID() % mh.WorkerPoolSize
	// fmt.Println("添加 ConnID=", request.GetConnection().GetConnID()," 请求 msgID=", request.GetMsgID(), "到 workerID=", workerID)
	// 将请求消息发送给任务队列
	mh.TaskQueue[workerID] <- request
}

// DoMsgHandler 马上以非阻塞方式处理消息
func (mh *MsgHandle) DoMsgHandler(request ziface.IRequest) {
	handler, ok := mh.Apis[request.GetMsgID()]
	if !ok {
		fmt.Println("接口路由 msgId = ", request.GetMsgID(), "没有发现！")
		return
	}

	// 执行对应处理方法
	handler.PreHandle(request)
	handler.Handle(request)
	handler.PostHandle(request)
}

// AddRouter 为消息添加具体的处理逻辑
func (mh *MsgHandle) AddRouter(msgId uint32, router ziface.IRouter) {
	// 1、判断当前msg绑定的API处理方法是否已经存在
	if _, ok := mh.Apis[msgId]; ok {
		panic("接口路由重复，msgId = " + strconv.Itoa(int(msgId)))
	}
	// 2、添加msg与api的绑定关系
	mh.Apis[msgId] = router
	fmt.Println("[Gf-Plus] 添加接口路由 msgId = ", msgId)
}

// StartOneWorker 启动一个Worker工作协程
func (mh *MsgHandle) StartOneWorker(workerID int, taskQueue chan ziface.IRequest) {
	fmt.Println("[启动] 工作协程 workerID = ", workerID)
	// 不断的等待队列中的消息
	for {
		select {
		// 有消息则取出队列的Request，并执行绑定的业务方法
		case request := <-taskQueue:
			mh.DoMsgHandler(request)
		}
	}
}

// StartWorkerPool 启动worker工作池
func (mh *MsgHandle) StartWorkerPool() {
	// 遍历需要启动worker的数量，依此启动
	for i := 0; i < int(mh.WorkerPoolSize); i++ {
		// 一个worker被启动
		// 给当前worker对应的任务队列开辟空间
		mh.TaskQueue[i] = make(chan ziface.IRequest, zutils.GlobalObject.MaxWorkerTaskLen)
		// 启动当前Worker，阻塞的等待对应的任务队列是否有消息传递进来
		go mh.StartOneWorker(i, mh.TaskQueue[i])
	}
}
