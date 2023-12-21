package pages

import (
	"TestProject/internal/lib/session"
	"errors"
	"fmt"
	"html/template"
	"net/http"
)

// Главная страница
func Home(w http.ResponseWriter, r *http.Request) {

	//Проверка наличия куки в запросе
	_, err := r.Cookie("NotesX")

	//Если нет, то создаём и отправляем
	if errors.Is(err, http.ErrNoCookie) {
		cookie := http.Cookie{}

		err = session.SetSession("NotesX", 16, 24, &cookie)

		if err != nil {
			http.Error(w, "Internal Server Error", http.StatusInternalServerError)
			return
		}
		http.SetCookie(w, &cookie)
	}

	//Проверка наличия страницы на сайте
	if r.URL.Path != "/" {
		http.NotFound(w, r)
		return
	}

	path := []string{"Home.html"}

	t, err := template.ParseFiles(path...)

	if err != nil {
		fmt.Println("Ошибка отправки HTML")
		panic(err)
	}

	//Variables := serverTypes.PageVariables{Title: "Notes", Hour: time.Now().Hour()}
	err = t.Execute(w, "HomePage")

	if err != nil {
		fmt.Println("Ошибка исполнения HTML")
		panic(err)
	}

}
