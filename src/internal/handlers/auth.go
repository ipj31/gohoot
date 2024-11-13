package handlers

import (
	"context"
	"fmt"
	"net/http"

	"github.com/ipj31/gohoot/internal/services"
	"github.com/ipj31/gohoot/web/templates"
)

type RegisterSubmit struct {
	userService *services.UserService
}

func NewRegisterSubmit(userService *services.UserService) *RegisterSubmit {
	return &RegisterSubmit{
		userService: userService,
	}
}

func (rs *RegisterSubmit) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	email := r.FormValue("email")
	password := r.FormValue("password")
	confirmPassword := r.FormValue("confirm-password")

	uniqueEmail, err := rs.userService.UniqueEmail(email)

	if err != nil {
		args := templates.RegisterFormArgs{
			Email:           "",
			Password:        password,
			ConfirmPassword: confirmPassword,
			EmailError:      "Error checking for unique email",
		}
		templates.RegisterForm(args).Render(context.Background(), w)
		http.Error(w, "error checking for unique email", http.StatusInternalServerError)
		return
	}

	if uniqueEmail {
		args := templates.RegisterFormArgs{
			Email:           "",
			Password:        password,
			ConfirmPassword: confirmPassword,
			EmailError:      "Email already exists",
		}
		templates.RegisterForm(args).Render(context.Background(), w)
		return
	}

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

	userID, err := rs.userService.CreateUser(email, password)
	if err != nil {
		args := templates.RegisterFormArgs{
			Email:                email,
			Password:             password,
			ConfirmPassword:      "",
			ConfirmPasswordError: "Error creating user",
		}
		templates.RegisterForm(args).Render(context.Background(), w)
		return
	}

	fmt.Println("Register successful for", email)

	jwtToken, err := services.GenerateJWT(userID)
	if err != nil {
		fmt.Println(err)
		http.Error(w, "error creating jwt token", http.StatusInternalServerError)
		return
	}
	services.SetTokenCookie(w, jwtToken)

	w.Header().Add("HX-Redirect", "/")
}

type LoginSubmit struct {
	userService *services.UserService
}

func NewLoginSubmit(userService *services.UserService) *LoginSubmit {
	return &LoginSubmit{
		userService: userService,
	}
}

func (uc *LoginSubmit) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	email := r.FormValue("email")
	password := r.FormValue("password")

	userID, valid, err := uc.userService.VerifyLogin(email, password)
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
	jwtToken, err := services.GenerateJWT(userID)
	if err != nil {
		http.Error(w, "error creating jwt token", http.StatusInternalServerError)
		return
	}
	services.SetTokenCookie(w, jwtToken)

	w.Header().Add("HX-Redirect", "/")
}

func HandleHome(w http.ResponseWriter, r *http.Request) {
	if r.URL.String() != "/" {
		// i could set up a toast that this sets to show the error
		http.NotFound(w, r)
		return
	}

	templates.Home(services.UserLoggedIn(r)).Render(context.Background(), w)
}

func HandleSignOut(w http.ResponseWriter, r *http.Request) {
	services.ClearTokenCookie(w)
	templates.LoginButton().Render(context.Background(), w)
}
