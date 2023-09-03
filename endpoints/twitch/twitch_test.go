package twitch

import (
	"auth/cache"
	"auth/database"
	"auth/tokens"
	"container/list"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

func Test_TwitchAuth(t *testing.T) {
	rr := httptest.NewRecorder()

	rb := `{"code": "code1"}`
	req := httptest.NewRequest("GET", "/twitch", strings.NewReader(rb))

	ta := TwitchAuthMock{
		CodeTokenMap: map[string]string{
			"code": "228",
		},
		TokenUserMap: map[string]string{
			"228": "1",
		},
	}
	db := database.InMemoryDB{
		Users: list.New(),
	}
	te := &Twitch{
		Client: &ta,
		DB:     &db,
		AT:     tokens.New("secret"),
		RDB: &cache.InMemoryTokenStore{
			Tokens: make(map[string]string),
		},
	}

	// check unixisted code
	if err := te.Authenticate(rr, req); err != nil {
		t.Fatal(err)
	}

	if rr.Result().StatusCode != http.StatusUnauthorized {
		t.Fatal("status code is not unauthorized")
	}

	// now register code and check auth process
	rb = `{"code":"code"}`
	req = httptest.NewRequest("GET", "/twitch", strings.NewReader(rb))
	rr = httptest.NewRecorder()

	if err := te.Authenticate(rr, req); err != nil {
		t.Fatal(err)
	}

	if rr.Result().StatusCode != http.StatusOK {
		t.Fatalf("expected OK, got: %d", rr.Result().StatusCode)
	}
	var authResponse AuthResponse
	json.NewDecoder(rr.Body).Decode(&authResponse)
	if authResponse.AccessToken == "" {
		t.Fatal("empty access_token")
	}
	if db.Users.Len() != 1 {
		t.Fatal("user was not added")
	}

	// authenticate same user
	ta.CodeTokenMap = map[string]string{
		"new code": "228",
	}
	rb = `{"code": "new code"}`
	req = httptest.NewRequest("GET", "/twitch", strings.NewReader(rb))
	rr = httptest.NewRecorder()

	if err := te.Authenticate(rr, req); err != nil {
		t.Fatal(err)
	}

	if rr.Result().StatusCode != http.StatusOK {
		t.Fatalf("expected OK, got: %d", rr.Result().StatusCode)
	}
	json.NewDecoder(rr.Body).Decode(&authResponse)
	if authResponse.AccessToken == "" {
		t.Fatal("empty access_token")
	}
	if db.Users.Len() != 1 {
		t.Fatal("user was not added")
	}
}