package services

import (
	"context"
	"fmt"

	"github.com/ipj31/gohoot/internal/database"
	"github.com/ipj31/gohoot/internal/models"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

type QuizzesService struct {
	dbClient       *database.MongoClient
	quizCollection *mongo.Collection
}

func NewQuizzesService(databaseClient *database.MongoClient) *QuizzesService {
	collection := databaseClient.Database.Collection("quizzes")
	return &QuizzesService{
		dbClient:       databaseClient,
		quizCollection: collection,
	}
}

func (qs *QuizzesService) CreateQuiz(quiz models.Quiz) (string, error) {
	result, err := qs.quizCollection.InsertOne(context.Background(), quiz)
	if err != nil {
		return "", err
	}

	quizID := result.InsertedID.(primitive.ObjectID).Hex()
	return quizID, nil
}

func (qs *QuizzesService) GetUserQuizzes(userID string) ([]models.Quiz, error) {
	cursor, err := qs.quizCollection.Find(context.Background(), bson.M{"userId": userID})
	if err != nil {
		return nil, err
	}

	quizzes := make([]models.Quiz, 0, 5)
	for cursor.Next(context.Background()) {
		var quiz models.Quiz
		if err := cursor.Decode(&quiz); err != nil {
			return nil, err
		}
		quizzes = append(quizzes, quiz)
	}
	return quizzes, nil
}

func (qs *QuizzesService) GetQuiz(quizID string) (models.Quiz, error) {
	id, err := primitive.ObjectIDFromHex(quizID)
	if err != nil {
		return models.Quiz{}, fmt.Errorf("error parsing quizID: %w", err)
	}

	var quiz models.Quiz
	err = qs.quizCollection.FindOne(context.Background(), bson.M{"_id": id}).Decode(&quiz)
	if err != nil {
		return models.Quiz{}, err
	}

	return quiz, nil
}
