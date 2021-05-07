package repository

import (
	"context"
	"fmt"

	"github.com/go-pg/pg/v10"
)

type DBConfig struct {
	Host     string
	Port     string
	Username string
	Password string
	DBName   string
	SSLMode  string
}

func NewPostgresDB(cfg DBConfig) (*pg.DB, error) {
	pgURL := fmt.Sprintf("postgres://%s:%s@%s:%s/%s?sslmode=%s", cfg.Username, cfg.Password, cfg.Host, cfg.Port, cfg.DBName, cfg.SSLMode)
	opt, err := pg.ParseURL(pgURL)
	if err != nil {
		return nil, err
	}

	db := pg.Connect(opt)

	ctx := context.Background()

	if err := db.Ping(ctx); err != nil {
		return nil, err
	}

	return db, nil
}
