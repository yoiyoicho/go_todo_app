package config

import (
	"github.com/caarlos0/env/v6"
)

type Config struct {
	// 構造体タグ
	Env  string `env:"ENV" envDefault:"dev"`
	Port int    `env:"PORT" envDefault:"8080"`
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
