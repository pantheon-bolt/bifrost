package application

import (
	"net/http"

	"github.com/go-chi/chi"
	"github.com/go-chi/chi/middleware"
	"github.com/pantheon-bolt/bifrost/handler"
	"github.com/pantheon-bolt/bifrost/repository/api"
)

func (a *App) loadRoutes() {
	router := chi.NewRouter()

	router.Use(middleware.Logger)

	router.Get("/", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	})

	router.Route("/apis", a.loadApiRoutes)

	a.router = router
}

func (a *App) loadApiRoutes(router chi.Router) {

	apiHandler := &handler.Api{
		Repo: &api.RedisRepo{
			Client: a.rdb,
		},
	}

	router.Post("/", apiHandler.Create)
	router.Get("/", apiHandler.List)
	router.Get("/{id}", apiHandler.GetByID)
	router.Put("/{id}", apiHandler.UpdateByID)
	router.Delete("/{id}", apiHandler.DeleteByID)
}
