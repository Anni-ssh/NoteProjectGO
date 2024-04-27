package handler

import (
	"NoteProject/internal/service"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"log/slog"
	"net/http"
)

type Handler struct {
	services *service.Service
	Logs     *slog.Logger
}

func NewHandler(services *service.Service) *Handler {
	return &Handler{services: services}
}
func (h *Handler) InitLogger(l *slog.Logger) {
	h.Logs = l
}
func (h *Handler) InitRouter() *chi.Mux {
	router := chi.NewRouter()
	router.Use(middleware.Recoverer)

	//TO DO
	//router.Use(logger.New(log))
	router.Get("/", h.home)

	//TO DO
	//Static file
	router.Route("/static", func(r chi.Router) {
		r.Get("/", h.static)
	})

	fs := http.FileServer(http.Dir("./web"))
	router.Handle("/web/*", http.StripPrefix("/web/", fs))

	// Аутентификация
	router.Route("/auth", func(r chi.Router) {
		r.Post("/sign-up", h.signUp)
		r.Post("/sign-in", h.signIn)
	})
	router.Route("/note", func(r chi.Router) {
		r.Use(h.authMiddleware)
		r.Get("/read", h.readNote)
		r.Post("/create", h.createNote)
		r.Post("/update", h.updateNote)
		r.Post("/delete", h.deleteNote)
	})

	return router
}
