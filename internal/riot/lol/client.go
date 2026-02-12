package lol

import (
	"net/http"
)

type LolClient struct {
	apiKey string
	client *http.Client
}

func NewClient(apiKey string) *LolClient {
	return &LolClient{
		apiKey: apiKey,
		client: &http.Client{},
	}
}
