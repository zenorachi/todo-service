package app

import (
	_ "github.com/zenorachi/todo-service/docs/swagger"
	"github.com/zenorachi/todo-service/internal/config"
	"github.com/zenorachi/todo-service/internal/database"
	"github.com/zenorachi/todo-service/internal/repository"
	"github.com/zenorachi/todo-service/internal/server"
	"github.com/zenorachi/todo-service/internal/service"
	"github.com/zenorachi/todo-service/internal/transport"
	"github.com/zenorachi/todo-service/pkg/auth"
	"github.com/zenorachi/todo-service/pkg/database/postgres"
	"github.com/zenorachi/todo-service/pkg/hash"
	"github.com/zenorachi/todo-service/pkg/logger"
	"os"
	"os/signal"
	"syscall"
)

// @title           			TO-DO service
// @version         			1.0
// @description     			This is a service that implements the functions of a daily planner

// @contact.name   				Maksim Sonkin
// @contact.email  				msonkin33@gmail.com

// @host      					localhost:8080
// @BasePath  					/

// @securityDefinitions.apikey  Bearer
// @in 						    header
// @name 					    Authorization
// @description					Type "Bearer" followed by a space and JWT token.

func Run(cfg *config.Config) {
	/* DO MIGRATIONS */
	if database.DoMigrations(&cfg.DB) != nil {
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
	tokenManager := auth.NewManager(cfg.Auth.Secret)

	/* INIT SERVICES & DEPS */
	services := service.New(service.Deps{
		Repos:           repository.New(db),
		Hasher:          hash.NewSHA1Hasher(cfg.Auth.Salt),
		TokenManager:    tokenManager,
		AccessTokenTTL:  cfg.Auth.AccessTokenTTL,
		RefreshTokenTTL: cfg.Auth.RefreshTokenTTL,
	})

	/* INIT HTTP HANDLER */
	handler := transport.NewHandler(services, tokenManager)

	/* INIT & RUN HTTP SERVER */
	srv := server.New(cfg, handler.InitRoutes())
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
	if srv.Shutdown() != nil {
		logger.Error("server", err.Error())
	}
}
