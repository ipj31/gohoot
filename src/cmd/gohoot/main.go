package main

import (
	"context"
	"net/http"

	"github.com/ipj31/gohoot/internal/database"
	"github.com/ipj31/gohoot/web/templates"
)

func handleHome(w http.ResponseWriter, r *http.Request) {
	if r.URL.String() != "/" {
		// i could set up a toast that this sets to show the error
		http.NotFound(w, r)
		return
	}

	templates.Home().Render(context.Background(), w)
}

func handleRegisterSubmit(w http.ResponseWriter, r *http.Request) {
	args := templates.RegisterFormArgs{
		Email:                r.FormValue("email"),
		Password:             r.FormValue("password"),
		ConfirmPassword:      r.FormValue("confirm-password"),
		ConfirmPasswordError: "Passwords do not match!",
	}

	templates.RegisterForm(args).Render(context.Background(), w)
}

func main() {
	_, err := database.NewMongoMatchClient("mongodb://admin:password@localhost:27017")
	if err != nil {
		panic(err)
	}

	// http.HandleFunc("/", handleHome)
	// http.Handle("/login", templ.Handler(templates.Login()))
	// http.Handle("/register", templ.Handler(templates.Register()))
	// http.HandleFunc("/register-submit", handleRegisterSubmit)

	// http.ListenAndServe("", nil)
}
