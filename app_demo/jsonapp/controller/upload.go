package controller

import (
	"context"
	"github.com/aloeproject/gcommon/logger"
	"github.com/aloeproject/tframe/iface"
	"time"
)

type IController interface {
	Action(ctx context.Context, req iface.IRequest)
}

var _ IController = (*UploadController)(nil)

func NewUploadController(log logger.ILogger) *UploadController {
	return &UploadController{log: log}
}

type UploadController struct {
	log logger.ILogger
}

func (u UploadController) Action(ctx context.Context, req iface.IRequest) {
	u.log.Debugw(ctx, "UploadController action %s", req.GetConnId())
	time.Sleep(1 * time.Second)
	req.SendMsg([]byte("收到数据"))
}
