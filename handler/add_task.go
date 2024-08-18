package handler

import (
	"encoding/json"
	"net/http"
	"time"

	"github.com/go-playground/validator/v10"
	"github.com/yoiyoicho/go_todo_app/entity"
	"github.com/yoiyoicho/go_todo_app/store"
)

type AddTask struct {
	Store     *store.TaskStore
	Validator *validator.Validate
}

// リクエストの処理が正常に完了する場合、RespondJSONを使ってレスポンスを返す
// なんらかのエラーがあった場合は、ErrResponse型に情報を含めてRespondJSONを使ってレスポンスを返す
func (at *AddTask) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	var b struct {
		Title string `json:"title" validate:"required"`
	}

	// リクエストボディからJSONデータを読み取り、b変数にデコードする
	if err := json.NewDecoder(r.Body).Decode(&b); err != nil {
		RespondJSON(ctx, w, ErrResponse{
			Message: err.Error(),
		}, http.StatusInternalServerError)
		return
	}

	// b変数のバリデーションを行う
	err := at.Validator.Struct(b)
	if err != nil {
		RespondJSON(ctx, w, &ErrResponse{
			Message: err.Error(),
		}, http.StatusBadRequest)
		return
	}

	t := &entity.Task{
		Title:   b.Title,
		Status:  entity.TaskStatusTodo,
		Created: time.Now(),
	}

	// インメモリストアにタスクを追加する
	id, err := store.Tasks.Add(t)
	if err != nil {
		RespondJSON(ctx, w, &ErrResponse{
			Message: err.Error(),
		}, http.StatusInternalServerError)
		return
	}
	resp := struct {
		ID int `json:"id"`
	}{ID: int(id)}
	RespondJSON(ctx, w, resp, http.StatusCreated)
}
