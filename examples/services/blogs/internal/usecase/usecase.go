package usecase

import (
	"context"
	"fmt"
	"github.com/google/uuid"
	"github.com/timickb/narration-engine/examples/services/blogs/internal/domain"
)

type BlogUsecase struct {
	publicationRepo domain.PublicationRepository
}

func NewBlogUsecase(publicationRepo domain.PublicationRepository) *BlogUsecase {
	return &BlogUsecase{publicationRepo: publicationRepo}
}

func (u *BlogUsecase) PublicationUpdate(ctx context.Context, dto *domain.PublicationUpdateDto) error {
	if err := u.publicationRepo.Update(ctx, dto); err != nil {
		return fmt.Errorf("publicationRepo.Update: %w", err)
	}
	return nil
}

func (u *BlogUsecase) PublicationGetById(ctx context.Context, id uuid.UUID) (*domain.Publication, error) {
	publication, err := u.publicationRepo.GetById(ctx, id)
	if err != nil {
		return nil, fmt.Errorf("publicationRepo.GetById: %w", err)
	}
	return publication, nil
}
