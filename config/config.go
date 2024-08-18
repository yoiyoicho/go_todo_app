package config

import (
	"github.com/caarlos0/env/v6"
)

type Config struct {
	// 構造体タグ
	Env        string `env:"ENV" envDefault:"dev"`
	Port       int    `env:"PORT" envDefault:"8080"`
	DBHost     string `env:"TODO_DB_HOST" envDefault:"127.0.0.1"`
	DBPort     int    `env:"TODO_DB_PORT" envDefault:"3306"`
	DBUser     string `env:"TODO_DB_USER" envDefault:"todo"`
	DBPassword string `env:"TODO_DB_PASSWORD" envDefault:"todo"`
	DBName     string `env:"TODO_DB_NAME" envDefault:"todo"`
}

func New() (*Config, error) {
	// Config構造体のポインタを作成
	cfg := &Config{}
	// 環境変数から設定値を読み込む
	if err := env.Parse(cfg); err != nil {
		return nil, err
	}
	return cfg, nil
}
