package main

import (
	"context"
	"fmt"
	"github.com/aloeproject/gcommon/logger"
	"github.com/aloeproject/tframe/app_demo/jsonapp/router"
	"github.com/aloeproject/tframe/frame"
	"runtime"
	"sync"
	"time"
)

func main() {
	wg := new(sync.WaitGroup)
	wg.Add(2)

	log := logger.NewZLogger(
		logger.NewLogger("name", "./", "logs"),
	)
	ser := frame.NewTFrame(log, frame.WithAddrTFrame("127.0.0.1:8088"), frame.WithMaxConnections(10), frame.WithWorkerConf(50, 5))
	ser.AddRouter(router.NewRouter(log))

	fmt.Printf("num:%d\n", runtime.NumGoroutine())

	go func() {
		defer wg.Done()
		ser.Start(context.Background())
	}()

	go func() {
		defer wg.Done()
		time.Sleep(15000 * time.Second)
		ser.Stop(context.Background())
	}()

	go func() {
		for {
			time.Sleep(500 * time.Millisecond)
			fmt.Printf("num:%d\n", runtime.NumGoroutine())
		}
	}()

	wg.Wait()
	for {
		fmt.Printf("num:%d\n", runtime.NumGoroutine())
		time.Sleep(1 * time.Second)
	}
	select {}
}
