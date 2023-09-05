package main

import (
	"auth/cache"
	"auth/database"
	"auth/endpoints/twitch"
	"context"
	"log"
	"net/http"
	"os"
	"os/signal"
	"time"

	"github.com/gorilla/mux"
	"github.com/joho/godotenv"
	"github.com/nicklaw5/helix"
	"github.com/redis/go-redis/v9"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func main() {
	ctx := context.Background()

	if err := godotenv.Load(); err != nil {
		log.Print("failed load .env file")
	}

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

	sig := make(chan os.Signal, 1)
	signal.Notify(sig, os.Interrupt)

	log.Println("server starting..")
	s := http.Server{
		Addr:    ":80",
		Handler: r,
	}

	go func() {
		if err := s.ListenAndServe(); err != nil {
			log.Fatal(err)
		}
	}()
	log.Println("server has started!")

	<-sig
	log.Println("stopping server")
	ctx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()

	if err := s.Shutdown(ctx); err != nil {
		log.Fatal(err)
	}
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
