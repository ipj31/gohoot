package main

import (
	"net/http"

	"github.com/a-h/templ"
	"github.com/ipj31/gohoot/internal/database"
	"github.com/ipj31/gohoot/internal/handlers"
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

	router.HandleFunc("/", handlers.HandleHome)
	router.Handle("/login", templ.Handler(templates.Login()))
	router.Handle("/register", templ.Handler(templates.Register()))
	router.Handle("/register-submit", registerSubmitRoute)
	router.Handle("/login-submit", loginSubmitRoute)
	router.HandleFunc("/sign-out", handlers.HandleSignOut)

	http.ListenAndServe("", router)
}
