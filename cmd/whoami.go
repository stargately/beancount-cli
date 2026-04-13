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

var whoamiCmd = &cobra.Command{
	Use:   "whoami",
	Short: "Print the currently authenticated user",
	RunE:  runWhoami,
}

func init() {
	rootCmd.AddCommand(whoamiCmd)
}

func runWhoami(cmd *cobra.Command, _ []string) error {
	creds, err := credentials.Load()
	if err != nil {
		return fmt.Errorf("failed to read credentials: %w", err)
	}
	if creds == nil || creds.Token == "" {
		return fmt.Errorf("not logged in. Run 'beancount-cli login' to authenticate")
	}
	if creds.IsExpired() {
		return fmt.Errorf("your session has expired. Run 'beancount-cli login' to re-authenticate")
	}

	cfg := config.Load()
	client := gqlclient.NewAuthed(cfg.APIURL, creds.Token)

	resp, err := generated.GetCurrentUser(context.Background(), client)
	if err != nil {
		if strings.Contains(err.Error(), "Not authenticated") {
			return fmt.Errorf("your session has expired. Run 'beancount-cli login' to re-authenticate")
		}
		return fmt.Errorf("cannot reach the Beancount API. Check your internet connection: %w", err)
	}

	profile := resp.UserProfile
	username := profile.Username
	if username == "" {
		username = "(not set)"
	}

	fmt.Fprintf(cmd.OutOrStdout(), "Email:    %s\n", profile.Email)
	fmt.Fprintf(cmd.OutOrStdout(), "Username: %s\n", username)
	fmt.Fprintf(cmd.OutOrStdout(), "Tier:     %s\n", profile.Tier)
	return nil
}
