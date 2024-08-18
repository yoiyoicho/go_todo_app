package main

import (
	"io"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestNewMux(t *testing.T) {
	// httptest.NewRecorderはレスポンスを受け取るための構造体
	w := httptest.NewRecorder()
	// httptest.NewRequestはテスト用のリクエストを生成する
	// 第1引数はHTTPメソッド、第2引数はリクエストURL、第3引数はリクエストボディ
	r := httptest.NewRequest(http.MethodGet, "/health", nil)
	// テスト対象のハンドラを生成
	// sutはsystem under test（テスト対象システム）の略
	sut := NewMux()
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
