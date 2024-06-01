package main

import (
	"context"
	"github.com/google/uuid"
	"github.com/shopspring/decimal"
	log "github.com/sirupsen/logrus"
	"github.com/timickb/narration-engine/pkg/db"
	"github.com/timickb/narration-engine/pkg/worker"
	schema "github.com/timickb/narration-engine/schema/v1/gen"
	"github.com/timickb/payments-example/internal/domain"
	"github.com/timickb/payments-example/internal/handler"
	"github.com/timickb/payments-example/internal/usecase"
	"github.com/timickb/payments-example/migrations"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
	"net"
	"time"
)

var (
	sampleAccount1 = &domain.Account{
		Id: uuid.MustParse("36b0dfde-1086-11ef-af3d-cb9dadff5e5c"),
		// Автор страницы Путешествия
		UserId:    uuid.MustParse("1adc580c-0eeb-11ef-8e57-57f257f46cf1"),
		Amount:    decimal.RequireFromString("12340.0"),
		CreatedAt: time.Now().Add(-time.Hour),
	}
	sampleAccount2 = &domain.Account{
		Id: uuid.MustParse("4e3c5c64-1086-11ef-84a8-bf6b6103af1e"),
		// Автор страницы Frontend
		UserId:    uuid.MustParse("43c47d58-0eeb-11ef-9766-d7668da30418"),
		Amount:    decimal.RequireFromString("811.6"),
		CreatedAt: time.Now().Add(-time.Hour),
	}
	sampleAccount3 = &domain.Account{
		Id: uuid.MustParse("6195d0e2-1086-11ef-91dd-ef5bbe508c8b"),
		// Автор страницы Финансы
		UserId:    uuid.MustParse("6b824942-0eeb-11ef-a3d5-af8a165635dd"),
		Amount:    decimal.RequireFromString("0.0"),
		CreatedAt: time.Now().Add(-time.Hour),
	}
)

func main() {
	ctx := context.Background()
	d, err := db.CreatePostgresConnection(ctx, &db.PostgresConfig{
		Host:               "localhost",
		Name:               "payments",
		User:               "payments",
		Password:           "qwerty",
		SSLMode:            "disable",
		Port:               5452,
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
	err = migrations.Migrator.Migrate(sqlDb, "payments")
	if err != nil {
		log.Fatalf("make migration: %s", err.Error())
	}

	uc := usecase.New()
	uc.AccountCreate(context.Background(), sampleAccount1)
	uc.AccountCreate(context.Background(), sampleAccount2)
	uc.AccountCreate(context.Background(), sampleAccount3)

	invoiceCreateHandler := handler.NewInvoiceCreateHandler(uc)
	accountAddFundsHandler := handler.NewAccountAddFundsHandler(uc)
	accountRemoveFundsHandler := handler.NewAccountRemoveFundsHandler(uc)
	accountCreateHandler := handler.NewAccountCreateHandler(uc)

	handlerApi := handler.New(map[string]worker.Worker{
		invoiceCreateHandler.Name():      invoiceCreateHandler,
		accountAddFundsHandler.Name():    accountAddFundsHandler,
		accountRemoveFundsHandler.Name(): accountRemoveFundsHandler,
		accountCreateHandler.Name():      accountCreateHandler,
	})

	listener, _ := net.Listen("tcp", ":5002")

	srv := grpc.NewServer()
	schema.RegisterWorkerServiceServer(srv, handlerApi)
	reflection.Register(srv)

	log.Println("Listen payments service at port 5002")
	srv.Serve(listener)
}
