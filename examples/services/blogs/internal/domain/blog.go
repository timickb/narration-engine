package domain

import (
	"github.com/google/uuid"
	"time"
)

// Blog Сущность блога.
type Blog struct {
	Id                uuid.UUID
	AuthorId          uuid.UUID
	AuthorEmail       string
	Name              string
	SubscribersCount  int
	PublicationsCount int
	DonationsCount    int
	CreatedAt         time.Time
	UpdatedAt         time.Time
}

// BlogUpdateDto Структура для обновления информации о блоге.
type BlogUpdateDto struct {
	Id                uuid.UUID
	Name              *string
	SubscribersCount  *int
	PublicationsCount *int
	DonationsCount    *int
}

// BlogUpdateStatsDto Структура для обновления статистики блога.
type BlogUpdateStatsDto struct {
	Id              uuid.UUID
	IncSubscribers  bool
	IntPublications bool
	IncDonations    bool
}
