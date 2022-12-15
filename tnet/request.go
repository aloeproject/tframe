package tnet

import (
	"context"
	"fmt"
	"github.com/aloeproject/tframe/iface"
	"net"
)

/*
这里的路由的数据结构是
{
	"method":方法
	.....
}
*/

var _ iface.IRouter = (*DefaultRouter)(nil)

type DefaultRouter struct {
}

func (d *DefaultRouter) GetPacket() iface.IPacket {
	return &DefaultPack{}
}

func (d *DefaultRouter) Handle(ctx context.Context, req iface.IRequest) {
	str := req.GetData()
	req.SendMsg([]byte(fmt.Sprintf("ok,收到消息%s", str)))
}

func (d *DefaultRouter) GetRequest(conn iface.IQConnection, msg iface.IRMessage) iface.IRequest {
	return &DefaultRequest{
		Conn: conn,
		Msg:  msg,
	}
}

var _ iface.IRequest = (*DefaultRequest)(nil)

func NewDefaultRequest(conn iface.IQConnection,
	msg iface.IRMessage) *DefaultRequest {
	return &DefaultRequest{
		Conn: conn,
		Msg:  msg,
	}
}

type DefaultRequest struct {
	Conn iface.IQConnection
	Msg  iface.IRMessage
}

func (d *DefaultRequest) GetMID() int32 {
	return d.Msg.GetMID()
}

func (d *DefaultRequest) GetHeadData() []byte {
	return d.Msg.GetHeadData()
}

func (d *DefaultRequest) SendMsg(bytes []byte) {
	d.Conn.SendMsg(bytes)
}

func (d *DefaultRequest) GetTCPConnection() *net.TCPConn {
	return d.Conn.GetTCPConnection()
}

func (d *DefaultRequest) GetConnId() string {
	return d.Conn.GetConnId()
}

func (d *DefaultRequest) SetProperty(k string, v interface{}) {
	d.Conn.SetProperty(k, v)
}

func (d *DefaultRequest) GetProperty(k string) (v interface{}, exists bool) {
	return d.Conn.GetProperty(k)
}

func (d *DefaultRequest) DelProperty(k string) {
	d.Conn.DelProperty(k)
}

func (d *DefaultRequest) GetDataLen() int {
	return d.Msg.GetDataLen()
}

func (d *DefaultRequest) GetData() []byte {
	return d.Msg.GetData()
}
