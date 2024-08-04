package handler

import (
	_ "NoteProject/docs"
	"NoteProject/internal/service"
	"fmt"
	"log/slog"
	"net/http"
	"os"
	"path/filepath"
	"strings"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/cors"
	httpSwagger "github.com/swaggo/http-swagger"
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

	c := cors.New(cors.Options{
		AllowedOrigins: []string{fmt.Sprintf("http://%s:%s", os.Getenv("SERVER_HOST"), os.Getenv("SERVER_PORT"))}, // Разрешаем только запросы с данных домена
		AllowedMethods: []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowedHeaders: []string{"Accept", "Authorization", "Content-Type", "Content-Length", "Cache-Control",
			"Connection", "Host", "Origin"},
		AllowCredentials: true,
		MaxAge:           300,
	})

	r.Use(c.Handler)

	r.Get("/swagger/*", httpSwagger.Handler(
		httpSwagger.URL("http://localhost:8080/swagger/doc.json"),
	))

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
		r.Get("/list", h.notesList)
		r.Post("/create", h.noteCreate)
		r.Put("/update", h.noteUpdate)
		r.Delete("/delete", h.noteDelete)
	})

	return r
}

func fileServer(r chi.Router, path string) {
	workDir, _ := os.Getwd()
	root := http.Dir(filepath.Join(workDir, webPath, staticPath))

	if strings.ContainsAny(path, "{}*") {
		//TODO
		panic("FileServer does not permit any URL parameters")
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
