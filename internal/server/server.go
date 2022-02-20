package server

import (
	"context"
	"net/http"
	"time"

	"github.com/jackc/pgx/v4"
	"github.com/juicyluv/ReadyRead/config"
	"github.com/juicyluv/ReadyRead/internal/author"
	"github.com/juicyluv/ReadyRead/internal/openapi"
	"github.com/juicyluv/ReadyRead/internal/user"
	"github.com/juicyluv/ReadyRead/pkg/logger"
	"github.com/julienschmidt/httprouter"
)

// Server represents http server.
type Server struct {
	server  *http.Server
	logger  *logger.Logger
	cfg     *config.Config
	handler *httprouter.Router
}

// NewServer returns a new Server instance.
func NewServer(cfg *config.Config, handler *httprouter.Router, logger *logger.Logger) *Server {
	return &Server{
		server: &http.Server{
			Handler:        handler,
			WriteTimeout:   time.Duration(cfg.Http.WriteTimeout) * time.Second,
			ReadTimeout:    time.Duration(cfg.Http.ReadTimeout) * time.Second,
			MaxHeaderBytes: cfg.Http.MaxHeaderBytes << 20,
			Addr:           ":" + cfg.Http.Port,
		},
		logger:  logger,
		cfg:     cfg,
		handler: handler,
	}
}

// Run initializes storages, services, handlers and then starts http server. Returns an error on failure.
func (s *Server) Run(dbConn *pgx.Conn) error {
	reqTimeout := s.cfg.DB.RequestTimeout

	s.logger.Info("initializing routes")

	userStorage := user.NewStorage(dbConn, reqTimeout)
	userService := user.NewService(userStorage, *s.logger)
	userHandler := user.NewHandler(*s.logger, userService)
	userHandler.Register(s.handler)
	s.logger.Info("initialized user routes")

	authorStorage := author.NewStorage(dbConn, reqTimeout)
	authorService := author.NewService(authorStorage, *s.logger)
	authorHandler := author.NewHandler(*s.logger, authorService)
	authorHandler.Register(s.handler)
	s.logger.Info("initialized author routes")

	openapi.InitSwagger(s.handler)
	s.logger.Info("initialized documentation")

	return s.server.ListenAndServe()
}

// Shutdown closes all connections and shuts down http server.
// It uses httpServer.Shutdown() method. Returns an error on failure.
func (s *Server) Shutdown(ctx context.Context) error {
	return s.server.Shutdown(ctx)
}
