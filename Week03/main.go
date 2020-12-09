package main

import (
	"context"
	"fmt"
	"golang.org/x/sync/errgroup"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"
)

func main() {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	g, _ := errgroup.WithContext(ctx)

	s := http.Server{}
	//启动服务
	g.Go(func() error {
		err := s.ListenAndServe()
		if err != nil && err != http.ErrServerClosed {
			//其他原因退出，通过ctx让监听signal的gourtine退出
			cancel()
		}
		return err
	})

	//监听signal信号
	g.Go(func() error {
		// 用quit监听signal信号
		quit := make(chan os.Signal)
		signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
		select {
		case <-quit:
			// 使用超时ctx，让服务能在最后秒处理shutdown前进来的请求
			ctx1, cancle1 := context.WithTimeout(ctx, 5*time.Second)
			defer cancle1()
			return s.Shutdown(ctx1)
		case <-ctx.Done():
			return nil
		}
	})

	fmt.Println(g.Wait())
}
