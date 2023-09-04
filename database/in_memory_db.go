package database

import (
	"context"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type InMemoryDB struct {
	Users []User
}

func (db *InMemoryDB) CreateUser(ctx context.Context, tid string) (*User, error) {
	u := User{
		Id:       primitive.NewObjectID(),
		TwitchId: tid,
	}
	db.Users = append(db.Users, u)
	return &u, nil
}

func (db *InMemoryDB) GetUsers(ctx context.Context, u *User) (*[]User, error) {
	res := make([]User, len(db.Users))
	var l int
	for i, u := range db.Users {
		if db.Users[i].TwitchId == u.TwitchId {
			res[l] = u
			l++
		}
	}
	res = res[:l]
	return &res, nil
}
