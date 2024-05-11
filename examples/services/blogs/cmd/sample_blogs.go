package main

import (
	"github.com/google/uuid"
	"github.com/timickb/blogs-example/internal/domain"
	"time"
)

var (
	SampleBlogTravel = &domain.Blog{
		Id:                uuid.MustParse("12434cd2-0eeb-11ef-b3a0-cf5d00058889"),
		AuthorId:          uuid.MustParse("1adc580c-0eeb-11ef-8e57-57f257f46cf1"),
		Name:              "Путешествия",
		SubscribersCount:  21,
		PublicationsCount: 6,
		CreatedAt:         time.Now().Add(-time.Hour),
		UpdatedAt:         time.Now().Add(-time.Hour),
	}

	SampleBlogFrontend = &domain.Blog{
		Id:                uuid.MustParse("4001d210-0eeb-11ef-a41f-071635578df5"),
		AuthorId:          uuid.MustParse("43c47d58-0eeb-11ef-9766-d7668da30418"),
		Name:              "Frontend",
		SubscribersCount:  8,
		PublicationsCount: 2,
		CreatedAt:         time.Now().Add(-time.Hour),
		UpdatedAt:         time.Now().Add(-time.Hour),
	}

	SampleBlogFinance = &domain.Blog{
		Id:                uuid.MustParse("66071ee8-0eeb-11ef-8aa6-c7962b59b80e"),
		AuthorId:          uuid.MustParse("6b824942-0eeb-11ef-a3d5-af8a165635dd"),
		Name:              "Финансы",
		SubscribersCount:  13,
		PublicationsCount: 1,
		CreatedAt:         time.Now().Add(-time.Hour),
		UpdatedAt:         time.Now().Add(-time.Hour),
	}
)
