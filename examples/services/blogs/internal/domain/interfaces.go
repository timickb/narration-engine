package domain

import (
	"context"
	"github.com/google/uuid"
)

// PublicationRepository Хранилище публикаций.
type PublicationRepository interface {
	Create(ctx context.Context, publication *Publication) error
	GetById(ctx context.Context, id uuid.UUID) (*Publication, error)
	Update(ctx context.Context, dto *PublicationUpdateDto) error
}

// BlogRepository Хранилище блогов.
type BlogRepository interface {
	Create(ctx context.Context, blog *Blog) error
	GetById(ctx context.Context, id uuid.UUID) (*Blog, error)
	Update(ctx context.Context, dto *BlogUpdateDto) error
}

// BlogUsecase Основной юзкейс сервиса.
type BlogUsecase interface {
	BlogGetById(ctx context.Context, id uuid.UUID) (*Blog, error)
	BlogCreate(ctx context.Context, blog *Blog) error
	BlogUpdateStats(ctx context.Context, dto *BlogUpdateStatsDto) error
	PublicationCreate(ctx context.Context, publication *Publication) error
	PublicationUpdate(ctx context.Context, dto *PublicationUpdateDto) error
	PublicationGetById(ctx context.Context, id uuid.UUID) (*Publication, error)
}
