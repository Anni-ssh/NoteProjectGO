package main

import (
	"NoteProject/internal/config"
	"NoteProject/internal/service"
	"NoteProject/internal/storage"
	"NoteProject/internal/storage/postgres"
	"NoteProject/internal/storage/redisDB"
	"NoteProject/internal/transport/http-server/handler"
	"NoteProject/internal/transport/http-server/server"
	"NoteProject/pkg/logger"
	"github.com/joho/godotenv"
	_ "github.com/lib/pq"
	"github.com/spf13/viper"
	"log/slog"
	"os"
	"os/signal"
	"syscall"
)

// @title NoteProject
// @version 1.0
// @description API Server for notes workspace
// @host localhost:8080
// @basePath /

// @SecurityDefinitions.apikey ApiKeyAuth
// @in header
// @name Authorization

const ctxTime = 5

func main() {

	if err := godotenv.Load(); err != nil {
		slog.Error("invalid .env file", slog.Any("error", err))
	}

	cfgPath := os.Getenv("CONFIG_PATH")
	cfgName := os.Getenv("CONFIG_NAME")
	err := config.InitConfig(cfgPath, cfgName)

	if err != nil {
		slog.Error("invalid config", slog.Any("error", err))
		panic(err)
	}

	log := logger.SetupLogger(viper.GetString("env"))

	DB, err := postgres.NewPostgresDB(postgres.Config{
		Host:     viper.GetString("db.Host"),
		Port:     viper.GetString("db.Port"),
		Username: os.Getenv("DB_USERNAME"),
		Password: os.Getenv("DB_PASSWORD"),
		DBName:   viper.GetString("db.DBName"),
		SSLMode:  viper.GetString("db.SSLMode"),
	})

	if err != nil {
		log.Error("failed to init PostgresDB", slog.Any("error", err))
		panic(err)
	}

	redisClient, err := redisDB.NewRedisClient(redisDB.Config{
		Addr:     viper.GetString("redisDB.Addr"),
		Password: viper.GetString("redisDB.Password"),
		DB:       viper.GetInt("redisDB.DB"),
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
		if err = srv.Run(viper.GetString("server"), handlers.InitRouter()); err != nil {
			log.Error("error starting server", slog.Any("error", err))
			panic(err)
		}
	}()

	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)

	sig := <-sigChan
	log.Info("Stopped by Admin", "Signal", sig)
}
