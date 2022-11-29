package tnet

import (
	"context"
	"github.com/aloeproject/gcommon/logger"
	"github.com/aloeproject/tframe/iface"
	"sync/atomic"
)

var _ iface.IHandleManger = (*HandleMsgManager)(nil)

type workPool struct {
	ctx    context.Context
	log    logger.ILogger
	router iface.IRouter
	req    iface.IRequest
}

type HandleMsgManager struct {
	ctx context.Context
	log logger.ILogger

	poolSize    uint32
	workTaskLen uint32
	workPool    []chan *workPool
	close       chan struct{}

	count uint32
}

func NewHandleMsgManager(ctx context.Context, log logger.ILogger, poolSize, workTaskLen uint32) *HandleMsgManager {
	return &HandleMsgManager{
		ctx:         ctx,
		log:         log,
		poolSize:    poolSize,
		workTaskLen: workTaskLen,
		workPool:    make([]chan *workPool, poolSize),
		close:       make(chan struct{}),
	}
}

func (h *HandleMsgManager) getPoolSize() uint32 {
	return h.poolSize
}

func (h *HandleMsgManager) StartWorkPool() {
	var i uint32
	for ; i < h.poolSize; i++ {
		h.workPool[i] = make(chan *workPool, h.workTaskLen)
		go h.runWork(h.workPool[i])
	}
}

/*
	工作退出
*/

func (h HandleMsgManager) workExit() {
	if err := recover(); err != nil {
		h.log.Errorw(h.ctx, "runWork recover error:[%+v]", err)
		return
	}
}

func (h *HandleMsgManager) runWork(work chan *workPool) {
	defer h.workExit()

	for {
		select {
		case <-h.close:
			return
		case <-h.ctx.Done():
			return
		case v := <-work:
			v.router.Handle(v.ctx, v.req)
		}
	}

}

func (h *HandleMsgManager) AutoHandle(ctx context.Context, router iface.IRouter, req iface.IRequest) {
	if h.getPoolSize() > 0 {
		h.AsyncHandle(ctx, router, req)
	} else {
		h.SynHandle(ctx, router, req)
	}
}

func (h *HandleMsgManager) SynHandle(ctx context.Context, router iface.IRouter, req iface.IRequest) {
	defer h.workExit()

	router.Handle(ctx, req)
}

func (h *HandleMsgManager) AsyncHandle(ctx context.Context, router iface.IRouter, req iface.IRequest) {
	//请求数量
	num := atomic.AddUint32(&h.count, 1)
	mod := num % h.poolSize
	h.workPool[mod] <- &workPool{
		ctx:    ctx,
		router: router,
		req:    req,
	}
	//防止溢出，到一亿次重置下
	if num > 1000*1000*100 {
		atomic.StoreUint32(&h.count, 0)
	}
}

func (h *HandleMsgManager) Close() {
	close(h.close)
}
