package auth

import (
	"context"
	"fmt"

	"github.com/spf13/cobra"

	"beancount.io/beancount-cli/generated"
	"beancount.io/beancount-cli/internal/config"
	"beancount.io/beancount-cli/internal/credentials"
	"beancount.io/beancount-cli/internal/gqlclient"
)

func NewCmdLogout() *cobra.Command {
	return &cobra.Command{
		Use:   "logout",
		Short: "Log out and revoke your local credentials",
		Long: `Revokes the stored token on the server and deletes ~/.beancount/credentials.json.

If the token has already expired server-side, the local credentials are still removed.`,
		RunE: runLogout,
	}
}

func runLogout(cmd *cobra.Command, _ []string) error {
	creds, err := credentials.Load()
	if err != nil {
		return fmt.Errorf("failed to read credentials: %w", err)
	}
	if creds == nil || creds.Token == "" {
		fmt.Fprintln(cmd.OutOrStdout(), "Not logged in")
		return nil
	}

	cfg := config.Load()
	client := gqlclient.NewAuthed(cfg.APIURL, creds.Token)

	// Revoke the token server-side; ignore any server error and always clear
	// local credentials so the client is never left in a broken state.
	_, _ = generated.Logout(context.Background(), client)

	if err := credentials.Clear(); err != nil {
		return fmt.Errorf("failed to clear credentials: %w", err)
	}

	fmt.Fprintln(cmd.OutOrStdout(), "Logged out")
	return nil
}
