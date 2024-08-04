package main

import (
	"context"
	"fmt"
	"io"
	"net/http"
	"testing"

	"golang.org/x/sync/errgroup"
)

func TestRun(t *testing.T) {
	ctx, cancel := context.WithCancel(context.Background())
	eg, ctx := errgroup.WithContext(ctx)
	eg.Go(func() error {
		return run(ctx)
	})
	in := "message"
	rsp, err := http.Get("http://localhost:18080/" + in)
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