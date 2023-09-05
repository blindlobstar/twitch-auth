package main

import (
	"auth/cache"
	"auth/database"
	"auth/endpoints/twitch"
	"context"
	"log"
	"net/http"
	"os"

	"github.com/gorilla/mux"
	"github.com/joho/godotenv"
	"github.com/nicklaw5/helix"
	"github.com/redis/go-redis/v9"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func main() {
	ctx := context.Background()

	//loading env variables
	godotenv.Load()

	// configuring mongodb
	log.Println("connecting db..")
	dbUri := os.Getenv("MONGO_URI")
	client, err := mongo.Connect(ctx, options.Client().ApplyURI(dbUri))
	if err != nil {
		log.Fatal(err)
	}
	defer func() {
		client.Disconnect(context.TODO())
	}()
	db := database.Create(client)

	// configuring redis
	log.Println("connecting rdb..")
	rdbAddr := os.Getenv("REDIS_URI")
	rdbClient := redis.NewClient(&redis.Options{
		Addr:     rdbAddr,
		Password: "",
		DB:       0,
	})
	rdb := cache.RedisTokenStore{RDB: rdbClient}

	// configuring twitch client
	log.Println("configuring twitch client..")
	twitchClientId := os.Getenv("TWITCH_CLIENT_ID")
	twitchSecret := os.Getenv("TWITCH_SECRET")
	twitchClient, err := helix.NewClient(&helix.Options{
		ClientID:     twitchClientId,
		ClientSecret: twitchSecret,
	})
	if err != nil {
		log.Fatal("error initializing twitch client")
	}

	r := mux.NewRouter()

	tokenSecret := os.Getenv("TOKEN_SECRET")
	twitch := twitch.Twitch{
		Client: twitchClient,
		DB:     &db,
		RDB:    &rdb,
		Secret: tokenSecret,
	}

	// register handlers
	r.HandleFunc("/twitch", errorHandlingMiddleware(twitch.Authenticate)).Methods("POST")

	log.Println("server starting..")
	http.ListenAndServe(":80", r)
}

func errorHandlingMiddleware(f func(w http.ResponseWriter, r *http.Request) error) func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		if err := f(w, r); err != nil {
			log.Print(err)
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
	}
}
