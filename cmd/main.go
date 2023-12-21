package main

import (
	"TestProject/internal/lib/logger"
	"TestProject/pages"
	"flag"
	"net/http"
)

func HandleRequest() {
	addr := flag.String("addr", ":8081", "Сетевой адрес HTTP")
	flag.Parse()

	//Создание логера
	slog := logger.SetupLogger("local")

	//Создание нового роутера
	router := http.NewServeMux()

	//Добавление статических данных
	//router.HandleFunc("/static", Static)
	//router.HandleFunc("/static/", Static)
	//
	////Горутины, отлавливают перемещение по сайту
	router.HandleFunc("/", pages.Home)
	router.HandleFunc("/note/save", pages.SaveNote)
	//router.HandleFunc("/note", pages.ShowNote)

	//второй параметр это настроки сервера, запускаем и сразу проверяем на ошибку

	slog.Info("Server starting", "address", "http://127.0.0.1"+*addr)

	err := http.ListenAndServe(*addr, router)
	if err != nil {
		slog.Error("Ошибка запуска сервера")
		panic(err)
	}
}

func main() {

	go HandleRequest()
	select {}

}
