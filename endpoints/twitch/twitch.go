package twitch

import (
	"auth/cache"
	"auth/database"
	"encoding/json"
	"log"
	"net/http"

	"github.com/golang-jwt/jwt"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type Twitch struct {
	Client TwitchAuth
	DB     database.UserRepo
	RDB    cache.TokenStore
	Secret string
}

type AuthRequest struct {
	Code  string
	State string
}

type AuthResponse struct {
	AccessToken  string
	RefreshToken string
}

func (t *Twitch) Authenticate(w http.ResponseWriter, r *http.Request) error {
	ctx := r.Context()

	var request AuthRequest
	json.NewDecoder(r.Body).Decode(&request)

	atr, err := t.Client.RequestUserAccessToken(request.Code)
	if err != nil {
		return err
	}
	if atr.Error != "" {
		w.WriteHeader(http.StatusUnauthorized)
		return nil
	}

	var userId string
	_, vr, _ := t.Client.ValidateToken(atr.Data.AccessToken)
	existingUsers, err := t.DB.GetUsers(ctx, &database.User{TwitchId: vr.Data.UserID})
	if err != nil {
		return err
	}

	// if user not exists, create one, publish event
	// and response with internal access and refresh tokens
	if len(*existingUsers) == 0 {
		user, err := t.DB.CreateUser(ctx, vr.Data.UserID)
		if err != nil {
			return err
		}
		// TODO: publish event

		userId = user.Id.String()
	} else if len(*existingUsers) == 1 {
		userId = (*existingUsers)[0].Id.String()
	} else {
		log.Printf("there is more than one user with same twitchId: %s \n", vr.Data.UserID)
		return err
	}

	at := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"userId": userId,
	})

	tokenString, err := at.SignedString([]byte(t.Secret))
	if err != nil {
		return err
	}

	refreshToken := primitive.NewObjectID().String()
	t.RDB.SaveTokens(ctx, tokenString, refreshToken)
	if err != nil {
		return err
	}
	respBytes, err := json.Marshal(AuthResponse{
		AccessToken:  tokenString,
		RefreshToken: refreshToken,
	})
	w.Write(respBytes)
	return nil
}
