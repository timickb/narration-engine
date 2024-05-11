package domain

import (
	"github.com/google/uuid"
	"time"
)

// PublicationStatus Статус публикации.
type PublicationStatus string

const (
	// PublicationStatusReview Публикация на рассмотрении.
	PublicationStatusReview PublicationStatus = "REVIEW"
	// PublicationStatusPending Публикация на доработке.
	PublicationStatusPending PublicationStatus = "PENDING"
	// PublicationStatusApproved Публикация утверждена и опубликована.
	PublicationStatusApproved PublicationStatus = "APPROVED"
	// PublicationStatusDeclined Публикация отклонена.
	PublicationStatusDeclined PublicationStatus = "DECLINED"
)

// Publication Публикация в блоге.
type Publication struct {
	Id        uuid.UUID
	AuthorId  uuid.UUID
	BlogId    uuid.UUID
	Title     string
	Body      string
	Status    PublicationStatus
	CreatedAt time.Time
	UpdatedAt time.Time
}

// PublicationUpdateDto Структура для обновления публикации в блоге.
type PublicationUpdateDto struct {
	Id     uuid.UUID
	Status *PublicationStatus
	Title  *string
	Body   *string
}
