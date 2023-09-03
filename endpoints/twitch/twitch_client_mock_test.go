package twitch

import (
	"net/http"
	"testing"
)

func Test_TwitchAuthMock(t *testing.T) {
	tm := TwitchAuthMock{
		CodeTokenMap: make(map[string]string),
		TokenUserMap: make(map[string]string),
	}

	code := "code"
	tm.RegisterCode(code)
	resp, _ := tm.RequestUserAccessToken(code)
	if resp.StatusCode != http.StatusOK || resp.Data.AccessToken == "" {
		t.Fatalf("error")
	}

	// try get access token with same code
	resp, _ = tm.RequestUserAccessToken(code)
	if resp.StatusCode != http.StatusBadRequest || resp.Data.AccessToken != "" {
		t.Fatalf("error 2")
	}
}
