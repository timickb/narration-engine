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

// BlogUsecase Основной юзкейс сервиса.
type BlogUsecase interface {
	PublicationUpdate(ctx context.Context, dto *PublicationUpdateDto) error
	PublicationGetById(ctx context.Context, id uuid.UUID) (*Publication, error)
}
