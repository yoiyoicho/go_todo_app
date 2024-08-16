package main

import (
	"context"
	"fmt"
	"log"
	"net"
	"net/http"
	"os"

	"golang.org/x/sync/errgroup"
)

// GolangCI-Lint のチェック用
const Url = "http://example.com"

func main() {
	if len(os.Args) < 2 {
		log.Printf("need a port number\n")
		os.Exit(1)
	}

	p := os.Args[1]
	l, err := net.Listen("tcp", ":"+p)
	if err != nil {
		log.Fatalf("failed to listen port %s: %v", p, err)
	}

	// context.Background() で空のコンテキストを生成
	if err := run(context.Background(), l); err != nil {
		log.Printf("failed to terminate server: %v", err)
		os.Exit(1)
	}
}

func run(ctx context.Context, l net.Listener) error {
	s := &http.Server{
		Handler: http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			fmt.Fprintf(w, "Hello, %s", r.URL.Path[1:])
		}),
	}

	nonuse := "aaa"
	
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
	err := eg.Wait()
	return err
}