package main

import (
	"context"
	"os"
	"os/signal"
	"syscall"
	"time"

	_ "github.com/google/uuid"
	_ "github.com/lib/pq"
	"song-library/internal/api"
	"song-library/internal/api/handler"
	"song-library/internal/config"
	"song-library/internal/migration"
	"song-library/internal/repository/postgres"
	"song-library/internal/service"
	"song-library/pkg/logger"

	_ "song-library/docs"
)

// @title Онлайн Библиотека Песен API
// @version 1.0
// @description API для управления библиотекой песен

// @host localhost:8080
// @BasePath /api/v1

// @Schemes http
// @Produce json
// @Consume json

func main() {
	cfg, err := config.LoadConfig()
	if err != nil {
		panic("Ошибка загрузки конфигурации: " + err.Error())
	}

	log := logger.NewLogger(cfg.LogLevel)
	log.Info("Запуск приложения")

	db, err := postgres.NewPostgresDB(cfg.DBHost, cfg.DBPort, cfg.DBUser, cfg.DBPassword, cfg.DBName, log)
	if err != nil {
		log.Error("Ошибка подключения к базе данных", "error", err)
		os.Exit(1)
	}

	if err = migration.RunMigrations(db.DB, log); err != nil {
		log.Error("Ошибка выполнения миграций", "error", err)
		os.Exit(1)
	}

	songRepo := postgres.NewSongRepository(db, log)
	apiClient := service.NewExternalAPIClient(cfg.ExternalAPIURL, log)
	songService := service.NewSongService(songRepo, apiClient, log)
	songHandler := handler.NewSongHandler(songService, log)

	router := api.NewRouter(songHandler, log, cfg.Environment)
	router.SetupRoutes()

	server := api.NewServer(router, cfg.ServerPort, log)
	go func() {
		if err = server.Run(); err != nil {
			log.Error("Ошибка запуска HTTP сервера", "error", err)
		}
	}()

	log.Info("Сервис успешно запущен", "port", cfg.ServerPort)

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	log.Info("Получен сигнал остановки, завершение работы...")

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	if err = server.Shutdown(ctx); err != nil {
		log.Error("Ошибка остановки сервера", "error", err)
	}

	log.Info("Сервер успешно остановлен")
}
