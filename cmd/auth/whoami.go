package auth

import (
	"context"
	"fmt"
	"strings"

	"github.com/spf13/cobra"

	"beancount.io/beancount-cli/generated"
	"beancount.io/beancount-cli/internal/utils"
)

func NewCmdWhoami() *cobra.Command {
	return &cobra.Command{
		Use:   "whoami",
		Short: "Print the currently authenticated user",
		RunE:  runWhoami,
	}
}

func runWhoami(cmd *cobra.Command, _ []string) error {
	client, err := utils.NewAuthedClient()
	if err != nil {
		return err
	}

	resp, err := generated.GetCurrentUser(context.Background(), client)
	if err != nil {
		if strings.Contains(err.Error(), "Not authenticated") {
			return fmt.Errorf("your session has expired. Run 'beancount-cli login' to re-authenticate")
		}
		return fmt.Errorf("cannot reach the Beancount API. Check your internet connection: %w", err)
	}

	profile := resp.UserProfile
	username := "(not set)"
	if profile.Username != nil && *profile.Username != "" {
		username = *profile.Username
	}

	fmt.Fprintf(cmd.OutOrStdout(), "Email:    %s\n", profile.Email)
	fmt.Fprintf(cmd.OutOrStdout(), "Username: %s\n", username)
	fmt.Fprintf(cmd.OutOrStdout(), "Tier:     %s\n", profile.Tier)
	return nil
}
