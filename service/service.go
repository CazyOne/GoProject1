package service

import (
	"context"
	"distributed-learn/registry"
	"fmt"
	"log"
	"net/http"
)

//注册HTTP处理器，启动HTTP服务器，向注册中心注册服务
func Start(ctx context.Context, host, port string, reg registry.Registration,
	registerHandlersFunc func()) (context.Context, error) {
	registerHandlersFunc()
	ctx = startService(ctx, reg.ServiceName, host, port)
	err := registry.RegisterService(reg)
	if err != nil {
		return ctx, err
	}

	return ctx, nil
}

func startService(ctx context.Context, serviceName registry.ServiceName,
	host, port string) context.Context {
	ctx, cancel := context.WithCancel(ctx)

	var srv http.Server
	srv.Addr = ":" + port

	//启动HTTP服务器，当服务器停止时，自动从注册中心注销服务
	go func() {
		log.Println(srv.ListenAndServe())
		err := registry.ShutdownService(fmt.Sprintf("http://%s:%s", host, port))
		if err != nil {
			log.Println(err)
		}
		cancel()
	}()

	//提示用户服务已启动，等待用户输入
	go func() {
		fmt.Printf("%v started. Press any key to stop\n", serviceName)
		var s string
		fmt.Scanln(&s) // 会等待用户输入，如果输入了就继续往下走，来Shutdown
		err := registry.ShutdownService(fmt.Sprintf("http://%s:%s", host, port))
		if err != nil {
			log.Println(err)
		}
		srv.Shutdown(ctx)
		cancel()
	}()

	return ctx
}
