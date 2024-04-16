package builder

import (
	"context"
	"fmt"
	"github.com/timickb/go-stateflow/internal/config"
	"github.com/timickb/go-stateflow/internal/controller"
	"github.com/timickb/go-stateflow/internal/domain"
	"github.com/timickb/go-stateflow/internal/usecase"
	"github.com/timickb/go-stateflow/migrations"
	"github.com/timickb/go-stateflow/pkg/db"
	"google.golang.org/grpc"
	"net"
)

// Builder DI контейнер приложения.
type Builder struct {
	ctx      context.Context
	cfg      *config.Config
	db       *db.Database
	listener net.Listener
	srv      *grpc.Server
	usecase  domain.Usecase
}

func New(ctx context.Context, cfg *config.Config) (*Builder, error) {
	b := &Builder{ctx: ctx, cfg: cfg}

	if err := b.initDatabase(); err != nil {
		return nil, fmt.Errorf("init database: %w", err)
	}
	if err := b.buildUsecase(); err != nil {
		return nil, fmt.Errorf("build usecase: %w", err)
	}
	if err := b.buildGrpcServer(); err != nil {
		return nil, fmt.Errorf("build grpc server: %w", err)
	}

	return b, nil
}

func (b *Builder) ServeGrpc() error {
	return b.srv.Serve(b.listener)
}

func (b *Builder) initDatabase() error {
	d, err := db.CreatePostgresConnection(b.ctx, b.cfg.Database)
	if err != nil {
		return fmt.Errorf("create postgres connection: %w", err)
	}

	if b.cfg.Database.AutoMigrate {
		sqlDb, err := d.SqlDB()
		if err != nil {
			return fmt.Errorf("get sql db: %w", err)
		}
		err = migrations.Migrator.Migrate(sqlDb, b.cfg.Database.Name)
		if err != nil {
			return fmt.Errorf("make migration: %w", err)
		}
	}

	b.db = d
	return nil
}

func (b *Builder) buildGrpcServer() error {
	listener, err := net.Listen("tcp", fmt.Sprintf(":%d", b.cfg.GrpcPort))
	if err != nil {
		return fmt.Errorf("net.Listen: %w", err)
	}
	b.listener = listener

	_ = controller.New(b.usecase)
	srv := grpc.NewServer()
	// TODO: register controller
	b.srv = srv
	return nil
}

func (b *Builder) buildUsecase() error {
	// TODO: implement
	b.usecase = usecase.New(nil, nil)
	return nil
}
