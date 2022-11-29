package router

import (
	"context"
	"encoding/json"
	"github.com/aloeproject/gcommon/logger"
	"github.com/aloeproject/tframe/app_demo/jsonapp/controller"
	"github.com/aloeproject/tframe/iface"
	"github.com/aloeproject/tframe/tnet"
)

var _ iface.IRouter = (*Router)(nil)

type Action struct {
	Method string `json:"method"`
}

func NewRouter(log logger.ILogger) *Router {
	return &Router{
		log: log,
	}
}

type Router struct {
	log logger.ILogger
}

func (r Router) Handle(ctx context.Context, req iface.IRequest) {
	body := req.GetData()
	action := Action{}
	err := json.Unmarshal(body, &action)
	if err != nil {
		r.log.Errorw(ctx, "action Unmarshal error:%v", err)
		return
	}
	var contr controller.IController
	if action.Method == "obu.upload" {
		contr = controller.NewUploadController(r.log)
	}

	if contr != nil {
		contr.Action(ctx, req)
	}
}

func (r Router) GetRequest(conn iface.IQConnection, msg iface.IRMessage) iface.IRequest {
	return tnet.NewDefaultRequest(conn, msg)
}

func (r Router) GetPacket() iface.IPacket {
	return &tnet.OnceQueryPack{}
}
