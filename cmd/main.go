package main

import (
	"TestProject/config"
	"TestProject/internal/lib/dataBaseSQL"
	"TestProject/internal/lib/dataBaseSQL/ServSQLite"
	"TestProject/internal/lib/logger"
	"TestProject/pages"
	"context"
	"flag"
	"fmt"
	_ "github.com/mattn/go-sqlite3"
	"net/http"
	"os"
	"os/signal"
	"syscall"
)

func HandleRequest(app pages.Application) {
	//Получаем порт через флаги при запуске
	addr := flag.String("addr", ":8081", "Сетевой адрес HTTP")
	flag.Parse()

	//Создание нового роутера
	router := http.NewServeMux()

	//Handle
	router.HandleFunc("/", app.Home)
	router.HandleFunc("/note/save", app.SaveNote)

	//Static
	//router.HandleFunc("/static", Static)
	//router.HandleFunc("/static/", Static)

	//Создание структуры сервера
	serv := &http.Server{
		Addr:    *addr,
		Handler: router,
	}

	//Инфо про сервер
	app.Slog.Info("Server starting", "address", "http://127.0.0.1"+*addr)
	//Запуск сервера - это горутина
	err := serv.ListenAndServe()
	if err != nil {
		app.Slog.Error("Ошибка запуска сервера")
		panic(err)
	}

}

func main() {
	//Создание конфига
	//Panic дальнейшая работа программы без данных компонентов невозможна
	cfg, err := config.CreateCfg("config.json")
	if err != nil {
		fmt.Println(err)
		panic(err)
	}

	//Создание логера
	slog := logger.SetupLogger(cfg.Env)

	//Создание пула соединений в БД
	db, err := dataBaseSQL.OpenDB(cfg.DataBasePath)
	if err != nil {
		slog.Error("Ошибка открытия соединения с БД")
		panic(err)
	}

	//FIX ME Создание контекста
	Ctx := context.Background()

	//Создание необходимых таблиц
	err = ServSQLite.DataBaseSQLiteNote{Storage: db}.CreateNotesTable(Ctx)
	if err != nil {
		slog.Error("Err Open Data base")
		panic(err)
	}

	app := pages.Application{
		DB:     db,
		Slog:   slog,
		Ctx:    Ctx,
		Config: cfg,
	}

	go HandleRequest(app)

	//GracefulShutdown
	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)

	select {
	case sig := <-sigChan:
		slog.Info("Stopped by Admin", "Signal", sig)
		return
	}

}
