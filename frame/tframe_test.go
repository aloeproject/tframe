package frame

import (
	"context"
	"github.com/aloeproject/gcommon/logger"
	"testing"
)

var log logger.ILogger

func TestMain(m *testing.M) {
	log = logger.NewZLogger(logger.NewLogger("name", "./", "logs"))
	m.Run()
}

func TestTFrame_Start(t1 *testing.T) {
	ser := NewTFrame(log)
	ser.Start(context.Background())
}
