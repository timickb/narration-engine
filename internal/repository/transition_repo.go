package repository

import (
	"github.com/timickb/narration-engine/pkg/db"
)

type transitionRepo struct {
	db *db.Database
}

func NewTransitionRepo(db *db.Database) *transitionRepo {
	return &transitionRepo{db: db}
}
