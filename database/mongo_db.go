package database

import "go.mongodb.org/mongo-driver/mongo"

type MongoDB struct {
	DB *mongo.Database
}

func Create(c *mongo.Client) MongoDB {
	db := c.Database("auth")
	return MongoDB{DB: db}
}
