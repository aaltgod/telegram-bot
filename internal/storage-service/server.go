package storageservice

import (
	"context"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	myhttp "github.com/aaltgod/telegram-bot/internal/storage-service/delivery/http"
	"github.com/aaltgod/telegram-bot/internal/storage-service/repository"
	"github.com/labstack/echo/v4"
	"github.com/labstack/gommon/log"
	"github.com/sirupsen/logrus"
)

type Server struct {
	logger  *logrus.Logger
	handler *myhttp.Handler
	storage repository.Repository
}

func NewServer(logger *logrus.Logger, handler *myhttp.Handler, storage repository.Repository) *Server {
	return &Server{
		logger:  logger,
		handler: handler,
		storage: storage,
	}
}

func (s *Server) Start() error {

	e := echo.New()

	e.Logger.SetLevel(log.INFO)

	e.GET("/users/:id", s.handler.GetUser)
	e.PUT("/users/:id", s.handler.UpdateUser)
	e.POST("/users/", s.handler.InsertUser)
	e.GET("/users", s.handler.GetUsers)

	e.POST("/requests/:user_id", s.handler.AppendRequest)
	e.DELETE("/requests", s.handler.DeleteRequest)
	e.GET("/requests/:user_id", s.handler.GetAllRequestsByID)

	go func() {
		if err := e.Start(":" + os.Getenv("HTTP_STORAGE_SERVICE_PORT")); err != nil && err != http.ErrServerClosed {
			s.logger.Warnln("The service is shutting down")
		}
	}()

	interrupt := make(chan os.Signal, 1)
	signal.Notify(interrupt, os.Interrupt, syscall.SIGTERM)

	killSignal := <-interrupt
	switch killSignal {
	case os.Interrupt:
		s.logger.Warnln("Got SIGINT")
	case syscall.SIGTERM:
		s.logger.Warnln("Got SIGTERM")
	}

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	return e.Shutdown(ctx)
}
