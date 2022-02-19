package main

import (
	"context"
	"errors"
	"flag"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/jackc/pgx/v4"
	"github.com/juicyluv/ReadyRead/config"
	"github.com/juicyluv/ReadyRead/internal/server"
	"github.com/juicyluv/ReadyRead/pkg/logger"
	"github.com/julienschmidt/httprouter"
)

var (
	configPath = flag.String("config-path", "config/config.yml", "path for application configuration file")
)

func main() {
	flag.Parse()

	logger.Init()

	logger := logger.GetLogger()
	logger.Info("logger initialized")

	cfg := config.Get(*configPath, ".env")
	logger.Info("loaded config file")

	router := httprouter.New()
	logger.Info("initialized httprouter")

	logger.Info("connecting to database")

	pgxConfig, err := pgx.ParseConfig(cfg.DB.DSN)
	if err != nil {
		logger.Fatalf("cannot parse database config from dsn: %v", err)
	}

	dbTimeout, dbCancel := context.WithTimeout(context.Background(), time.Duration(cfg.DB.ConnectionTimeout)*time.Second)
	defer dbCancel()
	dbConn, err := pgx.ConnectConfig(dbTimeout, pgxConfig)
	if err != nil {
		logger.Fatalf("cannot connect to database: %v", err)
	}

	if err := dbConn.Ping(dbTimeout); err != nil {
		logger.Fatalf("cannot ping database: %v", err)
	}

	logger.Info("connected to database")

	// userStorage := db.NewStorage(mongoClient, cfg.DB.Collection)
	// userService := user.NewService(userStorage, logger)

	// userHandler := user.NewHandler(logger, userService)
	// userHandler.Register(router)
	// logger.Info("initialized user routes")

	// logger.Info("initializing swagger documentation")
	// internal.InitSwagger(router)
	// logger.Info("initialized swagger documentation")

	logger.Info("starting the server")
	srv := server.NewServer(cfg, router, &logger)

	quit := make(chan os.Signal, 1)
	signals := []os.Signal{syscall.SIGABRT, syscall.SIGQUIT, syscall.SIGHUP, os.Interrupt, syscall.SIGTERM}
	signal.Notify(quit, signals...)

	go func() {
		if err := srv.Run(); err != nil && !errors.Is(err, http.ErrServerClosed) {
			logger.Fatalf("cannot run the server: %v", err)
		}
	}()
	logger.Infof("server has been started on port %s", cfg.Http.Port)

	<-quit
	logger.Warn("shutting down the server")

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer func() {
		// TODO: Close database connection.
		logger.Info("closed mongo database connection")
		cancel()
	}()

	if err := srv.Shutdown(ctx); err != nil {
		logger.Errorf("server shutdown failed: %v", err)
	}

	logger.Info("server has been shutted down")
}
