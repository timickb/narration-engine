package db

import (
	"context"
	"database/sql"
	"fmt"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

type dbKey struct{}

// Database Обертка над gorm.DB
type Database struct {
	db *gorm.DB
}

// CreatePostgresConnection Установить соединение с PostgreSQL и создать инстанс Database.
func CreatePostgresConnection(ctx context.Context, cfg *PostgresConfig) (*Database, error) {
	db, err := gorm.Open(postgres.Open(cfg.DSNString()))
	if err != nil {
		return nil, fmt.Errorf("open postgres connection: %w", err)
	}

	// TODO: custom logger
	// TODO: configure secondaries

	sqlDb, err := db.DB()
	if err != nil {
		return nil, fmt.Errorf("get sql db: %w", err)
	}

	sqlDb.SetMaxIdleConns(cfg.MaxIdleConnections)
	sqlDb.SetMaxOpenConns(cfg.MaxOpenConnections)

	return &Database{db: db}, nil
}

// WithTxSupport Инстанс БД в контексте запущенной транзакции.
func (d *Database) WithTxSupport(ctx context.Context) *Database {
	dbWithCtx := fetchDbFromCtx(ctx)
	if dbWithCtx == nil {
		return d
	}
	return dbWithCtx
}

// SqlDB Получить инстанс sql.DB
func (d *Database) SqlDB() (*sql.DB, error) {
	return d.db.DB()
}

func fetchDbFromCtx(ctx context.Context) *Database {
	value := ctx.Value(dbKey{})
	if value == nil {
		return nil
	}

	if unwrapped, ok := value.(*Database); ok {
		return unwrapped
	}
	return nil
}

func putDbToCtx(ctx context.Context, db *Database) context.Context {
	return context.WithValue(ctx, dbKey{}, db)
}
