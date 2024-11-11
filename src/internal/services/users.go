package services

import (
	"context"

	"github.com/ipj31/gohoot/internal/database"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"golang.org/x/crypto/bcrypt"
)

type User struct {
	ID             primitive.ObjectID `bson:"_id,omitempty"`
	Email          string             `bson:"email"`
	HashedPassword string             `bson:"hashedPassword"`
}

func HashPassword(password string) (string, error) {
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return "", err
	}
	return string(hashedPassword), err
}

func CheckPassword(hashedPassword, password string) error {
	return bcrypt.CompareHashAndPassword([]byte(hashedPassword), []byte(password))
}

type UserService struct {
	dbClient       *database.MongoClient
	userCollection *mongo.Collection
}

func NewUserService(databaseClient *database.MongoClient) *UserService {
	collection := databaseClient.Database.Collection("users")
	return &UserService{
		dbClient:       databaseClient,
		userCollection: collection,
	}
}

func (us *UserService) CreateUser(email, password string) (string, error) {
	hash, err := HashPassword(password)
	if err != nil {
		return "", err
	}

	result, err := us.userCollection.InsertOne(context.Background(), User{Email: email, HashedPassword: hash})
	if err != nil {
		return "", err
	}

	userID := result.InsertedID.(primitive.ObjectID).Hex()
	return userID, nil
}

func (us *UserService) VerifyLogin(email, password string) (string, bool, error) {
	var user User
	err := us.userCollection.FindOne(context.Background(), bson.M{"email": email}).Decode(&user)
	if err != nil {
		return "", false, err
	}

	err = CheckPassword(user.HashedPassword, password)
	if err != nil {
		return "", false, err
	}

	return user.ID.Hex(), true, err
}

func (us *UserService) UniqueEmail(email string) bool {
	var user User
	err := us.userCollection.FindOne(context.Background(), bson.M{"email": email}).Decode(&user)
	if err != nil {
		return false
	}
	return true
}
