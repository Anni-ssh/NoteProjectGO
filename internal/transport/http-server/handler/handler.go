package handler

import (
	"NoteProject/internal/service"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"log/slog"
	"net/http"
	"os"
	"path/filepath"
	"strings"
)

const (
	webPath    = "web"
	staticPath = "static"
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
	r := chi.NewRouter()
	r.Use(middleware.Recoverer) //recovery из panic
	r.Use(middleware.CleanPath) //исправление путей
	//TO DO
	//router.Use(logger.New(log))

	r.Get("/", h.home)

	//Static files
	fileServer(r, "/")

	// Аутентификация
	r.Route("/auth", func(r chi.Router) {
		r.Post("/sign-up", h.signUp)
		r.Post("/sign-in", h.signIn)
	})
	// Работа с заметками
	r.Route("/note", func(r chi.Router) {
		r.Use(h.authMiddleware)

		r.Get("/workspace", h.noteWorkspace)
		r.Post("/create", h.createNote)
		r.Put("/update", h.updateNote)
		r.Delete("/delete", h.deleteNote)
	})

	r.Get("/session", h.Session)
	return r
}

func fileServer(r chi.Router, path string) {
	workDir, _ := os.Getwd()
	root := http.Dir(filepath.Join(workDir, webPath, staticPath))

	if strings.ContainsAny(path, "{}*") {
		//TODO
		panic("FileServer does not permit any URL parameters.")
	}

	if path != "/" && path[len(path)-1] != '/' {
		r.Get(path, http.RedirectHandler(path+"/", 301).ServeHTTP)
		path += "/"
	}
	path += "*"

	r.Get(path, func(w http.ResponseWriter, r *http.Request) {
		rctx := chi.RouteContext(r.Context())
		pathPrefix := strings.TrimSuffix(rctx.RoutePattern(), "/*")
		fs := http.StripPrefix(pathPrefix, http.FileServer(root))
		fs.ServeHTTP(w, r)
	})
}
