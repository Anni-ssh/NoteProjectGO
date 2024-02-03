package pages

import (
	"TestProject/internal/helperFunc"
	"fmt"
	"log"
	"net/http"
	"os"
)

// SaveNote страничка сохранения заметки.
func (app *Application) SaveNote(w http.ResponseWriter, r *http.Request) {

	//Проверяем метод
	if r.Method != http.MethodPost {
		w.Header().Set("Allow", http.MethodPost)
		// Вызываем метод w.WriteHeader() для возвращения статус-кода 405
		// и вызывается метод w.Write() для возвращения тела-ответа с текстом "Метод запрещен".
		w.WriteHeader(405)
		//FIX ME
		_, err := w.Write([]byte("GET-Метод запрещен!"))
		if err != nil {
			return
		}
		// Завершаем работу функции - "return", чтобы
		// последующий код не выполнялся.
		return
	}

	err := r.ParseForm()
	if err != nil {
		http.Error(w, "Ошибка при обработке данных формы", http.StatusInternalServerError)
		// Логирование ошибки также может быть полезным для отслеживания проблем
		log.Println("Ошибка при парсинге формы:", err)
		return
	}

	title := r.FormValue("note-title")
	text := r.FormValue("note-text")

	noteBody := helperFunc.ConvStrInNote(title, text)

	err = helperFunc.SendNoteToDB(app.DB, app.Ctx, noteBody)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	http.Redirect(w, r, "/", http.StatusTemporaryRedirect)
}
