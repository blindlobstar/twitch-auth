package twitch

import (
	"fmt"
	"math/rand"
	"net/http"

	"github.com/nicklaw5/helix"
)

type TwitchAuth interface {
	RequestUserAccessToken(code string) (*helix.UserAccessTokenResponse, error)
	ValidateToken(accessToken string) (bool, *helix.ValidateTokenResponse, error)
}

type TwitchAuthMock struct {
	CodeTokenMap map[string]string
	TokenUserMap map[string]string
}

func (t *TwitchAuthMock) RegisterCode(code string) {
	at := fmt.Sprint(rand.Intn(1000))
	uid := fmt.Sprint(rand.Intn(1000))
	t.CodeTokenMap[code] = at
	t.TokenUserMap[at] = uid
}

func (t *TwitchAuthMock) RequestUserAccessToken(code string) (*helix.UserAccessTokenResponse, error) {
	var resp helix.UserAccessTokenResponse
	at := t.CodeTokenMap[code]
	if at == "" {
		resp.StatusCode = http.StatusBadRequest
		resp.ErrorStatus = http.StatusBadRequest
		resp.ErrorMessage = "wrong code"
		resp.Error = "wrong code"
		return &resp, nil
	}

	t.CodeTokenMap[code] = ""
	resp.StatusCode = http.StatusOK
	resp.Data = helix.AccessCredentials{
		AccessToken: at,
	}
	return &resp, nil
}

func (t *TwitchAuthMock) ValidateToken(accessToken string) (bool, *helix.ValidateTokenResponse, error) {
	var resp helix.ValidateTokenResponse
	uid := t.TokenUserMap[accessToken]
	if uid == "" {
		resp.StatusCode = http.StatusBadRequest
		resp.ErrorStatus = http.StatusBadRequest
		resp.ErrorMessage = "wrong access token"
		resp.Error = "wrong access token"
		return false, &resp, nil
	}
	resp.StatusCode = http.StatusOK
	resp.Data.UserID = uid
	return true, &resp, nil
}
