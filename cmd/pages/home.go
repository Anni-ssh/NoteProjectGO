package pages

import (
	"TestProject/internal/config"
	"TestProject/internal/helperFunc"
	"TestProject/internal/lib/session"
	"context"
	"database/sql"
	"errors"
	"html/template"
	"log/slog"
	"net/http"
)

type Application struct {
	DB     *sql.DB
	Err    error
	Slog   *slog.Logger
	Ctx    context.Context
	Config *config.StartupConfig
}

func (app *Application) Home(w http.ResponseWriter, r *http.Request) {

	//FIX ME ДОДЕЛАТЬ СЕССИИ
	//Проверка наличия куки в запросе
	_, err := r.Cookie("NotesX")

	//Если нет, то создаём и отправляем
	if errors.Is(err, http.ErrNoCookie) {
		cookie := http.Cookie{}

		//Создание токена
		err = session.SetSession("NotesX", 16, 24, &cookie)

		if err != nil {
			app.ServerError(w, err)
			return
		}
		http.SetCookie(w, &cookie)
	}

	//Проверка наличия страницы на сайте
	if r.URL.Path != "/" {
		app.NotFound(w)
		return
	}
	//т.к. всё парсится первым нужно обязательно home
	paths := []string{"./web/templates/home.html", "./web/templates/note.html", "./web/templates/header.html"}

	t, err := template.ParseFiles(paths...)

	//FIX ME
	if err != nil {
		if err != nil {
			app.Slog.Error("Ошибка отправки HTML")
			app.ServerError(w, err)
			return
		}
	}

	notesList, err := helperFunc.СonvExtractedNotesData(app.DB, app.Ctx)

	if err != nil {
		app.Slog.Error("Ошибка запроса БД")
		app.ServerError(w, err)
		return
	}

	err = t.Execute(w, notesList)
	if err != nil {
		app.Slog.Error("Ошибка исполнения HTML")
		app.ServerError(w, err)
		return
	}

}
