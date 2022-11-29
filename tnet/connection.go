package tnet

import (
	"context"
	"github.com/aloeproject/gcommon/logger"
	"github.com/aloeproject/tframe/iface"
	"github.com/google/uuid"
	"io"
	"net"
	"sync"
	"sync/atomic"
	"time"
)

var _ iface.IConnection = (*Connection)(nil)

type connOption struct {
	TimeOut time.Duration
}

type ConnOption func(opt *connOption)

func NewConnection(ctx context.Context, conn *net.TCPConn,
	server iface.IServer, router iface.IRouter, handleManger iface.IHandleManger,
	logger logger.ILogger, opts ...ConnOption) *Connection {

	defaultOpt := connOption{
		TimeOut: 30 * time.Second,
	}

	for _, o := range opts {
		o(&defaultOpt)
	}
	obj := &Connection{
		ctx:          ctx,
		logger:       logger,
		connId:       uuid.New().String(),
		conn:         conn,
		property:     new(sync.Map),
		server:       server,
		router:       router,
		manager:      MangerObj,
		dataPack:     router.GetPacket(),
		timeOut:      defaultOpt.TimeOut,
		handleManger: handleManger,
		msgChan:      make(chan []byte),
		isClose:      0,
	}

	MangerObj.Add(obj)
	return obj
}

type Connection struct {
	ctx    context.Context
	logger logger.ILogger
	cancel context.CancelFunc
	connId string

	conn     *net.TCPConn
	property *sync.Map

	server       iface.IServer
	router       iface.IRouter
	manager      iface.IConnManager
	dataPack     iface.IPacket
	handleManger iface.IHandleManger

	timeOut time.Duration
	//数据写入通道
	msgChan chan []byte
	isClose int32 //0 正常 1 关闭
}

func (c *Connection) GetConnId() string {
	return c.connId
}

func (c *Connection) GetTCPConnection() *net.TCPConn {
	return c.conn
}

func (c *Connection) Stop() {
	if c.cancel != nil {
		c.cancel()
	}
}

func (c *Connection) Start() error {
	c.ctx, c.cancel = context.WithCancel(c.ctx)
	c.logger.Infow(c.ctx, "%v connection is create,current connection nums:[%v]", c.GetConnId(), c.manager.Len())
	//连接时回调函数
	c.server.CallOnConnStart(c)

	go c.read(c.ctx)
	go c.write(c.ctx)

	select {
	case <-c.ctx.Done():
		defer c.closer(c.ctx)
		if err := c.ctx.Err(); err != nil && err != context.Canceled {
			return err
		}
		c.logger.Infow(c.ctx, "%s connection is close", c.GetConnId())
	}

	return nil
}

func (c *Connection) SendMsg(bytes []byte) {
	//已关闭不能发送
	if atomic.LoadInt32(&c.isClose) == 1 {
		return
	}
	c.msgChan <- bytes
}

func (c *Connection) SetProperty(k string, v interface{}) {
	c.property.Store(k, v)
}

func (c *Connection) GetProperty(k string) (v interface{}, exists bool) {
	return c.property.Load(k)
}

func (c *Connection) DelProperty(k string) {
	c.property.Delete(k)
}

func (c *Connection) read(ctx context.Context) {
	defer c.Stop()

	for {
		select {
		case <-ctx.Done():
			return
		default:
			//这种不需要拆包
			var err error
			if c.dataPack.GetHeadLen() == 0 {
				data := make([]byte, c.dataPack.GetMaxDataLen())
				var readSize int
				if readSize, err = c.conn.Read(data); err != nil {
					if err != io.EOF {
						c.logger.Errorw(ctx, "Connection_read_data_ReadFull error:[%v]", err)
					}
					return
				}
				msg, err := c.dataPack.Unpack(data[:readSize])
				if err != nil {
					c.logger.Errorw(ctx, "Connection_read_data_Unpack error:[%v]", err)
					continue
				}
				msg.SetData(data[:readSize])
				c.handleManger.AutoHandle(c.ctx, c.router, c.router.GetRequest(c, msg))
			} else {
				//进行第一次拆包
				headData := make([]byte, c.dataPack.GetHeadLen())
				if _, err := c.conn.Read(headData); err != nil {
					if err != io.EOF {
						c.logger.Errorw(ctx, "Connection_read_head_ReadFull error:[%v]", err)
					}
					return
				}
				//头拆包
				msg, err := c.dataPack.Unpack(headData)
				if err != nil {
					c.logger.Errorw(ctx, "Connection_read_head_Unpack error:[%v]", err)
					return
				}
				//第二次拆包
				if msg.GetDataLen() > 0 {
					data := make([]byte, msg.GetDataLen())
					if _, err = c.conn.Read(data); err != nil {
						if err != io.EOF {
							c.logger.Errorw(ctx, "Connection_read_data_ReadFull error:[%v]", err)
						}
						return
					}
					msg.SetData(data)
					c.handleManger.AutoHandle(c.ctx, c.router, c.router.GetRequest(c, msg))
				}
			}
		}
	}
}

func (c *Connection) write(ctx context.Context) {
	for {
		select {
		case data := <-c.msgChan:
			if _, err := c.conn.Write(data); err != nil {
				c.logger.Errorw(ctx, "Connection_write_Write error:[%v]", err)
				continue
			}
		case <-ctx.Done():
			return
		}
	}
}

func (c *Connection) closer(ctx context.Context) {
	c.server.CallOnConnStop(c)
	//管理连接关闭
	c.manager.Close(c)
	//handle管理关闭
	c.handleManger.Close()
	//设置关闭状态
	atomic.StoreInt32(&c.isClose, 1)

	err := c.conn.Close()
	if err != nil {
		c.logger.Errorw(ctx, "Connection_closer_Close error:[%v]", err)
	}
	close(c.msgChan)
}
