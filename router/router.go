package router

import (
	//...

	"app/config"
	"app/controller"
	middlewares "app/middlewares"
	"app/model"
	"log"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/cors"
	"github.com/go-chi/jwtauth/v5"
	"github.com/go-chi/render"
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
	profileController := controller.NewQueryController[model.Profile]()
	scheduleController := controller.NewScheduleController()

	roomControllerCustom := controller.NewRoomController()
	authController := controller.NewAuthController()

	middlewares := middlewares.NewMiddlewares()

	app.Route("/api/v1", func(router chi.Router) {

		router.Get("/ping", func(w http.ResponseWriter, r *http.Request) {
			render.JSON(w, r, map[string]interface{}{
				"mess": "done",
			})
		})

		router.Route("/public", func(public chi.Router) {
			public.Get("/refresh-admin", authController.CreateAdmin)
			public.Post("/login", authController.Login)
		})

		router.Route("/auth", func(auth chi.Router) {
			auth.Post("/register", authController.Register)
			auth.Post("/auth-face", authController.AuthFace)
			auth.Post("/send-file-auth", authController.SendFileAuth)
			auth.Post("/save-process", authController.SaveProcess)
			auth.Post("/create-socket-auth-face", authController.CreateSocketAuthFace)
		})

		router.Route("/protected", func(protected chi.Router) {
			protected.Use(jwtauth.Verifier(config.GetJWT()))
			protected.Use(jwtauth.Authenticator(config.GetJWT()))
			protected.Use(middlewares.ValidateExpAccessToken())

			protected.Post("/refresh-token", authController.RefreshToken)

			protected.Route("/query", func(query chi.Router) {
				query.Post("/room", roomController.Query)
				query.Post("/field", fieldController.Query)
				query.Post("/schedule", scheduleController.Query)
				query.Post("/department", departmentController.Query)
				query.Post("/profile", profileController.Query)
				query.Post("/profile-department", profileDepartmentController.Query)
			})

			protected.Route("/schedule", func(schedule chi.Router) {
				schedule.Get("/call-medical-file", scheduleController.CallMedicalFile)
				schedule.Post("/pull-medical-file", scheduleController.PullMedicalFile)
				schedule.Post("/transit", scheduleController.Transit)
			})

			protected.Route("/room", func(room chi.Router) {
				room.Get("/call-step", roomControllerCustom.CallStep)
				room.Post("/pull-step", roomControllerCustom.PullStep)
				room.Post("/save-step", roomControllerCustom.SaveStep)
				room.Post("/add-account", roomControllerCustom.AddAccount)
			})
		})
	})

	log.Printf(
		"Server art-pixel starting success! URL: http://%s:%s",
		config.GetAppHost(),
		config.GetAppPort(),
	)

	return app
}
