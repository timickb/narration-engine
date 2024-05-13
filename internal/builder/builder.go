package builder

import (
	"context"
	"fmt"
	"github.com/google/uuid"
	log "github.com/sirupsen/logrus"
	"github.com/timickb/narration-engine/internal/adapter/handler"
	"github.com/timickb/narration-engine/internal/config"
	"github.com/timickb/narration-engine/internal/controller"
	"github.com/timickb/narration-engine/internal/core"
	"github.com/timickb/narration-engine/internal/domain"
	"github.com/timickb/narration-engine/internal/repository"
	"github.com/timickb/narration-engine/internal/usecase"
	"github.com/timickb/narration-engine/migrations"
	"github.com/timickb/narration-engine/pkg/db"
	schema "github.com/timickb/narration-engine/schema/v1/gen"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/grpc/reflection"
	"net"
	"sync"
)

// Builder DI контейнер приложения.
type Builder struct {
	ctx           context.Context
	cfg           *config.Config
	log           *log.Logger
	db            *db.Database
	listener      net.Listener
	srv           *grpc.Server
	usecase       domain.Usecase
	runner        *core.InstanceRunner
	handlers      map[string]*handler.Handler
	handlersConns []*grpc.ClientConn
	wg            *sync.WaitGroup
}

func New(ctx context.Context, cfg *config.Config) (*Builder, error) {
	b := &Builder{ctx: ctx, cfg: cfg}
	b.log = log.WithContext(ctx).Logger
	b.wg = &sync.WaitGroup{}
	b.wg.Add(cfg.AsyncWorker.Count + 1)

	if err := b.initDatabase(); err != nil {
		return nil, fmt.Errorf("init database: %w", err)
	}
	b.buildUsecase()
	if err := b.buildGrpcServer(); err != nil {
		return nil, fmt.Errorf("build grpc server: %w", err)
	}
	if err := b.buildExternalHandlers(); err != nil {
		return nil, fmt.Errorf("build workers clients: %w", err)
	}
	b.buildInstanceRunner()

	return b, nil
}

func (b *Builder) WaitGroup() *sync.WaitGroup {
	return b.wg
}

func (b *Builder) ServeGrpc() error {
	return b.srv.Serve(b.listener)
}

func (b *Builder) StartInstanceRunner(ctx context.Context) {
	b.runner.Start(ctx)
}

func (b *Builder) ServerPort() int {
	return b.cfg.GrpcPort
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

func (b *Builder) buildInstanceRunner() {
	instanceChan := make(chan uuid.UUID)
	instanceRepo := repository.NewInstanceRepo(b.db)
	transitionRepo := repository.NewTransitionRepo(b.db)
	transactor := db.NewTransactor(b.db)

	handlerAdapters := make(map[string]core.HandlerAdapter)
	for key, value := range b.handlers {
		handlerAdapters[key] = value
	}

	b.runner = core.NewInstanceRunner(
		b.cfg,
		b.cfg,
		transactor,
		instanceRepo,
		transitionRepo,
		handlerAdapters,
		instanceChan,
		b.wg,
	)
}

func (b *Builder) buildGrpcServer() error {
	listener, err := net.Listen("tcp", fmt.Sprintf(":%d", b.cfg.GrpcPort))
	if err != nil {
		return fmt.Errorf("net.Listen: %w", err)
	}
	b.listener = listener

	ctrl := controller.New(b.usecase)
	srv := grpc.NewServer()
	schema.RegisterStateflowServiceServer(srv, ctrl)
	reflection.Register(srv)
	b.srv = srv
	return nil
}

func (b *Builder) buildUsecase() {
	instanceRepo := repository.NewInstanceRepo(b.db)
	eventRepo := repository.NewPendingEventRepo(b.db)
	transactor := db.NewTransactor(b.db)
	b.usecase = usecase.New(instanceRepo, eventRepo, transactor, b.cfg)
}

func (b *Builder) buildExternalHandlers() error {
	b.handlers = make(map[string]*handler.Handler)

	for name, conf := range b.cfg.Handlers {
		addr := fmt.Sprintf("%s:%d", conf.Host, conf.Port)
		conn, err := grpc.DialContext(b.ctx, addr, grpc.WithTransportCredentials(insecure.NewCredentials()))
		if err != nil {
			return err
		}
		b.handlers[name] = handler.NewHandlerClient(
			schema.NewWorkerServiceClient(conn),
			name,
		)
	}

	return nil
}

func (b *Builder) GracefulStop() error {
	for _, conn := range b.handlersConns {
		if err := conn.Close(); err != nil {
			b.log.Errorf(err.Error())
		}
	}
	return nil
}
