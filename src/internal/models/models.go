package models

import "go.mongodb.org/mongo-driver/bson/primitive"

type User struct {
	ID             primitive.ObjectID `bson:"_id,omitempty"`
	Email          string             `bson:"email"`
	HashedPassword string             `bson:"hashedPassword"`
}

type Question struct {
	Question      string   `bson:"question"`
	CorrectAnswer string   `bson:"correctAnswer"`
	Answers       []string `bson:"answers"`
}

type Quiz struct {
	ID          primitive.ObjectID `bson:"_id,omitempty"`
	UserID      string             `bson:"userId"`
	Name        string             `bson:"name"`
	Description string             `bson:"description"`
	Questions   []Question         `bson:"question"`
}
