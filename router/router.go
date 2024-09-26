package router

import (
	//...

	"app/config"
	"app/controller"
	"app/model"
	"log"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/cors"
)

func AppRouter() http.Handler {
	app := chi.NewRouter()

	// A good base middleware stack
	app.Use(middleware.RequestID)
	app.Use(middleware.RealIP)
	app.Use(middleware.Logger)
	app.Use(middleware.Recoverer)

	app.Use(cors.Handler(cors.Options{
		AllowedOrigins:   []string{"*"},
		AllowedMethods:   []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowedHeaders:   []string{"Accept", "Authorization", "Content-Type", "X-CSRF-Token"},
		ExposedHeaders:   []string{"Link"},
		AllowCredentials: false,
		MaxAge:           300,
	}))

	departmentController := controller.NewQueryController[model.Department]()
	fieldController := controller.NewQueryController[model.Field]()
	roomController := controller.NewQueryController[model.Room]()
	profileDepartmentController := controller.NewQueryController[model.ProfileDepartment]()

	app.Route("/api/v1", func(r chi.Router) {
		r.Post("/department", departmentController.Query)
		r.Post("/field", fieldController.Query)
		r.Post("/room", roomController.Query)
		r.Post("/profile-department", profileDepartmentController.Query)
	})

	log.Printf(
		"Server art-pixel starting success! URL: http://%s:%s",
		config.GetAppHost(),
		config.GetAppPort(),
	)

	return app
}
