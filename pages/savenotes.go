package pages

import (
	"log"
	"net/http"
)

// SaveNote страничка сохранения заметки.
func SaveNote(w http.ResponseWriter, r *http.Request) {

	//Проверяем метод
	if r.Method != http.MethodPost {
		w.Header().Set("Allow", http.MethodPost)
		// Вызываем метод w.WriteHeader() для возвращения статус-кода 405
		// и вызывается метод w.Write() для возвращения тела-ответа с текстом "Метод запрещен".
		w.WriteHeader(405)
		w.Write([]byte("GET-Метод запрещен!"))
		// Затем мы завершаем работу функции вызвав "return", чтобы
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

	http.Redirect(w, r, "/", http.StatusTemporaryRedirect)
}
