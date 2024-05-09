package main

import (
	"github.com/timickb/narration-engine/examples/services/blogs/internal/handler"
	"github.com/timickb/narration-engine/examples/services/blogs/internal/usecase"
	"github.com/timickb/narration-engine/examples/services/blogs/repository"
	"github.com/timickb/narration-engine/pkg/worker"
	schema "github.com/timickb/narration-engine/schema/v1/gen"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
	"log"
	"net"
)

func main() {
	publicationRepo := repository.NewPublicationRepo()
	blogUsecase := usecase.NewBlogUsecase(publicationRepo)

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
