package application

import (
	"net/http"

	"github.com/go-chi/chi"
	"github.com/go-chi/chi/middleware"
	"github.com/pantheon-bolt/bifrost/handler"
)

func loadRoutes() *chi.Mux {
	router := chi.NewRouter()

	router.Use(middleware.Logger)

	router.Get("/", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	})

	router.Route("/apis", loadApiRoutes)

	return router
}

func loadApiRoutes(router chi.Router) {
	apiHandler := &handler.Api{}

	router.Post("/", apiHandler.Create)
	router.Get("/", apiHandler.List)
	router.Get("/{id}", apiHandler.GetByID)
	router.Put("/{id}", apiHandler.UpdateByID)
	router.Delete("/{id}", apiHandler.DeleteByID)
}
