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

	// Authentication
	router.Handle("/login", templ.Handler(templates.Login()))
	router.Handle("/register", templ.Handler(templates.Register()))
	router.Handle("/register-submit", registerSubmitRoute)
	router.Handle("/login-submit", loginSubmitRoute)
	router.HandleFunc("/sign-out", handlers.HandleSignOut)

	// Quizzes
	router.Handle("/quizzes", middleware.AuthMiddleware(http.HandlerFunc(userQuizzesRoute.HandleUserQuizzes)))

	// TODO add routes to handle all the operations for quizzes with correct verbs

	router.Handle("GET /quiz/new", middleware.AuthMiddleware(http.HandlerFunc(userQuizzesRoute.HandleNewQuiz)))
	router.Handle("GET /quiz/{id}", middleware.AuthMiddleware(http.HandlerFunc(userQuizzesRoute.HandleGetQuiz)))
	router.Handle("POST /quiz/{id}", middleware.AuthMiddleware(http.HandlerFunc(userQuizzesRoute.HandleSaveQuiz)))

	// TODO add endpoint instead to add a blank question to the quiz
	// this complies with the nature of htmx of the html being the reflection of the server state
	// so if the user needs to create a blank question that should be what is on the server
	// imagine everything as backend

	// router.Handle("GET /quiz/add-question", templ.Handler())

	http.ListenAndServe("", router)
}
