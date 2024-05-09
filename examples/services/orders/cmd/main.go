package main

import (
	"github.com/timickb/narration-engine/examples/services/orders/internal/handler"
	"github.com/timickb/narration-engine/pkg/worker"
	schema "github.com/timickb/narration-engine/schema/v1/gen"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
	"log"
	"net"
)

func main() {
	orderCreateHandler := handler.NewOrderCreateHandler()
	h := handler.New(map[string]worker.Worker{
		orderCreateHandler.Name(): orderCreateHandler,
	})

	listener, _ := net.Listen("tcp", ":5002")

	srv := grpc.NewServer()
	schema.RegisterWorkerServiceServer(srv, h)
	reflection.Register(srv)

	log.Println("Listen orders service at port 5002")
	srv.Serve(listener)
}
