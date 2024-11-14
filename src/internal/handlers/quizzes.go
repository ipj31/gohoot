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

	var quiz models.Quiz
	var err error
	if quizID == "new" {
		// TODO add the blank quiz to the db i think
		quiz = models.NewQuiz()
	} else {
		quiz, err = uq.quizzesService.GetQuiz(quizID)
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
	}

	templates.Quiz(quiz).Render(context.Background(), w)
}
