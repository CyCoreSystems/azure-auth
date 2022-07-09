package token

import (
	"context"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/exec"
	"path"
	"strings"

	"github.com/CyCoreSystems/azure-auth/pkg/config"

	"golang.org/x/oauth2"
)

func tokenFileName() string {
	cfgDir, err := os.UserConfigDir()
	if err != nil {
		log.Fatalln("failed to determine configuration directory:", err)
	}

	return path.Join(cfgDir, "azure", "token.json")
}

// Get returns a Token, delegating the to browser, if needed.
func Get(pCtx context.Context, cfg *config.Config) (*oauth2.Token,error) {
	ctx, cancel := context.WithCancel(pCtx)
	defer cancel()

	ts := oauth2.ReuseTokenSource(loadFileToken(), &BrowserTokenSource{
		ctx:    ctx,
		cancel: cancel,
		cfg: cfg,
	})

	return ts.Token()
}

// XOAuth2Token returns an XOAuth2 token formatted for use in SMTP and IMAP settings.
func XOAuth2Token(ctx context.Context, cfg *config.Config) (string, error) {
	tok, err := Get(ctx, cfg)
	if err != nil {
		return "", fmt.Errorf("failed to retrieve token: %w", err)
	}

	return fmt.Sprintf("user=%s\x01auth=Bearer %s\x01\x01", cfg.Username, tok.AccessToken), nil
}

func exporter(cfg *config.Config) {
	var ts oauth2.TokenSource

	token, err := ts.Token()
	if err != nil {
		log.Fatalln("failed to retrieve token:", err)
	}

	xoauth2 := fmt.Sprintf("user=%s\x01auth=Bearer %s\x01\x01", cfg.Username, token.AccessToken)

	out := struct {
		AccessToken  string `json:"accessToken"`
		RefreshToken string `json:"refreshToken"`
		XOAuth2      string `json:"xoauth2"`
	}{
		AccessToken:  token.AccessToken,
		RefreshToken: token.RefreshToken,
		XOAuth2:      base64.StdEncoding.EncodeToString([]byte(xoauth2)),
	}

	if err := json.NewEncoder(os.Stdout).Encode(&out); err != nil {
		log.Fatalln("failed to encode output:", err)
	}

	return
}

func loadFileToken() *oauth2.Token {

	tokenFileData, err := os.ReadFile(tokenFileName())
	if err != nil {
		return nil
	}

	token := new(oauth2.Token)

	if err := json.Unmarshal(tokenFileData, token); err != nil {
		log.Printf("failed to parse token from file %q: %s", tokenFileName(), err.Error())
		return nil
	}

	return token
}

type BrowserTokenSource struct {
	ctx    context.Context
	cancel context.CancelFunc

	cfg *config.Config

	token *oauth2.Token
}

// Token retrieves a token from a browser interaction.
func (s *BrowserTokenSource) Token() (*oauth2.Token, error) {
	srv := &http.Server{
		Addr: fmt.Sprintf(":%d", s.cfg.Redirect.Port),
	}

	go func() {
		http.HandleFunc(fmt.Sprintf("/%s", strings.Trim(s.cfg.Redirect.Path,"/")), s.redirectHandler())

		if err := srv.ListenAndServe(); err != http.ErrServerClosed {
			log.Fatalln("HTTP server died unexpectedly:", err)
		}
	}()

	if err := exec.Command("xdg-open", s.cfg.OAuth2().AuthCodeURL("state", oauth2.AccessTypeOffline)).Run(); err != nil {
		fmt.Println("failed to launch browser automatically:", err)
		fmt.Println()
		fmt.Println("1. Ensure that you are logged in as your user in your browser.")
		fmt.Println()
		fmt.Println("2. Open the following link and authorise XOAUTH2 token:")
		fmt.Println(s.cfg.OAuth2().AuthCodeURL("state", oauth2.AccessTypeOffline))
		fmt.Println()
	}

	<-s.ctx.Done()

	srv.Shutdown(s.ctx)

	return s.token, nil
}

func (s *BrowserTokenSource) redirectHandler() func(http.ResponseWriter, *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		code := r.FormValue("code")

		token, err := s.cfg.OAuth2().Exchange(s.ctx, code)
		if err != nil {
			log.Fatalf("Failed to exchange authorisation code for token: %v.", err)
		}

		s.token = token

		// Store new token to file
		if err := s.storeToken(); err != nil {
			log.Println("failed to store updated token:", err)
		}

		defer s.cancel()

		if _, err := w.Write([]byte("OK")); err != nil {
			log.Println("failed to write response:", err)
		}
	}
}

func (s *BrowserTokenSource) storeToken() error {
	if s.token == nil {
		return fmt.Errorf("no token")
	}

	data, err := json.Marshal(s.token)
	if err != nil {
		return fmt.Errorf("failed to encode token to JSON: %w", err)
	}

	if err := os.WriteFile(tokenFileName(), data, os.FileMode(0600)); err != nil {
		return fmt.Errorf("failed to write token file %q: %w", tokenFileName(), err)
	}

	return nil
}
