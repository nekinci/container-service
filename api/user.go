package main

import (
	"context"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"time"
)

type User struct {
	Id                 string     `bson:"_id"`
	CreatedAt          time.Time  `bson:"created_at"`
	UpdatedAt          *time.Time `bson:"updated_at"`
	Name               string     `bson:"name"`
	Surname            string     `bson:"surname"`
	Email              string     `bson:"email"`
	Password           string     `bson:"password"`
	RefreshToken       string     `bson:"refresh_token"`
	RefreshTokenExpire *time.Time `bson:"refresh_token_expire"`
}

type userRepository struct {
}

func (userRepo userRepository) GetOne(username, password string) (*User, error) {
	var user = User{}
	filter := bson.D{primitive.E{Key: "email", Value: username}, primitive.E{Key: "password", Value: password}}
	collection := GetMongoClient().Database("ContainerService").Collection("users")
	err := collection.FindOne(context.TODO(), filter).Decode(&user)
	if err != nil {
		return nil, err
	}
	return &user, err
}

func (userRepo userRepository) Get(email string) (*User, error) {
	var user = User{}
	filter := bson.D{primitive.E{Key: "email", Value: email}}
	collection := GetMongoClient().Database("ContainerService").Collection("users")
	err := collection.FindOne(context.TODO(), filter).Decode(&user)
	if err != nil {
		return nil, err
	}
	return &user, nil
}

func (userRepo userRepository) GetByRefreshToken(refreshToken string) (*User, error) {
	var user = User{}
	filter := bson.D{primitive.E{Key: "refresh_token", Value: refreshToken}}
	collection := GetMongoClient().Database("ContainerService").Collection("users")
	err := collection.FindOne(context.TODO(), filter).Decode(&user)
	if err != nil {
		return nil, err
	}
	return &user, nil
}

func (userRepo userRepository) Update(u User) (*User, error) {
	_, updatedUserErr := GetContainerDatabase().Collection("users").UpdateByID(context.TODO(), u.Id, bson.M{"$set": u})
	if updatedUserErr != nil {
		return nil, updatedUserErr
	}
	return &u, nil
}

func (userRepo userRepository) Save(u User) (*User, error) {
	us, err := GetContainerDatabase().Collection("users").InsertOne(context.TODO(), u)
	if err != nil {
		return nil, err
	}
	u.Id = us.InsertedID.(string)
	return &u, nil
}

func newUserRepository() userRepository {
	return userRepository{}
}
