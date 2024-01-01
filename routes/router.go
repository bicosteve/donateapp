package routes

import (
	"donateapp/controllers"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/cors"
	"net/http"
)

func Routes() http.Handler {
	router := chi.NewRouter()

	router.Use(middleware.Recoverer)
	router.Use(cors.Handler(cors.Options{
		AllowedOrigins: []string{"http://*", "https://*"},
		AllowedMethods: []string{"GET", "POST", "PUT", "PATCH", "DELETE", "OPTIONS"},
		AllowedHeaders: []string{
			"Accept", "Authorization", "Content-Type", "X-CSRF-Token",
		},
		ExposedHeaders:   []string{"Link"},
		AllowCredentials: true,
		MaxAge:           300,
	}))

	// User Routes
	router.Post("/api/v1/users/auth/register", controllers.CreateUser)
	router.Post("/api/v1/users/auth/login", controllers.LoginUser)
	router.Get("/api/v1/users/profile", controllers.GetProfile)

	// Donation Routes /api/v1/donations/donate
	router.Post("/api/v1/donations/donate", controllers.CreateDonation)
	// Donation /api/v1/donations/donation/{id}
	router.Get("/api/v1/donations/donation/{id}", controllers.GetDonationByID)

	return router
}
