package usecase

import (
	"context"
	"fmt"
	"github.com/google/uuid"
	"github.com/timickb/blogs-example/internal/domain"
	"github.com/timickb/narration-engine/pkg/utils"
)

type BlogUsecase struct {
	blogRepo        domain.BlogRepository
	publicationRepo domain.PublicationRepository
}

func NewBlogUsecase(
	blogRepo domain.BlogRepository, publicationRepo domain.PublicationRepository,
) *BlogUsecase {
	return &BlogUsecase{blogRepo: blogRepo, publicationRepo: publicationRepo}
}

func (u *BlogUsecase) PublicationCreate(ctx context.Context, publication *domain.Publication) error {
	publication.Status = domain.PublicationStatusReview
	if err := u.publicationRepo.Create(ctx, publication); err != nil {
		return fmt.Errorf("publicationRepo.Create: %w", err)
	}
	return nil
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

func (u *BlogUsecase) BlogCreate(ctx context.Context, blog *domain.Blog) error {
	if err := u.blogRepo.Create(ctx, blog); err != nil {
		return fmt.Errorf("blogRepo.Create: %w", err)
	}
	return nil
}

func (u *BlogUsecase) BlogUpdateStats(ctx context.Context, blogId uuid.UUID, incPublications, incSubscribers bool) error {
	blog, err := u.blogRepo.GetById(ctx, blogId)
	if err != nil {
		return fmt.Errorf("blogRepo.GetById: %w", err)
	}

	dto := &domain.BlogUpdateDto{Id: blogId}
	if incSubscribers {
		dto.SubscribersCount = utils.Ptr(blog.SubscribersCount + 1)
	}
	if incPublications {
		dto.PublicationsCount = utils.Ptr(blog.PublicationsCount + 1)
	}

	if err = u.blogRepo.Update(ctx, dto); err != nil {
		return fmt.Errorf("blogRepo.Update")
	}
	return nil
}
