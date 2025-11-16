package infrastructure

import (
	"errors"
	"shortner/config"

	"github.com/wb-go/wbf/dbpg"
)

var (
	ErrURLNotFound = errors.New("url not found")
	ErrAliasExists = errors.New("alias already exists")
)

type PostgresStorage struct {
	db *dbpg.DB
}

func New(db *dbpg.DB) *PostgresStorage {
	return &PostgresStorage{db: db}
}

func NewPostgres(cfg config.PostgresConfig) (*dbpg.DB, error) {
	opts := &dbpg.Options{
		MaxOpenConns:    cfg.MaxOpenConns,
		MaxIdleConns:    cfg.MaxIdleConns,
		ConnMaxLifetime: cfg.ConnMaxLifetime,
	}

	return dbpg.New(cfg.DSN, cfg.SlavesDSN, opts)
}
