package pages

import (
	"TestProject/internal/helperFunc"
	"encoding/json"
	"net/http"
)

// SaveNote страничка сохранения заметки.
func (app *Application) ChangeNote(w http.ResponseWriter, r *http.Request) {
	//Проверяем метод
	if r.Method != http.MethodPost {
		// Вызываем метод w.WriteHeader() для возвращения статус-кода 405
		// и вызывается метод w.Write() для возвращения тела-ответа с текстом "Метод запрещен".
		w.WriteHeader(405)
		w.Header().Set("Allow", http.MethodPost)
		//FIX ME
		_, err := w.Write([]byte("GET-Метод запрещен!"))
		if err != nil {
			return
		}
		// Завершаем работу функции - "return", чтобы
		// последующий код не выполнялся.
		return
	}

	if r.Method == http.MethodPost {
		decoder := json.NewDecoder(r.Body)
		var data map[string]string
		err := decoder.Decode(&data)
		//FIX ME
		if err != nil {
			http.Error(w, "Полученные данные содержат ошибку или пустые", http.StatusBadRequest)
			return
		}

		title := data["Title"]
		text := data["Text"]

		//Отправляем клиенту ответ, что успешно обработали запрос. Формат Json.
		responseData := map[string]string{"message": "Success"}
		w.WriteHeader(http.StatusOK)
		w.Header().Set("Content-Type", "application/json")
		//FIX ME err
		err = json.NewEncoder(w).Encode(responseData)

		if err != nil {
			app.Slog.Error("Не удалось отправить ответ клиенту")
			app.ServerError(w, err)
			return
		}

		noteBody := helperFunc.ConvStrInNote(title, text)

		err = helperFunc.SendNoteToDB(app.DB, app.Ctx, noteBody)
		if err != nil {
			//FIX ME !!!!
			app.Slog.Error("Ошибка. Не удалось отправить данные в БД")
			app.ServerError(w, err)
			return
		}
		http.Redirect(w, r, "/", http.StatusTemporaryRedirect)
	}
}
