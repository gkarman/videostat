package app

import (
	"context"
	"errors"
	"fmt"
	"log/slog"
	"sync"
	"time"

	"github.com/gkarman/demo/internal/config"
	"github.com/gkarman/demo/internal/infrastructure/dispatcher"
	"github.com/gkarman/demo/internal/infrastructure/events"
	"github.com/gkarman/demo/internal/infrastructure/mq"
	grpc2 "github.com/gkarman/demo/internal/infrastructure/transport/grpc"
	"github.com/gkarman/demo/internal/infrastructure/transport/http"
	"github.com/gkarman/demo/internal/platform"
	"github.com/jackc/pgx/v5/pgxpool"
)

type Api struct {
	log          *slog.Logger
	db           *pgxpool.Pool
	serverHttp   *http.Server
	grpcServer   *grpc2.Server
	rabbitPusher *mq.RabbitPublisher
}

func NewApi(ctx context.Context) (*Api, error) {
	cfg, err := config.LoadConfig()
	if err != nil {
		return nil, fmt.Errorf("load config: %w", err)
	}

	log := platform.NewLogger(cfg)

	log.Info("db connect...")
	postgresDB, err := platform.NewPostgres(ctx, cfg)
	if err != nil {
		return nil, fmt.Errorf("connect to postgresDB: %w", err)
	}
	log.Info("db connected")

	log.Info("rabbitPusher connect...")
	rabbitPublisher, err := platform.NewRabbitPublisher(cfg)
	if err != nil {
		return nil, fmt.Errorf("rabbitPusher init: %w", err)
	}
	log.Info("rabbitPusher connected")

	d := dispatcher.New()
	events.RegisterEventHandlers(d, log, rabbitPublisher)

	serverHttp := platform.NewHTTPServer(log, postgresDB, cfg, d)
	serverGrpc, err := platform.NewGRPCServer(log, postgresDB, cfg, d)
	if err != nil {
		rabbitPublisher.Close()
		return nil, fmt.Errorf("create gRPC server: %w", err)
	}

	return &Api{
		log:          log,
		db:           postgresDB,
		serverHttp:   serverHttp,
		grpcServer:   serverGrpc,
		rabbitPusher: rabbitPublisher,
	}, nil
}

func (a *Api) Run(ctx context.Context) error {
	defer a.db.Close()
	defer a.rabbitPusher.Close()
	a.serverHttp.Start()
	a.grpcServer.Start()

	<-ctx.Done()

	a.log.Info("shutting down application", "reason", ctx.Err())

	shutdownCtx, cancel := context.WithTimeout(context.Background(), 15*time.Second)
	defer cancel()

	return a.shutdownServers(shutdownCtx)
}

func (a *Api) shutdownServers(ctx context.Context) error {

	var (
		wg        sync.WaitGroup
		errCh     = make(chan error, 2)
		joinedErr error
	)

	wg.Add(2)
	go func() {
		defer wg.Done()
		if err := a.serverHttp.Stop(ctx); err != nil {
			a.log.Error("serverHttp shutdown failed", "error", err)
			errCh <- fmt.Errorf("http shutdown: %w", err)
		}
	}()
	go func() {
		defer wg.Done()
		if err := a.grpcServer.Stop(ctx); err != nil {
			a.log.Error("gRPC server shutdown failed", "error", err)
			errCh <- fmt.Errorf("gRPC shutdown: %w", err)
		}
	}()

	wg.Wait()
	close(errCh)

	for err := range errCh {
		joinedErr = errors.Join(joinedErr, err)
	}

	return joinedErr
}
