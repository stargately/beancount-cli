package cmd

import (
	"context"
	"fmt"
	"strings"

	"github.com/spf13/cobra"

	"beancount.io/beancount-cli/generated"
	"beancount.io/beancount-cli/internal/config"
	"beancount.io/beancount-cli/internal/credentials"
	"beancount.io/beancount-cli/internal/gqlclient"
)

var logoutCmd = &cobra.Command{
	Use:   "logout",
	Short: "Log out and revoke your local credentials",
	Long: `Revokes the stored token on the server and deletes ~/.beancount/credentials.json.

If the token has already expired server-side, the local credentials are still removed.`,
	RunE: runLogout,
}

func init() {
	rootCmd.AddCommand(logoutCmd)
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

	// Revoke the token server-side. If the server rejects the token (already
	// expired), we still proceed to clear local credentials.
	_, err = generated.Logout(context.Background(), client)
	if err != nil && !isAuthError(err) {
		return fmt.Errorf("logout request failed: %w", err)
	}

	if err := credentials.Clear(); err != nil {
		return fmt.Errorf("failed to clear credentials: %w", err)
	}

	fmt.Fprintln(cmd.OutOrStdout(), "Logged out")
	return nil
}

// isAuthError reports whether err is an "Not authenticated" response from the server.
func isAuthError(err error) bool {
	return strings.Contains(err.Error(), "Not authenticated")
}
