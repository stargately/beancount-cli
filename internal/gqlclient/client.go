package gqlclient

import (
	"net/http"

	"github.com/Khan/genqlient/graphql"
)

// New returns an unauthenticated GraphQL client for the given endpoint.
func New(apiURL string) graphql.Client {
	return graphql.NewClient(apiURL, http.DefaultClient)
}

// NewAuthed returns a GraphQL client that attaches a Bearer token to every request.
func NewAuthed(apiURL, token string) graphql.Client {
	httpClient := &http.Client{
		Transport: &bearerTransport{token: token},
	}
	return graphql.NewClient(apiURL, httpClient)
}

type bearerTransport struct {
	token string
}

func (t *bearerTransport) RoundTrip(req *http.Request) (*http.Response, error) {
	req = req.Clone(req.Context())
	req.Header.Set("Authorization", "Bearer "+t.token)
	return http.DefaultTransport.RoundTrip(req)
}
