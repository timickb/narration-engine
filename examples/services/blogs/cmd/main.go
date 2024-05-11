package main

import (
	"context"
	"github.com/timickb/blogs-example/internal/handler"
	"github.com/timickb/blogs-example/internal/usecase"
	"github.com/timickb/blogs-example/repository"
	"github.com/timickb/narration-engine/pkg/worker"
	schema "github.com/timickb/narration-engine/schema/v1/gen"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
	"log"
	"net"
)

func main() {
	blogRepo := repository.NewBlogRepo()
	publicationRepo := repository.NewPublicationRepo()
	blogUsecase := usecase.NewBlogUsecase(blogRepo, publicationRepo)

	blogUsecase.BlogCreate(context.Background(), SampleBlogFinance)
	blogUsecase.BlogCreate(context.Background(), SampleBlogFrontend)
	blogUsecase.BlogCreate(context.Background(), SampleBlogTravel)

	publicationUpdateHandler := handler.NewPublicationUpdateHandler(blogUsecase)
	stubHandler := handler.NewStubHandler()

	h := handler.New(map[string]worker.Worker{
		publicationUpdateHandler.Name(): publicationUpdateHandler,
		stubHandler.Name():              stubHandler,
	})

	listener, _ := net.Listen("tcp", ":5003")

	srv := grpc.NewServer()
	schema.RegisterWorkerServiceServer(srv, h)
	reflection.Register(srv)

	log.Println("Listen blogs service at port 5003")
	srv.Serve(listener)
}
