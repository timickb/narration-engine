package repository

import (
	"context"
	"errors"
	"github.com/google/uuid"
	"github.com/timickb/blogs-example/internal/domain"
	"sync"
	"time"
)

type publicationRepo struct {
	sync.RWMutex
	data map[uuid.UUID]*domain.Publication
}

func NewPublicationRepo() *publicationRepo {
	return &publicationRepo{data: make(map[uuid.UUID]*domain.Publication)}
}

// Create Создать публикацию.
func (r *publicationRepo) Create(ctx context.Context, publication *domain.Publication) error {
	r.Lock()
	defer r.Unlock()

	if _, ok := r.data[publication.Id]; ok {
		return errors.New("publication already exists")
	}
	r.data[publication.Id] = publication
	return nil
}

// GetById Получить публикацию по идентификатору.
func (r *publicationRepo) GetById(ctx context.Context, id uuid.UUID) (*domain.Publication, error) {
	r.RLock()
	defer r.RUnlock()

	publication, ok := r.data[id]
	if !ok {
		return nil, errors.New("publication not found")
	}
	return publication, nil
}

// Update Обновить публикацию.
func (r *publicationRepo) Update(ctx context.Context, dto *domain.PublicationUpdateDto) error {
	r.Lock()
	defer r.Unlock()

	if _, ok := r.data[dto.Id]; !ok {
		return errors.New("publication not found")
	}

	if dto.Status != nil {
		r.data[dto.Id].Status = *dto.Status
	}
	if dto.Title != nil {
		r.data[dto.Id].Title = *dto.Title
	}
	if dto.Body != nil {
		r.data[dto.Id].Body = *dto.Body
	}

	r.data[dto.Id].UpdatedAt = time.Now()
	return nil
}
