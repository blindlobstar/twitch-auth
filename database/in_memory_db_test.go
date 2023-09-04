package database

import (
	"context"
	"testing"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

func TestInMemoryDb(t *testing.T) {
	db := InMemoryDB{
		Users: make([]User, 0),
	}
	u, _ := db.CreateUser(context.TODO(), "123")
	if u.Id == primitive.NilObjectID {
		t.Fatal("empty id")
	}

	users, _ := db.GetUsers(context.TODO(), &User{TwitchId: u.TwitchId})
	if len(*users) != 1 {
		t.Fatalf("wrong users count. expected: %d, got: %d", 1, len(*users))
	}
}
