package main

import (
	"context"
	"github.com/timickb/narration-engine/pkg/db"
	"github.com/timickb/narration-engine/pkg/worker"
	schema "github.com/timickb/narration-engine/schema/v1/gen"
	"github.com/timickb/notifications-example/internal/adapter/mail"
	"github.com/timickb/notifications-example/internal/handler"
	"github.com/timickb/notifications-example/migrations"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
	"log"
	"net"
)

func main() {
	ctx := context.Background()
	d, err := db.CreatePostgresConnection(ctx, &db.PostgresConfig{
		Host:               "localhost",
		Name:               "notifications",
		User:               "notifications",
		Password:           "qwerty",
		SSLMode:            "disable",
		Port:               5451,
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
	err = migrations.Migrator.Migrate(sqlDb, "notifications")
	if err != nil {
		log.Fatalf("make migration: %s", err.Error())
	}
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
