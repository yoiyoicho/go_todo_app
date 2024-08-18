package main

import (
	"context"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/yoiyoicho/go_todo_app/config"
)

func TestNewMux(t *testing.T) {
	// ポートのエラーが出るので、一旦スキップ
	// mux_test.go:28: failed to NewMux: dial tcp 127.0.0.1:3306: connect: connection refused
	t.Skip()

	// httptest.NewRecorderはレスポンスを受け取るための構造体
	w := httptest.NewRecorder()
	// httptest.NewRequestはテスト用のリクエストを生成する
	// 第1引数はHTTPメソッド、第2引数はリクエストURL、第3引数はリクエストボディ
	r := httptest.NewRequest(http.MethodGet, "/health", nil)
	// テスト対象のハンドラを生成
	// sutはsystem under test（テスト対象システム）の略
	cfg, err := config.New()
	if err != nil {
		t.Fatalf("failed to initiate config: %v", err)
	}
	sut, cleanup, err := NewMux(context.Background(), cfg)
	defer cleanup()
	if err != nil {
		t.Fatalf("failed to NewMux: %v", err)
	}
	// ServeHTTPメソッドを呼び出すことで、テスト対象のハンドラを実行する
	sut.ServeHTTP(w, r)
	resp := w.Result()
	t.Cleanup(func() { _ = resp.Body.Close() })

	if resp.StatusCode != http.StatusOK {
		t.Error("want status code 200, but", resp.StatusCode)
	}
	got, err := io.ReadAll(resp.Body)
	if err != nil {
		t.Fatalf("failed to read body: %v", err)
	}
	want := `{"status": "ok"}`
	if string(got) != want {
		t.Errorf("want %q, but got %q", want, got)
	}
}
