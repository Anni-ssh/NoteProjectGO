package handler

import (
	"net/http"
)

// @Summary Home Возвращает домашнюю страницу
// @Description Обрабатывает запрос на получение домашней страницы.
// @Tags Home
// @Accept  json
// @Produce  json
// @Success 200 {object} Response "Успешный ответ"
// @Failure 500 {object} ErrorResponse "Внутренняя ошибка сервера"
// @Router /home [get]

func (h *Handler) home(w http.ResponseWriter, r *http.Request) {
	const op = "handler.home"
	// log := h.Logs.With(slog.String("operation", op))

	_, _ = w.Write([]byte("It is Home"))

	// err := pages.Home(w)
	// if err != nil {
	// 	NewErrResponse(w, http.StatusInternalServerError, "Internal Server Error")
	// 	log.Error("Server error", slog.Any("error", err))
	// 	return
	// }

}
