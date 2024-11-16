package handlers

import (
	"context"
	"net/http"

	"github.com/ipj31/gohoot/internal/models"
	"github.com/ipj31/gohoot/internal/services"
	"github.com/ipj31/gohoot/web/templates"
	"go.mongodb.org/mongo-driver/mongo"
)

type UserQuizzes struct {
	quizzesService *services.QuizzesService
}

func NewUserQuizzes(quizzesService *services.QuizzesService) *UserQuizzes {
	return &UserQuizzes{
		quizzesService: quizzesService,
	}
}

func (uq *UserQuizzes) HandleUserQuizzes(w http.ResponseWriter, r *http.Request) {
	userID, ok := r.Context().Value("user_id").(string)
	if !ok {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	quizzes, err := uq.quizzesService.GetUserQuizzes(userID)
	if err != nil {
		http.Error(w, "error fetching quizzes", http.StatusInternalServerError)
	}

	templates.UserQuizzes(quizzes).Render(r.Context(), w)
}

func (uq *UserQuizzes) HandleNewQuiz(w http.ResponseWriter, r *http.Request) {
	userID, ok := r.Context().Value("user_id").(string)
	if !ok {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	quiz := models.Quiz{
		Name:      "New Quiz",
		UserID:    userID,
		Questions: []models.Question{},
	}
	quizID, err := uq.quizzesService.CreateQuiz(quiz)
	if err != nil {
		http.Error(w, "error creating quiz", http.StatusInternalServerError)
		return
	}

	w.Header().Add("HX-Redirect", "/quiz/"+quizID)
}

func (uq *UserQuizzes) HandleGetQuiz(w http.ResponseWriter, r *http.Request) {
	userID, ok := r.Context().Value("user_id").(string)
	if !ok {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	quizID := r.PathValue("id")
	if quizID == "" {
		http.Error(w, "error parsing quiz id from path params", http.StatusNotFound)
		return
	}

	quiz, err := uq.quizzesService.GetQuiz(quizID)
	if err == mongo.ErrNoDocuments {
		http.NotFound(w, r)
		return
	}
	if err != nil {
		http.Error(w, "error fetching quiz from db", http.StatusInternalServerError)
		return
	}

	if quiz.UserID != userID {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	templates.Quiz(quiz).Render(context.Background(), w)
}

func (uq *UserQuizzes) HandleSaveQuiz(w http.ResponseWriter, r *http.Request) {
	userID, ok := r.Context().Value("user_id").(string)
	if !ok {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	quizID := r.PathValue("id")
	if quizID == "" {
		http.Error(w, "error parsing quiz id from path params", http.StatusNotFound)
		return
	}

	// TODO figure out how to handle all the questions
	updateQuiz := models.Quiz{
		Name:        r.FormValue("name"),
		Description: r.FormValue("description"),
		Questions:   []models.Question{},
	}

	err := uq.quizzesService.UpdateQuiz(userID, quizID, updateQuiz)
	if err == services.ErrUnauthorized {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}
