package main

import (
	"context"
	"fmt"
	"io"
	"net"
	"net/http"
	"testing"

	"golang.org/x/sync/errgroup"
)

func TestRun(t *testing.T) {
	t.Skip("リファクタリング中")

	// ポート番号に0を指定すると利用可能なポートを動的に選択してくれる
	l, err := net.Listen("tcp", "localhost:0")
	if err != nil {
		t.Fatalf("failed to listen port %v", err)
	}
	ctx, cancel := context.WithCancel(context.Background())
	eg, ctx := errgroup.WithContext(ctx)
	eg.Go(func() error {
		return run(ctx)
	})
	in := "message"
	url := fmt.Sprintf("http://%s/%s", l.Addr().String(), in)
	t.Logf("try request to %q", url)
	rsp, err := http.Get(url)
	if err != nil {
		t.Errorf("failed to get: %+v", err)
	}

	got, err := io.ReadAll(rsp.Body)
	defer rsp.Body.Close()
	if err != nil {
		t.Errorf("failed to read body: %v", err)
	}

	want := fmt.Sprintf("Hello, %s", in)
	// 期待通りにHTTPサーバーが起動しているかテストする
	if string(got) != want {
		t.Errorf("want %q, but got %q", want, got)
	}

	cancel()
	// テストコードから意図通りに終了するかテストする
	if err := eg.Wait(); err != nil {
		// Fatalfはテストを失敗としてマークし、そのテストの実行を中止する
		t.Fatal(err)
	}
}
