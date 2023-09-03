package database

import (
	"container/list"
	"context"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type InMemoryDB struct {
	Users *list.List
}

func (db *InMemoryDB) CreateUser(ctx context.Context, tid string) (*User, error) {
	u := User{
		Id:       primitive.NewObjectID(),
		TwitchId: tid,
	}
	db.Users.PushBack(u)
	return &u, nil
}

func (db *InMemoryDB) GetUsers(ctx context.Context, u *User) (*[]User, error) {
	ul := list.New()
	for e := db.Users.Front(); e != nil; e = e.Next() {
		if e.Value.(User).TwitchId == u.TwitchId {
			ul.PushBack(e.Value)
		}
	}

	res := make([]User, ul.Len())
	i := 0
	for e := db.Users.Front(); e != nil; e = e.Next() {
		res[i] = e.Value.(User)
		i++
	}

	return &res, nil
}
