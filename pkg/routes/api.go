package routes

import (
	"donateapp/pkg/controllers"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/cors"
	"net/http"
)

func Routes() http.Handler {
	router := chi.NewRouter()
	router.Use(middleware.Recoverer)

	setCors(router)

	router.Post("/api/v1/users/auth/register", controllers.CreateUser)
	router.Post("/api/v1/users/auth/login", controllers.LoginUser)
	router.Get("/server/v1/users/profile", controllers.GetProfile)

	router.Post("/api/v1/donations/donate", controllers.CreateDonation)
	router.Get("/api/v1/donations/donation/{id}", controllers.GetDonationByID)
	router.Get("/api/v1/donations/", controllers.GetDonations)
	router.Put("/api/v1/donations/donation/{id}", controllers.UpdateDonation)
	router.Delete("/api/v1/donations/donation/{id}", controllers.DeleteDonation)

	return router
}

func setCors(r *chi.Mux) {
	r.Use(cors.Handler(cors.Options{
		AllowedOrigins: []string{"http://*", "https://*"},
		AllowedMethods: []string{"GET", "POST", "PUT", "PATCH", "DELETE", "OPTIONS"},
		AllowedHeaders: []string{
			"Accept", "Authorization", "Content-Type", "X-CSRF-Token",
		},
		ExposedHeaders:   []string{"Link"},
		AllowCredentials: true,
		MaxAge:           300,
	}))
}
