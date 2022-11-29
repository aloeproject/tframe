package iface

import (
	"context"
	"net"
)

type IServer interface {
	Start(ctx context.Context) error
	Stop(ctx context.Context) error

	AddRouter(IRouter)

	SetOnConnStart(func(IConnection)) //设置该Server的连接创建时Hook函数
	SetOnConnStop(func(IConnection))  //设置该Server的连接断开时的Hook函数
	CallOnConnStart(conn IConnection)
	CallOnConnStop(conn IConnection)
}

//请求中使用的连接方法
type IQConnection interface {
	SendMsg([]byte)
	GetTCPConnection() *net.TCPConn
	GetConnId() string
	SetProperty(k string, v interface{})
	GetProperty(k string) (v interface{}, exists bool)
	DelProperty(k string)
}

type IConnection interface {
	IQConnection
	Start() error
	Stop()
}

type IConnManager interface {
	Add(IConnection)
	Close(IConnection)
	Get(string) IConnection
	Len() int
	ClearConn()
}

/*
读取消息
*/
type IRMessage interface {
	GetDataLen() int //获取消息数据段长度
	GetData() []byte //获取消息内容

	GetHeadData() []byte
}

/*
	将请求的一个消息封装到message中，定义抽象层接口
*/
type IMessage interface {
	IRMessage
	SetHeadData([]byte) //设置head数据
	SetData([]byte)     //设计消息内容
	SetHeadLen(int)
}

/*
数据拆包
*/
type IPacket interface {
	Unpack(binaryData []byte) (IMessage, error)
	Pack(msg IMessage) ([]byte, error)
	GetHeadLen() int
	GetMaxDataLen() int //获取一次请求的最大数据长度 非拆包使用
}

type IRequest interface {
	IQConnection
	IRMessage
}

type IRouter interface {
	Handle(ctx context.Context, req IRequest)
	GetRequest(conn IQConnection, msg IRMessage) IRequest
	GetPacket() IPacket //获取数据打包方式
}

/*
  消息管理
*/
type IHandleManger interface {
	StartWorkPool() //启动work 工作池
	//自动同步异步选择
	AutoHandle(ctx context.Context, router IRouter, req IRequest)
	//同步处理
	SynHandle(ctx context.Context, router IRouter, req IRequest)
	//异步处理
	AsyncHandle(ctx context.Context, router IRouter, req IRequest)
	//关闭连接管理
	Close()
}
