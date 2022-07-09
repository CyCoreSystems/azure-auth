package smtp

import (
	"context"
	"fmt"
	"io"

	"github.com/CyCoreSystems/azure-auth/pkg/config"
	"github.com/CyCoreSystems/azure-auth/pkg/token"
	"github.com/emersion/go-smtp"
)

// Send sends an email.
func Send(ctx context.Context, cfg *config.Config, svc string, recipients []string, body io.Reader) error {
	if len(recipients) < 1 {
		return fmt.Errorf("no recipients")
	}

	tok, err := token.Get(ctx, cfg)
	if err != nil {
		return fmt.Errorf("failed to get token: %w", err)
	}

	auth := NewOAuthBearerClient(&OAuthBearerOptions{
		Username:  cfg.Username,
		AccessToken: tok.AccessToken,
	})


	if err := smtp.SendMail(svc, auth, cfg.Username, recipients, body); err != nil {
		return fmt.Errorf("failed to send message: %w", err)
	}

	return nil
}
