package main

import (
	"context"
	"github.com/google/uuid"
	"github.com/timickb/blogs-example/internal/domain"
	"github.com/timickb/blogs-example/internal/handler"
	"github.com/timickb/blogs-example/internal/usecase"
	"github.com/timickb/blogs-example/migrations"
	"github.com/timickb/blogs-example/repository"
	"github.com/timickb/narration-engine/pkg/db"
	"github.com/timickb/narration-engine/pkg/worker"
	schema "github.com/timickb/narration-engine/schema/v1/gen"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
	"log"
	"net"
	"time"
)

var (
	sampleBlogTravel = &domain.Blog{
		Id:                uuid.MustParse("12434cd2-0eeb-11ef-b3a0-cf5d00058889"),
		AuthorId:          uuid.MustParse("1adc580c-0eeb-11ef-8e57-57f257f46cf1"),
		AuthorEmail:       "somebody1@somewhere.com",
		Name:              "Путешествия",
		SubscribersCount:  21,
		PublicationsCount: 6,
		DonationsCount:    2,
		CreatedAt:         time.Now().Add(-time.Hour),
		UpdatedAt:         time.Now().Add(-time.Hour),
	}

	sampleBlogFrontend = &domain.Blog{
		Id:                uuid.MustParse("4001d210-0eeb-11ef-a41f-071635578df5"),
		AuthorId:          uuid.MustParse("43c47d58-0eeb-11ef-9766-d7668da30418"),
		AuthorEmail:       "somebody2@somewhere.com",
		Name:              "Frontend",
		SubscribersCount:  8,
		PublicationsCount: 2,
		DonationsCount:    3,
		CreatedAt:         time.Now().Add(-time.Hour),
		UpdatedAt:         time.Now().Add(-time.Hour),
	}

	sampleBlogFinance = &domain.Blog{
		Id:                uuid.MustParse("66071ee8-0eeb-11ef-8aa6-c7962b59b80e"),
		AuthorId:          uuid.MustParse("6b824942-0eeb-11ef-a3d5-af8a165635dd"),
		AuthorEmail:       "somebody3@somewhere.com",
		Name:              "Финансы",
		SubscribersCount:  13,
		PublicationsCount: 1,
		DonationsCount:    1,
		CreatedAt:         time.Now().Add(-time.Hour),
		UpdatedAt:         time.Now().Add(-time.Hour),
	}

	sampleFinancePublication = &domain.Publication{
		Id:        uuid.MustParse("f27b7b1a-12da-11ef-89a6-cb3a18c1abb8"),
		AuthorId:  sampleBlogFinance.AuthorId,
		BlogId:    sampleBlogFinance.Id,
		Title:     "Первая статья",
		Body:      "...",
		Status:    domain.PublicationStatusReview,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}
)

func main() {
	ctx := context.Background()
	d, err := db.CreatePostgresConnection(ctx, &db.PostgresConfig{
		Host:               "localhost",
		Name:               "blogs",
		User:               "blogs",
		Password:           "qwerty",
		SSLMode:            "disable",
		Port:               5453,
		MaxOpenConnections: 20,
		MaxIdleConnections: 20,
	})
	if err != nil {
		log.Fatal(err)
	}
	sqlDb, err := d.SqlDB()
	if err != nil {
		log.Fatalf("get sql db: %s", err.Error())
	}
	err = migrations.Migrator.Migrate(sqlDb, "blogs")
	if err != nil {
		log.Fatalf("make migration: %s", err.Error())
	}
	blogRepo := repository.NewBlogRepo()
	publicationRepo := repository.NewPublicationRepo()
	blogUsecase := usecase.NewBlogUsecase(blogRepo, publicationRepo)

	blogUsecase.BlogCreate(ctx, sampleBlogFinance)
	blogUsecase.BlogCreate(ctx, sampleBlogFrontend)
	blogUsecase.BlogCreate(ctx, sampleBlogTravel)
	blogUsecase.PublicationCreate(ctx, sampleFinancePublication)

	publicationUpdateHandler := handler.NewPublicationUpdateHandler(blogUsecase)
	statsUpdateHandler := handler.NewStatsUpdateHandler(blogUsecase)
	blogCreateHandler := handler.NewBlogCreateHandler(blogUsecase)
	stubHandler := handler.NewStubHandler()

	h := handler.New(map[string]worker.Worker{
		publicationUpdateHandler.Name(): publicationUpdateHandler,
		statsUpdateHandler.Name():       statsUpdateHandler,
		blogCreateHandler.Name():        blogCreateHandler,
		stubHandler.Name():              stubHandler,
	})

	listener, _ := net.Listen("tcp", ":5003")

	srv := grpc.NewServer()
	schema.RegisterWorkerServiceServer(srv, h)
	reflection.Register(srv)

	log.Println("Listen blogs service at port 5003")
	srv.Serve(listener)
}
