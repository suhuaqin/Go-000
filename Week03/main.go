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
	s := http.Server{}

	g, _ := errgroup.WithContext(context.Background())

	//启动服务
	g.Go(func() error {
		return s.ListenAndServe()
	})

	//监听signal信号
	g.Go(func() error {
		// 用quit监听signal信号
		quit := make(chan os.Signal)
		signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
		<-quit

		// 使用超时 ctx，让服务能在最后 5S 处理 shutdown 前进来的请求
		// 为了避免集联取消，这里特地新生成一下root ctx
		ctx, cancle := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancle()
		return s.Shutdown(ctx)
	})

	err := g.Wait()
	fmt.Println(err)
}
