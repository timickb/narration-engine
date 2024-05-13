package main

import (
	"github.com/timickb/narration-engine/pkg/worker"
	schema "github.com/timickb/narration-engine/schema/v1/gen"
	"github.com/timickb/notifications-example/internal/adapter/mail"
	"github.com/timickb/notifications-example/internal/handler"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
	"log"
	"net"
)

func main() {
	mailAdapter := mail.NewAdapter()

	sendMessageHandler := handler.NewSendMessageHandler(mailAdapter)
	sendBlogCreatedMailHandler := handler.NewSendBlogCreateMailHandler(mailAdapter)
	sendDonationReceivedMailHandler := handler.NewSendDonationReceivedMail(mailAdapter)
	sendDonationSentMailHandler := handler.NewSendDonationSentMail(mailAdapter)

	h := handler.New(map[string]worker.Worker{
		sendMessageHandler.Name():              sendMessageHandler,
		sendBlogCreatedMailHandler.Name():      sendBlogCreatedMailHandler,
		sendDonationReceivedMailHandler.Name(): sendDonationReceivedMailHandler,
		sendDonationSentMailHandler.Name():     sendDonationSentMailHandler,
	})

	listener, _ := net.Listen("tcp", ":5001")

	srv := grpc.NewServer()
	schema.RegisterWorkerServiceServer(srv, h)
	reflection.Register(srv)

	log.Println("Listen notifications service at port 5001")
	srv.Serve(listener)
}
