package repository

import (
	"context"
	"fmt"
	"github.com/google/uuid"
	"github.com/timickb/blogs-example/internal/domain"
	"sync"
	"time"
)

type blogRepo struct {
	sync.RWMutex
	data map[uuid.UUID]*domain.Blog
}

func NewBlogRepo() *blogRepo {
	return &blogRepo{data: make(map[uuid.UUID]*domain.Blog)}
}

// Create Добавить новый блог.
func (r *blogRepo) Create(ctx context.Context, blog *domain.Blog) error {
	r.Lock()
	r.Unlock()

	if _, ok := r.data[blog.Id]; ok {
		return fmt.Errorf("blog %s already exists", blog.Id.String())
	}
	r.data[blog.Id] = blog
	return nil
}

// GetById Получить блог по идентификатору.
func (r *blogRepo) GetById(ctx context.Context, id uuid.UUID) (*domain.Blog, error) {
	r.RLock()
	defer r.RUnlock()

	blog, ok := r.data[id]
	if !ok {
		return nil, fmt.Errorf("blog %s not found", id.String())
	}
	return blog, nil
}

// Update Обновиить информацию о блоге.
func (r *blogRepo) Update(ctx context.Context, dto *domain.BlogUpdateDto) error {
	r.Lock()
	defer r.Unlock()

	if _, ok := r.data[dto.Id]; !ok {
		return fmt.Errorf("blog %s not found", dto.Id.String())
	}

	if dto.Name != nil {
		r.data[dto.Id].Name = *dto.Name
	}
	if dto.PublicationsCount != nil {
		r.data[dto.Id].PublicationsCount = *dto.PublicationsCount
	}
	if dto.SubscribersCount != nil {
		r.data[dto.Id].SubscribersCount = *dto.SubscribersCount
	}

	r.data[dto.Id].UpdatedAt = time.Now()
	return nil
}
