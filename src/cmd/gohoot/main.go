package main

import (
	"net/http"

	"github.com/a-h/templ"
	"github.com/ipj31/gohoot/internal/database"
	"github.com/ipj31/gohoot/internal/handlers"
	"github.com/ipj31/gohoot/internal/middleware"
	"github.com/ipj31/gohoot/internal/services"
	"github.com/ipj31/gohoot/web/templates"
)

func main() {
	mongoClient, err := database.NewMongoMatchClient("mongodb://admin:password@localhost:27017")
	if err != nil {
		panic(err)
	}

	router := http.NewServeMux()

	userService := services.NewUserService(mongoClient)
	registerSubmitRoute := handlers.NewRegisterSubmit(userService)
	loginSubmitRoute := handlers.NewLoginSubmit(userService)

	quizzesService := services.NewQuizzesService(mongoClient)
	userQuizzesRoute := handlers.NewUserQuizzes(quizzesService)

	router.HandleFunc("/", handlers.HandleHome)
	router.Handle("/login", templ.Handler(templates.Login()))
	router.Handle("/register", templ.Handler(templates.Register()))
	router.Handle("/register-submit", registerSubmitRoute)
	router.Handle("/login-submit", loginSubmitRoute)
	router.HandleFunc("/sign-out", handlers.HandleSignOut)
	router.Handle("/quizzes", middleware.AuthMiddleware(http.HandlerFunc(userQuizzesRoute.HandleUserQuizzes)))
	// TODO add routes to handle all the operations for quizzes with correct verbs

	http.ListenAndServe("", router)
}
