package cache

import "context"

type InMemoryTokenStore struct {
	Tokens map[string]string
}

func (ts *InMemoryTokenStore) SaveTokens(ctx context.Context, at string, rt string) error {
	ts.Tokens[rt] = at
	return nil
}

func (ts *InMemoryTokenStore) GetToken(ctx context.Context, rt string) (string, error) {
	return ts.Tokens[rt], nil
}
