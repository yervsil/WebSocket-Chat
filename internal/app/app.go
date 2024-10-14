package app

import (
	"context"
	"fmt"
	"log/slog"
	"os"
	"os/signal"
	"strings"
	"sync"
	"time"

	"github.com/yervsil/auth_service/internal/configs"
	"github.com/yervsil/auth_service/internal/producer"
	"github.com/yervsil/auth_service/internal/repository"
	"github.com/yervsil/auth_service/internal/repository/postgres"
	"github.com/yervsil/auth_service/internal/router"
	"github.com/yervsil/auth_service/internal/server"
	"github.com/yervsil/auth_service/internal/service"
	"github.com/yervsil/auth_service/internal/transport/http"
	"github.com/yervsil/auth_service/internal/transport/websocket"
)

const (
	envLocal = "local"
	envDev   = "dev"
	envProd  = "prod"
)

func Run(cfg *configs.Config) {
	log := setupLogger(cfg.Env)

	db, err := postgres.New(cfg)
	if err != nil {
		log.Error(fmt.Sprintf("failed to db connection: %s", err))
		os.Exit(1)
	}

	producer := producer.NewProducer(strings.Split(cfg.Kafka.Brokers, ","), cfg.Kafka.Topic)
	
	repository := repository.New(db)
	service := service.New(repository, log)
	httpHandler := http.New(service, producer, log)

	websocketHandler := websocket.New(producer, log)

	srv := server.New(cfg, router.Routes(httpHandler, websocketHandler))
	go func() {
		if err := srv.Run(); err != nil {
			log.Error(err.Error())
			os.Exit(1)
		}
	}()

	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt)
	<-c

	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
	defer cancel()

	var wg sync.WaitGroup
	wg.Add(3)

	go func() {
        defer wg.Done()
        if err := srv.Stop(ctx); err != nil {
            log.Error(fmt.Sprintf("HTTP server Shutdown with error: %v", err.Error()))
        }
        log.Info("HTTP server gracefully stopped")
    }()

	go func() {
        defer wg.Done()
        if err := db.Close(); err != nil {
            log.Error(fmt.Sprintf("Error closing DB connection: %v", err))
        }
        log.Info("PostgreSQL connection closed")
    }()

	go func() {
        defer wg.Done()
        if err := producer.Close(); err != nil {
            log.Error(fmt.Sprintf("Error closing Kafka producer: %v", err))
        }
        log.Info("Kafka producer closed")
    }()

    wg.Wait()

    log.Info("Application gracefully stopped")
}



func setupLogger(env string) *slog.Logger {
	var log *slog.Logger

	switch env {
	case envLocal:
		log = slog.New(
			slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelDebug}),
		)
	case envDev:
		log = slog.New(
			slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelDebug}),
		)
	case envProd:
		log = slog.New(
			slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelInfo}),
		)
	}

	return log
}