package app

import (
	"github.com/zenorachi/todo-service/internal/config"
	"github.com/zenorachi/todo-service/internal/database"
	"github.com/zenorachi/todo-service/internal/server"
	"github.com/zenorachi/todo-service/pkg/database/postgres"
	"github.com/zenorachi/todo-service/pkg/logger"
	"os"
	"os/signal"
	"syscall"
)

func Run(cfg *config.Config) {
	/* DO MIGRATIONS */
	err := database.DoMigrations(&cfg.DB)
	if err != nil {
		logger.Fatal("migrations", "migrations failed")
	}
	logger.Info("migrations", "migrations done")

	/* INIT POSTGRES-DB */
	db, err := postgres.NewDB(&cfg.DB)
	defer func() { _ = db.Close() }()
	if err != nil {
		logger.Fatal("database-connection", err)
	}
	logger.Info("database", "postgres connected")

	/* INIT TOKEN MANAGER */
	//tokenManager := auth.NewManager(cfg.Auth.Secret)

	/* INIT SERVICES & DEPS */

	/* INIT HTTP HANDLER */

	/* INIT HTTP SERVER */
	srv := server.New(cfg, nil)
	srv.Run()
	logger.Info("server", "http server started")

	/* GRACEFUL SHUTDOWN */
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGTERM, syscall.SIGINT)

	/* WAITING FOR SYSCALL */
	select {
	case s := <-quit:
		logger.Info("app - Run - signal: " + s.String())
	case err = <-srv.Notify():
		logger.Error("server", err.Error())
	}

	/* SHUTTING DOWN */
	logger.Info("Shutting down...")
	err = srv.Shutdown()
	if err != nil {
		logger.Error("server", err.Error())
	}
}
