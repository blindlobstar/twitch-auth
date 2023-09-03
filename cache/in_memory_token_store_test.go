package cache

import (
	"context"
	"testing"
)

func TestInMemoryTokenStore(t *testing.T) {
	ts := InMemoryTokenStore{
		Tokens: make(map[string]string),
	}

	ts.SaveTokens(context.TODO(), "abc", "1")
	at, _ := ts.GetToken(context.TODO(), "1")
	if at == "" {
		t.Fatal("expected token was not found")
	}

	uat, _ := ts.GetToken(context.TODO(), "2")
	if uat != "" {
		t.Fatal("not existed token was found")
	}
}
