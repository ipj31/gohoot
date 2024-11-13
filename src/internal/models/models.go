package models

import "go.mongodb.org/mongo-driver/bson/primitive"

type User struct {
	ID             primitive.ObjectID `bson:"_id,omitempty"`
	Email          string             `bson:"email"`
	HashedPassword string             `bson:"hashedPassword"`
}

type Question struct {
	Question      string   `bson:"_id,omitempty"`
	CorrectAnswer string   `bson:"correctAnswer"`
	Answers       []string `bson:"answers"`
}

type Quiz struct {
	Name      string     `bson:"name"`
	Questions []Question `bson:"question"`
	UserID    string     `bson:"userId"`
}
