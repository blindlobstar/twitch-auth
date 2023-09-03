package main

import (
	"auth/endpoints/twitch"
	"auth/tokens"
)

func main() {
	t := twitch.Twitch{}
	t.AT = &tokens.AccessToken{}
}
