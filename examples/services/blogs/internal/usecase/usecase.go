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

func (u *BlogUsecase) BlogGetById(ctx context.Context, id uuid.UUID) (*domain.Blog, error) {
	blog, err := u.blogRepo.GetById(ctx, id)
	if err != nil {
		return nil, fmt.Errorf("blogRepo.GetById: %w", err)
	}
	return blog, nil
}

func (u *BlogUsecase) BlogCreate(ctx context.Context, blog *domain.Blog) error {
	if err := u.blogRepo.Create(ctx, blog); err != nil {
		return fmt.Errorf("blogRepo.Create: %w", err)
	}
	return nil
}

func (u *BlogUsecase) BlogUpdateStats(ctx context.Context, dto *domain.BlogUpdateStatsDto) error {
	blog, err := u.blogRepo.GetById(ctx, dto.Id)
	if err != nil {
		return fmt.Errorf("blogRepo.GetById: %w", err)
	}

	updateDto := &domain.BlogUpdateDto{Id: dto.Id}
	if dto.IncSubscribers {
		updateDto.SubscribersCount = utils.Ptr(blog.SubscribersCount + 1)
	}
	if dto.IntPublications {
		updateDto.PublicationsCount = utils.Ptr(blog.PublicationsCount + 1)
	}
	if dto.IncDonations {
		updateDto.DonationsCount = utils.Ptr(blog.DonationsCount + 1)
	}

	if err = u.blogRepo.Update(ctx, updateDto); err != nil {
		return fmt.Errorf("blogRepo.Update")
	}
	return nil
}
