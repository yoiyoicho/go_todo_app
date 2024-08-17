package main

import (
	"context"
	"fmt"
	"log"
	"net"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"golang.org/x/sync/errgroup"

	"github.com/yoiyoicho/go_todo_app/config"
)

func main() {
	// context.Background() で空のコンテキストを生成
	if err := run(context.Background()); err != nil {
		log.Printf("failed to terminate server: %v", err)
		os.Exit(1)
	}
}

func run(ctx context.Context) error {
	ctx, stop := signal.NotifyContext(ctx, os.Interrupt, syscall.SIGTERM)
	// run関数内の処理が完了し、エラーハンドリングやリソースの解放などの後処理が行われた後に stop() が実行される
	defer stop()

	cfg, err := config.New()
	if err != nil {
		return err
	}

	l, err := net.Listen("tcp", fmt.Sprintf(":%d", cfg.Port))
	if err != nil {
		log.Fatalf("failed to listen port %d: %v", cfg.Port, err)
	}

	url := fmt.Sprintf("http://%s", l.Addr().String())
	log.Printf("start with: %v", url)

	s := &http.Server{
		Handler: http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			// コマンドラインで実験するため
			time.Sleep(5 * time.Second)
			fmt.Fprintf(w, "Hello, %s", r.URL.Path[1:])
		}),
	}
	
	// errgroup.Groupは複数のゴルーチンを管理し、
	// そのいずれかがエラーを返した場合に全てのゴルーチンをキャンセルする
	eg, ctx := errgroup.WithContext(ctx)
	// 別ゴルーチンでHTTPサーバーを起動する
	eg.Go(func() error {
		// クロージャの特性により、この無名関数は s への参照を保持し、
		// 実行時にその値にアクセスすることができる
		if err := s.Serve(l); err != nil &&
		// http.ErrServerClosed は http.Server.Shutdown() が正常に終了したことを示すので異常ではない
		err != http.ErrServerClosed {
			log.Printf("failed to close: %+v", err)
			return err
		}
		return nil
	})

	// チャネルからの終了通知を待機する
	<-ctx.Done()
	if err := s.Shutdown(context.Background()); err != nil {
		log.Printf("failed to shutdown: %+v", err)
	}

	// eg.Wait() は、errgroup.Group の全てのゴルーチンの実行が完了するまで待機し
	// 最初に発生したエラーを返す
	return eg.Wait()
}