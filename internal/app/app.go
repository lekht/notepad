package app

import (
	"os"
	"os/signal"
	"syscall"

	"github.com/lekht/notepad/config"
	"github.com/lekht/notepad/internal/auth"
	"github.com/lekht/notepad/internal/controllers"
	"github.com/lekht/notepad/internal/repository"
	"github.com/lekht/notepad/internal/usecase"
	"github.com/lekht/notepad/pkg/httpserver"
	"github.com/lekht/notepad/pkg/logger"
	"github.com/lekht/notepad/pkg/postgres"
	"github.com/lekht/notepad/pkg/speller"
)

func Run(cfg *config.Config) {
	l := logger.New(&cfg.Logger)

	pg, err := postgres.New(&cfg.Postgres)
	if err != nil {
		l.Fatal(err, "failed to create new postgres connection")
	}

	r := repository.New(pg)

	a := auth.New()

	s := speller.New(&cfg.Speller)

	u := usecase.New(r, a, s)

	router := controllers.NewRouter(u, l)

	server := httpserver.New(router, httpserver.Port(cfg.Server.Port))

	l.Info("Server is running on port " + cfg.Server.Port)

	interrupt := make(chan os.Signal, 1)
	signal.Notify(interrupt, os.Interrupt, syscall.SIGTERM)

	select {
	case s := <-interrupt:
		l.Info(s.String(), "interrupt signal recieved")
	case err = <-server.Notify():
		l.Error(err, "server error recieved")
	}

	err = server.Shutdown()
	if err != nil {
		l.Error(err, "failed to shutdown server")
	}

	l.Info("all is still good")
}
