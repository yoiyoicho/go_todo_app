package store

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"time"

	_ "github.com/go-sql-driver/mysql"
	"github.com/jmoiron/sqlx"
	"github.com/yoiyoicho/go_todo_app/clock"
	"github.com/yoiyoicho/go_todo_app/config"
)

// sqlx.DBオブジェクトのポインタ、クリーンアップ関数、エラーを返す
// *sqlx.DB型の値はRDBMSの利用終了後に*sql.DB.Closeメソッドを呼び出してコネクションを正しく終了する必要がある
// New関数の呼び出し元で終了処理をできるよう*sql.DB.Closeメソッドを実行する無名関数を返す
func New(ctx context.Context, cfg *config.Config) (*sqlx.DB, func(), error) {
	// MySQLデータベースへの接続を開始する
	// parseTime=true を忘れるとtime.Time型のフィールドに正しい時刻が入らない
	db, err := sql.Open("mysql",
		fmt.Sprintf(
			"%s:%s@tcp(%s:%d)/%s?parseTime=true",
			cfg.DBUser, cfg.DBPassword,
			cfg.DBHost, cfg.DBPort,
			cfg.DBName,
		),
	)
	if err != nil {
		return nil, nil, err
	}
	// 接続テスト
	// 2秒のタイムアウトでデータベースに対してPingを実行する
	ctx, cancel := context.WithTimeout(ctx, 2*time.Second)
	defer cancel()
	if err := db.PingContext(ctx); err != nil {
		return nil, func() { _ = db.Close() }, err
	}
	// 標準のsql.DBオブジェクトをsqlx.DBオブジェクトにラップする
	// データベース接続プールを表すオブジェクトで、複数のクエリやトランザクションを実行できる
	xdb := sqlx.NewDb(db, "mysql")
	return xdb, func() { _ = db.Close() }, nil
}

type Beginner interface {
	BeginTx(ctx context.Context, opts *sql.TxOptions) (*sql.Tx, error)
}

type Preparer interface {
	PreparexContext(ctx context.Context, query string) (*sqlx.Stmt, error)
}

// 書き込み系の操作を集めたインターフェース
type Execer interface {
	ExecContext(ctx context.Context, query string, args ...any) (sql.Result, error)
	NamedExecContext(ctx context.Context, query string, arg interface{}) (sql.Result, error)
}

// 参照系の主要なメソッドを集めたインターフェース
type Queryer interface {
	Preparer
	QueryxContext(ctx context.Context, query string, args ...any) (*sqlx.Rows, error)
	QueryRowxContext(ctx context.Context, query string, args ...any) *sqlx.Row
	GetContext(ctx context.Context, dest interface{}, query string, args ...any) error
	SelectContext(ctx context.Context, dest interface{}, query string, args ...any) error
}

// インターフェースが期待通りに宣言されている確認するコード
var (
	// sqlx.DB 型へのnilポインタをBeginnerインターフェース型の変数に代入する
	// 型チェックのためのコードで、変数自身は不要なので、ブランク識別子に代入する
	_ Beginner = (*sqlx.DB)(nil)
	_ Preparer = (*sqlx.DB)(nil)
	_ Queryer  = (*sqlx.DB)(nil)
	_ Queryer  = (*sqlx.Tx)(nil)
	_ Execer   = (*sqlx.DB)(nil)
	_ Execer   = (*sqlx.Tx)(nil)
)

type Repository struct {
	Clocker clock.Clocker
}

const (
	ErrCodeMySQLDuplicateEntry = 1062
)

var (
	ErrAlreadyEntry = errors.New("duplicate entry")
)
