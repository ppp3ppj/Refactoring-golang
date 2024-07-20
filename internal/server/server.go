package server

import (
	"context"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"sync"
	"syscall"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"github.com/labstack/gommon/log"
	"github.com/ppp3ppj/go-refactoring-workshop/config"
	"github.com/ppp3ppj/go-refactoring-workshop/db"
	server_middlewares "github.com/ppp3ppj/go-refactoring-workshop/internal/middlewares"
)

type echoServer struct {
    app *echo.Echo
    conf *config.Config
    db db.IDatabase
}

var (
    server *echoServer
    once sync.Once
)

func NewEchoServer(conf *config.Config, db db.IDatabase) *echoServer {
    echoApp := echo.New()
    echoApp.Logger.SetLevel(log.DEBUG)

    once.Do(func() {
        server = &echoServer{
            app: echoApp,
            conf: config.ConfigGetting(),
            db: db,
        }
    })

    return server
}

func (s * echoServer) Start() {
    timeOutMiddleware := server_middlewares.GetTimeOutMiddleware(s.conf.Server.Timeout)
    corsMiddleware := server_middlewares.GetCORSMiddleware(s.conf.Server.AllowOrigins)

    s.app.Use(middleware.Recover())
    s.app.Use(middleware.Logger())

    s.app.Use(timeOutMiddleware)
    s.app.Use(corsMiddleware)

    // Graceful Shutdown
    quitCh := make(chan os.Signal, 1)
    signal.Notify(quitCh, syscall.SIGINT, syscall.SIGTERM)
    go s.gracefullyShutdown(quitCh)

    s.httpListening()
}

func (s *echoServer) gracefullyShutdown(quitCh <-chan os.Signal) {
    ctx := context.Background()

    <-quitCh
    s.app.Logger.Info("Shutting down the service...")

    if err := s.app.Shutdown(ctx); err != nil {
        s.app.Logger.Fatalf("Error: %s", err.Error())
    }
}

func (s *echoServer) httpListening() {
    url := fmt.Sprintf(":%d", s.conf.Server.Port)
    if err := s.app.Start(url); err != nil && err != http.ErrServerClosed {
        s.app.Logger.Fatalf("shutting down the server: %v", err)
    }
}

func (s *echoServer) healthCheck(c echo.Context) error {
    return c.String(http.StatusOK, "OK")
}
