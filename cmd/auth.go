package cmd

import (
	"fmt"

	"github.com/Khan/genqlient/graphql"

	"beancount.io/beancount-cli/internal/config"
	"beancount.io/beancount-cli/internal/credentials"
	"beancount.io/beancount-cli/internal/gqlclient"
)

// newAuthedClient validates the stored credentials and returns an authenticated
// GraphQL client. Returns an error if the user is not logged in or the session
// has expired.
func newAuthedClient() (graphql.Client, error) {
	creds, err := credentials.Load()
	if err != nil {
		return nil, fmt.Errorf("failed to read credentials: %w", err)
	}
	if creds == nil || creds.Token == "" {
		return nil, fmt.Errorf("not logged in. Run 'beancount-cli login' to authenticate")
	}
	if creds.IsExpired() {
		return nil, fmt.Errorf("your session has expired. Run 'beancount-cli login' to re-authenticate")
	}
	cfg := config.Load()
	return gqlclient.NewAuthed(cfg.APIURL, creds.Token), nil
}
