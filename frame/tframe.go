package frame

import (
	"context"
	"github.com/aloeproject/tframe/iface"
	"github.com/aloeproject/tframe/tnet"
	"github.com/aloeproject/toolbox/logger"
	"net"
	"time"
)

type option struct {
	Addr           string //格式 127.0.0.1:8080
	TimeOut        time.Duration
	MaxConnections int    //最大连接数
	WorkerPoolSize uint32 //work工作池数
	WorkTaskLen    uint32 //任务队列长度
}

type Option func(opt *option)

func WithAddrTFrame(s string) Option {
	return func(opt *option) {
		opt.Addr = s
	}
}

func WithAddrTimeOut(t time.Duration) Option {
	return func(opt *option) {
		opt.TimeOut = t
	}
}

func WithMaxConnections(maxConnections int) Option {
	return func(opt *option) {
		opt.MaxConnections = maxConnections
	}
}

func WithWorkerConf(poolSize, taskLen uint32) Option {
	return func(opt *option) {
		opt.WorkerPoolSize = poolSize
		opt.WorkTaskLen = taskLen
	}
}

func NewTFrame(logger logger.ILogger, opts ...Option) *TFrame {
	defaultOpt := option{
		Addr:           "127.0.0.1:8990",
		TimeOut:        30 * time.Second,
		MaxConnections: 1000,
		WorkerPoolSize: 10,
		WorkTaskLen:    5,
	}

	for _, o := range opts {
		o(&defaultOpt)
	}

	return &TFrame{
		option: option{
			Addr:           defaultOpt.Addr,
			TimeOut:        defaultOpt.TimeOut,
			MaxConnections: defaultOpt.MaxConnections,
			WorkerPoolSize: defaultOpt.WorkerPoolSize,
			WorkTaskLen:    defaultOpt.WorkTaskLen,
		},
		log:         logger,
		connManager: tnet.MangerObj,
		router:      &tnet.DefaultRouter{}, //默认走default
	}
}

var _ iface.IServer = (*TFrame)(nil)

type TFrame struct {
	option
	ctx    context.Context
	cancel context.CancelFunc

	log logger.ILogger

	lister *net.TCPListener

	onConnStart func(iface.IConnection)
	onConnStop  func(iface.IConnection)

	connManager iface.IConnManager
	router      iface.IRouter
}

func (t *TFrame) AddRouter(router iface.IRouter) {
	t.router = router
}

func (t *TFrame) SetOnConnStart(f func(iface.IConnection)) {
	t.onConnStart = f
}

func (t *TFrame) SetOnConnStop(f func(iface.IConnection)) {
	t.onConnStop = f
}

func (t *TFrame) CallOnConnStart(conn iface.IConnection) {
	if t.onConnStart != nil {
		t.onConnStart(conn)
	}
}

func (t *TFrame) CallOnConnStop(conn iface.IConnection) {
	if t.onConnStop != nil {
		t.onConnStop(conn)
	}
}

func (t *TFrame) start() error {
	t.log.Infof("tframe start address:[%v]", t.Addr)

	addr, err := net.ResolveTCPAddr("tcp", t.Addr)
	if err != nil {
		return err
	}

	t.lister, err = net.ListenTCP("tcp", addr)
	if err != nil {
		return err
	}

	go func() {
		for {
			conn, err := t.lister.AcceptTCP()
			if err != nil {
				t.log.Errorw(t.ctx, "TFrame_start_Accept err:%v", err)
				return
			}

			//超过最大连接
			if t.connManager.Len() >= t.MaxConnections {
				conn.Close()
				t.log.Errorw(t.ctx, "TFrame_NewConnection Maximum connections exceeded")
				continue
			}

			handleManager := tnet.NewHandleMsgManager(t.ctx, t.log, t.option.WorkerPoolSize, t.option.WorkTaskLen)
			handleManager.StartWorkPool()

			connSer := tnet.NewConnection(t.ctx, conn, t, t.router, handleManager, t.log)
			go func() {
				err = connSer.Start()
				if err != nil {
					t.log.Errorw(t.ctx, "TFrame_NewConnection_Start err:%v", err)
					return
				}
			}()
		}
	}()

	return nil
}

// 框架启动
func (t *TFrame) Start(ctx context.Context) error {
	t.ctx, t.cancel = context.WithCancel(ctx)

	err := t.start()
	if err != nil {
		t.log.Errorw(ctx, "TFrame_Start error err:%v", err)
	}

	select {
	case <-t.ctx.Done():
		return nil
	}
}

// 框架停止
func (t *TFrame) Stop(ctx context.Context) error {
	t.connManager.ClearConn()
	if t.lister != nil {
		err := t.lister.Close()
		if err != nil {
			t.log.Errorw(t.ctx, "TFrame_start_lister_Close err:%v", err)
		}
	}
	t.cancel()
	return nil
}
