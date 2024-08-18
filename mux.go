package main

import (
	"context"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/go-playground/validator/v10"
	"github.com/yoiyoicho/go_todo_app/auth"
	"github.com/yoiyoicho/go_todo_app/clock"
	"github.com/yoiyoicho/go_todo_app/config"
	"github.com/yoiyoicho/go_todo_app/handler"
	"github.com/yoiyoicho/go_todo_app/service"
	"github.com/yoiyoicho/go_todo_app/store"
)

func NewMux(ctx context.Context, cfg *config.Config) (http.Handler, func(), error) {
	mux := chi.NewRouter()
	mux.HandleFunc("/health", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json; charset=utf-8")
		_, _ = w.Write([]byte(`{"status": "ok"}`))
	})
	v := validator.New()
	db, cleanup, err := store.New(ctx, cfg)
	if err != nil {
		return nil, cleanup, err
	}
	clocker := clock.RealClocker{}
	r := store.Repository{Clocker: clocker}
	at := &handler.AddTask{
		Service:   &service.AddTask{DB: db, Repo: &r},
		Validator: v,
	}
	// Postの第2引数はhttp.HandlerFunc型の関数を受け取る
	// なぜhttp.HandlerFunc(at.ServeHTTP)としなくてもOKなのか？
	mux.Post("/tasks", at.ServeHTTP)
	lt := &handler.ListTask{
		Service: &service.ListTask{DB: db, Repo: &r},
	}
	mux.Get("/tasks", lt.ServeHTTP)
	ru := &handler.RegisterUser{
		Service:   &service.RegisterUser{DB: db, Repo: &r},
		Validator: v,
	}
	mux.Post("/register", ru.ServeHTTP)
	rcli, err := store.NewKVS(ctx, cfg)
	if err != nil {
		return nil, cleanup, err
	}
	jwter, err := auth.NewJWTer(rcli, clocker)
	if err != nil {
		return nil, cleanup, err
	}
	l := &handler.Login{
		Service: &service.Login{
			DB: db,
			Repo: &r,
			TokenGenerator: jwter,
		},
		Validator: v,
	}
	mux.Post("/login", l.ServeHTTP)
	return mux, cleanup, nil
}
