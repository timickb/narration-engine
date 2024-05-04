package main

import (
	"github.com/timickb/narration-engine/pkg/worker"
	schema "github.com/timickb/narration-engine/schema/v1/gen"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
	"log"
	"net"
	"notifications/internal/handler"
)

func main() {
	sendMessageHandler := handler.NewSendMessageHandler()
	h := handler.New(map[string]worker.Worker{
		sendMessageHandler.Name(): sendMessageHandler,
	})

	listener, _ := net.Listen("tcp", ":5001")

	srv := grpc.NewServer()
	schema.RegisterWorkerServiceServer(srv, h)
	reflection.Register(srv)

	log.Println("Listen notifications service at port 5001")
	srv.Serve(listener)
}
