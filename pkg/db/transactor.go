package db

import (
	"context"
	"errors"
	"gorm.io/gorm"
)

type transactor struct {
	db *Database
}

func NewTransactor(db *Database) *transactor {
	return &transactor{db: db}
}

func (t *transactor) Transaction(ctx context.Context, fn func(ctx context.Context) error) error {
	db := fetchDbFromCtx(ctx)
	if db != nil {
		return errors.New("transaction is already started")
	}

	return t.db.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		return fn(putDbToCtx(ctx, &Database{db: tx}))
	})
}
