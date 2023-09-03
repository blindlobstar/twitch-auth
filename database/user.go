package database

import (
	"context"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type User struct {
	Id       primitive.ObjectID `bson:"_id"`
	TwitchId string             `bson:"twitch_id"`
}

type UserRepo interface {
	CreateUser(ctx context.Context, tid string) (*User, error)
	GetUsers(ctx context.Context, u *User) (*[]User, error)
}

func (db *MongoDB) CreateUser(ctx context.Context, twitch_id string) (*User, error) {
	u := User{
		Id:       primitive.NewObjectID(),
		TwitchId: twitch_id,
	}
	_, err := db.DB.Collection("user").InsertOne(ctx, u)
	return &u, err
}

func (db *MongoDB) GetUsers(ctx context.Context, m *User) (*[]User, error) {

	var result []User
	filter := bson.D{{Key: "twitch_id", Value: m.TwitchId}}
	cursor, err := db.DB.Collection("user").Find(ctx, filter)
	if err != nil {
		return nil, err
	}
	if err = cursor.All(ctx, &result); err != nil {
		return nil, err
	}

	return &result, nil
}
