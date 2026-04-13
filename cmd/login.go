package cmd

import (
	"context"
	"fmt"
	"time"

	"github.com/Khan/genqlient/graphql"
	"github.com/pkg/browser"
	"github.com/spf13/cobra"

	"beancount.io/beancount-cli/generated"
	"beancount.io/beancount-cli/internal/config"
	"beancount.io/beancount-cli/internal/credentials"
	"beancount.io/beancount-cli/internal/gqlclient"
)

var loginCmd = &cobra.Command{
	Use:   "login",
	Short: "Log in to Beancount via your browser",
	Long: `Initiates a browser-based device authorization flow.

The CLI creates a session, opens your default browser to the Beancount
authorization page, then polls until you approve or deny the request.
On success, the token is saved to ~/.beancount/credentials.json.`,
	RunE: runLogin,
}

func init() {
	rootCmd.AddCommand(loginCmd)
}

func runLogin(cmd *cobra.Command, _ []string) error {
	cfg := config.Load()
	client := gqlclient.New(cfg.APIURL)
	ctx := context.Background()

	// Step 1: create a session on the backend.
	resp, err := generated.CreateCliAuthSession(ctx, client)
	if err != nil {
		return fmt.Errorf("cannot reach the Beancount API. Check your internet connection: %w", err)
	}
	sessionID := resp.CreateCliAuthSession.SessionId

	// Step 2: open the browser to the device auth page.
	authURL := fmt.Sprintf("%s/auth/login/device?session_id=%s", cfg.DashboardURL, sessionID)
	fmt.Fprintln(cmd.OutOrStdout(), "Opening browser to authorize...")
	if err := browser.OpenURL(authURL); err != nil {
		fmt.Fprintln(cmd.OutOrStdout(), "Could not open browser automatically.")
	}
	fmt.Fprintf(cmd.OutOrStdout(), "If the browser did not open, visit:\n  %s\n", authURL)

	// Step 3: poll every 2 seconds until the session resolves.
	return pollForAuth(ctx, cmd, client, sessionID)
}

const pollInterval = 2 * time.Second

func pollForAuth(ctx context.Context, cmd *cobra.Command, client graphql.Client, sessionID string) error {
	for {
		time.Sleep(pollInterval)

		resp, err := generated.GetCliAuthSession(ctx, client, sessionID)
		if err != nil {
			return fmt.Errorf("cannot reach the Beancount API. Check your internet connection: %w", err)
		}

		session := resp.GetCliAuthSession
		switch session.Status {
		case generated.CliAuthStatusAuthorized:
			// Step 4: consume the token (single-use — clears it from the server).
			consumeResp, err := generated.ConsumeCliAuthSession(ctx, client, sessionID)
			if err != nil {
				return fmt.Errorf("failed to retrieve auth token: %w", err)
			}
			consumed := consumeResp.ConsumeCliAuthSession
			if err := credentials.Save(consumed.Token, consumed.ExpireAt); err != nil {
				return fmt.Errorf("failed to save credentials: %w", err)
			}
			fmt.Fprintln(cmd.OutOrStdout(), "Logged in successfully")
			return nil

		case generated.CliAuthStatusDenied:
			return fmt.Errorf("authorization denied")

		case generated.CliAuthStatusExpired:
			return fmt.Errorf("session expired. Run 'beancount-cli login' to try again")

		default:
			// CliAuthStatusPending — keep waiting.
		}
	}
}
