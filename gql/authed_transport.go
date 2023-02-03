package gql

import (
	"net/http"
)

type AuthedTransport struct {
	apiKey  string
	wrapped http.RoundTripper
}

func NewAuthedTransport(apiKey string) *AuthedTransport {
	return &AuthedTransport{apiKey: apiKey, wrapped: http.DefaultTransport}
}

func (t *AuthedTransport) RoundTrip(req *http.Request) (*http.Response, error) {
	req.Header.Set("Authorization", "bearer "+t.apiKey)
	return t.wrapped.RoundTrip(req)
}
