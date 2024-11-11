package main

import (
	"context"
	"fmt"
	"net/http"

	"github.com/a-h/templ"
	"github.com/ipj31/gohoot/internal/database"
	"github.com/ipj31/gohoot/internal/services"
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

type RegisterSubmit struct {
	userService *services.UserService
}

func (rs *RegisterSubmit) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	email := r.FormValue("email")
	password := r.FormValue("password")
	confirmPassword := r.FormValue("confirm-password")

	if password != confirmPassword {
		args := templates.RegisterFormArgs{
			Email:                email,
			Password:             password,
			ConfirmPassword:      "",
			ConfirmPasswordError: "Passwords do not match!",
		}
		templates.RegisterForm(args).Render(context.Background(), w)
		return
	}

	err := rs.userService.CreateUser(email, password)
	if err != nil {
		args := templates.RegisterFormArgs{
			Email:                email,
			Password:             password,
			ConfirmPassword:      "",
			ConfirmPasswordError: "Passwords do not match!",
		}
		templates.RegisterForm(args).Render(context.Background(), w)
		return
	}

	fmt.Println("Register successful for", email)
	w.Header().Add("HX-Redirect", "/")
}

type LoginSubmit struct {
	userService *services.UserService
}

func (uc *LoginSubmit) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	email := r.FormValue("email")
	password := r.FormValue("password")

	valid, err := uc.userService.VerifyLogin(email, password)
	// TODO don't love the flow, no differentiation between invalid login and error ocurred
	if !valid || err != nil {
		args := templates.LoginFormArgs{
			Email:         email,
			Password:      password,
			PasswordError: "Invalid email/password",
		}
		templates.LoginForm(args).Render(context.Background(), w)
		return
	}

	fmt.Println("Login successful for", email)
	w.Header().Add("HX-Redirect", "/")
}

func main() {
	mongoClient, err := database.NewMongoMatchClient("mongodb://admin:password@localhost:27017")
	if err != nil {
		panic(err)
	}

	userService := services.NewUserService(mongoClient)

	registerSubmitRoute := &RegisterSubmit{
		userService: userService,
	}

	loginSubmitRoute := &LoginSubmit{
		userService: userService,
	}

	http.HandleFunc("/", handleHome)
	http.Handle("/login", templ.Handler(templates.Login()))
	http.Handle("/register", templ.Handler(templates.Register()))
	http.Handle("/register-submit", registerSubmitRoute)
	http.Handle("/login-submit", loginSubmitRoute)

	http.ListenAndServe("", nil)
}
