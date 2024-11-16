package handlers

import (
	"context"
	"fmt"
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

	if err := r.ParseForm(); err != nil {
		http.Error(w, "Failed to parse form", http.StatusBadRequest)
		return
	}

	// Count how many questions we have
	questionCount := 0
	for i := 0; ; i++ {
		if _, exists := r.Form[fmt.Sprintf("question[%d].text", i)]; !exists {
			break
		}
		questionCount++
	}

	questions := make([]models.Question, 0, questionCount)
	for i := 0; i < questionCount; i++ {
		questionText := r.FormValue(fmt.Sprintf("question[%d].text", i))
		if questionText == "" {
			continue
		}

		answers := r.Form[fmt.Sprintf("question[%d].answers[]", i)]
		if len(answers) == 0 {
			continue
		}

		correctAnswer := r.FormValue(fmt.Sprintf("question[%d].correct_answer", i))

		// Validate that the correct answer is one of the answers
		hasCorrectAnswer := false
		for _, answer := range answers {
			if answer == correctAnswer {
				hasCorrectAnswer = true
				break
			}
		}

		if !hasCorrectAnswer {
			http.Error(w, fmt.Sprintf("Question %d: correct answer must be one of the provided answers", i+1), http.StatusBadRequest)
			return
		}

		question := models.Question{
			Question:      questionText,
			Answers:       answers,
			CorrectAnswer: correctAnswer,
		}
		questions = append(questions, question)
	}

	updateQuiz := models.Quiz{
		Name:        r.FormValue("name"),
		Description: r.FormValue("description"),
		Questions:   questions,
	}

	if updateQuiz.Name == "" {
		http.Error(w, "Quiz name is required", http.StatusBadRequest)
		return
	}
	if len(updateQuiz.Questions) == 0 {
		http.Error(w, "Quiz must have at least one question", http.StatusBadRequest)
		return
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

func (uq *UserQuizzes) HandleDeleteQuiz(w http.ResponseWriter, r *http.Request) {
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

	err := uq.quizzesService.DeleteQuiz(userID, quizID)
	if err == mongo.ErrNoDocuments {
		http.Error(w, err.Error(), http.StatusNotFound)
		return
	}
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}
