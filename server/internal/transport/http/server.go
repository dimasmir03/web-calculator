package http

import (
	"context"
	"net/http"
	"time"

	"github.com/dimasmir03/web-calculator-server/internal/storage/sqlite"
	"github.com/labstack/echo/v4"
	"github.com/sirupsen/logrus"
)

type Server struct {
	Handler    http.Handler
	HttpServer *echo.Echo
	Address    string
	Logger     *logrus.Logger
	Storage    *sqlite.Storage
}

// NewServer create Server object
func NewServer(address string, e *echo.Echo, storage *sqlite.Storage, logger *logrus.Logger) *Server {
	logger.SetFormatter(&logrus.JSONFormatter{})
	return &Server{
		Address:    address,
		HttpServer: e,
		Logger:     logger,
		Storage:    storage,
	}
}

// Run http server
func (s *Server) Run() {
	s.Logger.Infof("Server starting at %s", s.Address)
	s.HttpServer.Start(s.Address)
}

// Stop http server
func (s *Server) Stop(ctx context.Context) {
	ctx, cancel := context.WithTimeout(ctx, 5*time.Second)
	if err := s.HttpServer.Shutdown(ctx); err != nil {
		s.Logger.Fatalf("Error with stop server: %s", err.Error())
	}
	cancel()
	s.Logger.Infoln("Server stopped")
}
