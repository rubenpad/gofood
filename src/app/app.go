package app

import (
	"net/http"

	"github.com/go-chi/chi"
	"github.com/go-chi/chi/middleware"
)

type App struct {
	Router *chi.Mux
}

func New() *App {
	ap := &App{}

	router := chi.NewRouter()
	router.Use(middleware.RequestID)
	router.Use(middleware.RealIP)
	router.Use(middleware.Logger)
	router.Use(middleware.Recoverer)

	ap.Router = router
	return ap
}

func (ap *App) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	ap.Router.ServeHTTP(w, r)
}
