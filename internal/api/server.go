package api

import (
	"context"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	myhttp "github.com/aaltgod/telegram-bot/internal/api/delivery/http"
	"github.com/labstack/echo/v4"
	"github.com/labstack/gommon/log"
	"github.com/sirupsen/logrus"
)

type Server struct {
	logger  *logrus.Logger
	handler *myhttp.Handler
}

func NewServer(logger *logrus.Logger, handler *myhttp.Handler) *Server {
	return &Server{
		logger:  logger,
		handler: handler,
	}
}

func (s *Server) Start() error {

	e := echo.New()

	e.Logger.SetLevel(log.INFO)

	e.GET("/api/get_user", s.handler.GetUser)
	e.GET("/api/get_history_by_tg", s.handler.GetUserHistory)
	e.DELETE("/api/delete_record", s.handler.DeleteOneRequest)
	e.GET("/api/get_users", s.handler.GetUsers)

	go func() {
		if err := e.Start(":" + os.Getenv("HTTP_API_PORT")); err != nil && err != http.ErrServerClosed {
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
