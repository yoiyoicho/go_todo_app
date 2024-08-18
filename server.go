package main

import (
	"context"
	"log"
	"net"
	"net/http"
	"os"
	"os/signal"
	"syscall"

	"golang.org/x/sync/errgroup"
)

type Server struct {
	srv *http.Server
	l   net.Listener
}

func NewServer(l net.Listener, mux http.Handler) *Server {
	return &Server{
		srv: &http.Server{Handler: mux},
		l:   l,
	}
}

func (s *Server) Run(ctx context.Context) error {
	ctx, stop := signal.NotifyContext(ctx, os.Interrupt, syscall.SIGTERM)
	// run関数内の処理が完了し、エラーハンドリングやリソースの解放などの後処理が行われた後に stop() が実行される
	defer stop()

	// errgroup.Groupは複数のゴルーチンを管理し、
	// そのいずれかがエラーを返した場合に全てのゴルーチンをキャンセルする
	eg, ctx := errgroup.WithContext(ctx)
	// 別ゴルーチンでHTTPサーバーを起動する
	eg.Go(func() error {
		// クロージャの特性により、この無名関数は s への参照を保持し、
		// 実行時にその値にアクセスすることができる
		if err := s.srv.Serve(s.l); err != nil &&
			// http.ErrServerClosed は http.Server.Shutdown() が正常に終了したことを示すので異常ではない
			err != http.ErrServerClosed {
			log.Printf("failed to close: %+v", err)
			return err
		}
		return nil
	})

	// チャネルからの終了通知を待機する
	<-ctx.Done()
	if err := s.srv.Shutdown(context.Background()); err != nil {
		log.Printf("failed to shutdown: %+v", err)
	}

	// eg.Wait() は、errgroup.Group の全てのゴルーチンの実行が完了するまで待機し
	// 最初に発生したエラーを返す
	return eg.Wait()
}
