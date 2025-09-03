package main

import (
	"context"
	"database/sql"
	"errors"
	"net/http"
	"os"
	"os/signal"
	"sync"
	"syscall"
	"time"

	_ "github.com/jackc/pgx/v5/stdlib"

	"github.com/Thiht/transactor"
	"github.com/nasik90/gophkeeper/cmd/gophkeeper-server/settings"
	"github.com/nasik90/gophkeeper/internal/common/logger"
	"github.com/nasik90/gophkeeper/internal/server"
	"github.com/nasik90/gophkeeper/internal/server/handler"
	"github.com/nasik90/gophkeeper/internal/server/service"
	"github.com/nasik90/gophkeeper/internal/server/storage/pg"
	"go.uber.org/zap"
)

func main() {
	options := parseOptions()
	if err := logger.Initialize(options.LogLevel); err != nil {
		panic(err)
	}
	store, transactor := initStore(options)
	service := service.NewService(store, transactor)
	handler := handler.NewHandler(service)
	server := server.NewServer(handler, options.ServerAddress)

	runServers(server, store)
}

func parseOptions() *settings.Options {
	options := new(settings.Options)
	settings.ParseFlags(options)
	return options
}

func initStore(options *settings.Options) (service.Store, transactor.Transactor) {
	var (
		store      service.Store
		transactor transactor.Transactor
	)

	if options.DatabaseDSN != "" {
		conn, err := sql.Open("pgx", options.DatabaseDSN)
		if err != nil {
			logger.Log.Fatal("open pgx conn", zap.String("DatabaseDSN", options.DatabaseDSN), zap.Error(err))
		}
		store, transactor, err = pg.NewStore(conn)
		if err != nil {
			logger.Log.Fatal("create pg repo", zap.String("DatabaseDSN", options.DatabaseDSN), zap.Error(err))
		}
	} else {
		logger.Log.Fatal("can not initialize store: db settings are empty")
	}

	return store, transactor
}

func runServers(server *server.Server, store service.Store) {
	var wg sync.WaitGroup

	wg.Add(1)
	go func() {
		defer wg.Done()
		if err := server.RunServer(); err != nil && !errors.Is(err, http.ErrServerClosed) {
			logger.Log.Fatal("run server", zap.Error(err))
		}
	}()

	// wg.Add(1)
	// go func() {
	// 	defer wg.Done()
	// 	if err := grpcServer.RunServer(); err != nil {
	// 		logger.Log.Fatal("run grpc server", zap.Error(err))
	// 	}
	// }()

	sigs := make(chan os.Signal, 1)
	signal.Notify(sigs, syscall.SIGINT, syscall.SIGTERM, syscall.SIGQUIT)

	wg.Add(1)
	go func() {
		defer wg.Done()
		<-sigs
		ctx, cancel := context.WithTimeout(context.Background(), 20*time.Second)
		defer cancel()

		logger.Log.Info("closing http server")
		if err := server.StopServer(ctx); err != nil {
			logger.Log.Error("stop http server", zap.Error(err))
		}

		// logger.Log.Info("closing grpc server")
		// grpcServer.StopServer()

		logger.Log.Info("closing the storage")
		if err := store.Close(); err != nil {
			logger.Log.Error("close storage", zap.Error(err))
		}

		logger.Log.Info("ready to exit")
	}()

	wg.Wait()
	logger.Log.Info("closed gracefully")
}
