package main

import (
	"context"
	"fmt"
	"net/http"
	"time"

	"github.com/a-h/templ"
	"github.com/ipj31/gohoot/internal/database"
	"github.com/ipj31/gohoot/internal/services"
	"github.com/ipj31/gohoot/web/templates"
)

// TODO code cleanup

func setTokenCookie(w http.ResponseWriter, token string) {
	http.SetCookie(w, &http.Cookie{
		Name:     "token",
		Value:    token,
		Expires:  time.Now().Add(24 * time.Hour),
		HttpOnly: true,
		Secure:   false, // TODO change to true when using https
	})
}

func AuthMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		cookie, err := r.Cookie("token")
		if err != nil {
			http.Error(w, "Unauthorized", http.StatusUnauthorized)
			return
		}

		userID, err := services.ValidateJWT(cookie.Value)
		if err != nil {
			http.Error(w, "Unauthorized", http.StatusUnauthorized)
			return
		}

		ctx := context.WithValue(r.Context(), "user_id", userID)

		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

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

	if rs.userService.UniqueEmail(email) {
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
	setTokenCookie(w, jwtToken)

	w.Header().Add("HX-Redirect", "/")
}

type LoginSubmit struct {
	userService *services.UserService
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
	setTokenCookie(w, jwtToken)

	w.Header().Add("HX-Redirect", "/")
}

func authTesting(w http.ResponseWriter, r *http.Request) {
	userID, ok := r.Context().Value("user_id").(string)
	if !ok {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	fmt.Println(userID)
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
	http.Handle("/auth-test", AuthMiddleware(http.HandlerFunc(authTesting)))

	http.ListenAndServe("", nil)
}
