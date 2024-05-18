package routes

import (
	controllers2 "donateapp/pkg/controllers"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/cors"
	"net/http"
)

func Routes() http.Handler {
	router := chi.NewRouter()
	router.Use(middleware.Recoverer)

	setCors(router)

	router.Post("/api/v1/users/auth/register", controllers2.CreateUser)
	router.Post("/api/v1/users/auth/login", controllers2.LoginUser)
	router.Get("/server/v1/users/profile", controllers2.GetProfile)

	router.Post("/api/v1/donations/donate", controllers2.CreateDonation)
	router.Get("/api/v1/donations/donation/{id}", controllers2.GetDonationByID)
	router.Get("/api/v1/donations/", controllers2.GetDonations)
	router.Put("/api/v1/donations/donation/{id}", controllers2.UpdateDonation)
	router.Delete("/api/v1/donations/donation/{id}", controllers2.DeleteDonation)

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
