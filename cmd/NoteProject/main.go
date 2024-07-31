package main

import (
	"NoteProject/internal/service"
	"NoteProject/internal/storage"
	"NoteProject/internal/storage/postgres"
	"NoteProject/internal/storage/redisDB"
	"NoteProject/internal/transport/http-server/handler"
	"NoteProject/internal/transport/http-server/server"
	"NoteProject/pkg/logger"
	migrations "NoteProject/pkg/migration"
	"log/slog"
	"os"
	"os/signal"
	"strconv"
	"syscall"
	"time"

	"github.com/joho/godotenv"
	_ "github.com/lib/pq"
)

// @title NoteProject
// @version 1.0
// @description API Server for notes workspace
// @host localhost:8080
// @basePath /

// @SecurityDefinitions.apikey ApiKeyAuth
// @in header
// @name Authorization

func main() {

	// Загрузка переменных окружения из файла .env
	if err := godotenv.Load(); err != nil {
		slog.Error("invalid .env file", slog.Any("error", err))
	}

	// Инициализация логгера
	log := logger.SetupLogger(os.Getenv("ENV"))

	// Чтение параметров базы данных из переменных окружения
	dbHost := os.Getenv("DB_HOST")
	dbPort := os.Getenv("DB_PORT")
	dbUsername := os.Getenv("DB_USERNAME")
	dbPassword := os.Getenv("DB_PASSWORD")
	dbName := os.Getenv("DB_NAME")
	dbSSLMode := os.Getenv("DB_SSLMODE")

	attempts, err := strconv.Atoi(os.Getenv("DB_ATTEMPTS"))

	if err != nil {
		log.Error("failed conv str to int", slog.Any("error", err))
		panic(err)
	}

	delay, err := strconv.Atoi(os.Getenv("DB_DELAY"))

	if err != nil {
		log.Error("failed conv str to int", slog.Any("error", err))
		panic(err)
	}

	DB, err := postgres.NewPostgresDB(postgres.Config{
		Host:     dbHost,
		Port:     dbPort,
		Username: dbUsername,
		Password: dbPassword,
		DBName:   dbName,
		SSLMode:  dbSSLMode,
	}, attempts, time.Duration(delay))

	if err != nil {
		log.Error("failed to init PostgresDB", slog.Any("error", err))
		panic(err)
	}

	// Миграции БД
	err = migrations.RunMigrations(DB)
	if err != nil {
		log.Error("Failed to create create migrations", slog.Any("error", err))
		panic(err)
	}

	log.Info("Migrations applied successfully!")

	// Чтение параметров Redis из переменных окружения
	redisAddr := os.Getenv("REDIS_ADDR")
	redisPassword := os.Getenv("REDIS_PASSWORD")
	redisNum, err := strconv.Atoi(os.Getenv("REDIS_DB_NUM"))

	if err != nil {
		log.Error("failed conv str to int", slog.Any("error", err))
		panic(err)
	}

	redisClient, err := redisDB.NewRedisClient(redisDB.Config{
		Addr:     redisAddr,
		Password: redisPassword,
		DB:       redisNum,
	})

	if err != nil {
		log.Error("failed to init RedisDB", slog.Any("error", err))
		panic(err)
	}

	repositories := storage.NewStorage(DB, redisClient)
	services := service.NewService(repositories)
	handlers := handler.NewHandler(services)

	handlers.InitLogger(log)

	srv := &server.Server{}
	log.Info("Starting server...")

	go func() {
		if err = srv.Run(os.Getenv("SERVER_HOST")+":"+os.Getenv("SERVER_PORT"), handlers.InitRouter()); err != nil {
			log.Error("error starting server", slog.Any("error", err))
			panic(err)
		}
	}()

	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)

	sig := <-sigChan
	log.Info("Stopped by Admin", "Signal", sig)

}
