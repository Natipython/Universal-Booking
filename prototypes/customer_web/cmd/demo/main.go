package main

import (
	"context"
	"fmt"
	"net"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	web "ubi.local/customer-web-v0"
)

func main() {
	code := run()
	os.Exit(code)
}
func run() (code int) {
	db, err := web.StartDisposable()
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		return 1
	}
	defer func() {
		if err := db.Close(); err != nil {
			fmt.Fprintln(os.Stderr, err)
			code = 1
		}
	}()
	listener, err := net.Listen("tcp", "127.0.0.1:0")
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		return 1
	}
	origin := "http://" + listener.Addr().String()
	app, err := web.New(origin, web.SQLGateway{Pool: db.Runtime})
	if err != nil {
		_ = listener.Close()
		fmt.Fprintln(os.Stderr, err)
		return 1
	}
	server := &http.Server{Handler: app, ReadHeaderTimeout: 5 * time.Second, ReadTimeout: 10 * time.Second, WriteTimeout: 15 * time.Second, IdleTimeout: 30 * time.Second, MaxHeaderBytes: 8192}
	errors := make(chan error, 1)
	go func() { errors <- server.Serve(listener) }()
	fmt.Println("Local synthetic demo:", origin, "(no token, contact or request-body logs; automatic shutdown in 30 minutes)")
	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer stop()
	select {
	case <-ctx.Done():
	case <-time.After(30 * time.Minute):
	case err := <-errors:
		if err != http.ErrServerClosed {
			fmt.Fprintln(os.Stderr, "demo server stopped")
			return 1
		}
	}
	closeCtx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	_ = server.Shutdown(closeCtx)
	fmt.Println("Synthetic demo stopped; disposable database cleanup follows.")
	return 0
}
