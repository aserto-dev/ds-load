package hubspotclient

import (
	"context"
	"errors"
	"fmt"
	"log"
	"net/http"
	"time"

	"golang.org/x/oauth2"
	"golang.org/x/oauth2/google"
	admin "google.golang.org/api/admin/directory/v1"
)

type HubspotClient struct {
	token     string
	users     map[string]string
	companies map[string]Company
}

func GetRefreshToken(ctx context.Context, clientID, clientSecret string, port int) (string, error) {
	redirectURL := fmt.Sprintf("http://localhost:%d", port)
	config := &oauth2.Config{
		ClientID:     clientID,
		ClientSecret: clientSecret,
		RedirectURL:  redirectURL,
		Scopes: []string{
			admin.AdminDirectoryUserScope,
			admin.AdminDirectoryGroupScope,
		},
		Endpoint: google.Endpoint,
	}

	var authCode string

	// Generate the authorization URL with "access_type=offline" and "prompt=consent"
	authURL := config.AuthCodeURL("state-token", oauth2.AccessTypeOffline, oauth2.SetAuthURLParam("prompt", "consent"))

	fmt.Printf("Go to the following URL to authorize the application:\n\n%s\n\n", authURL)
	fmt.Println("Waiting for authorization code...")

	// Create an HTTP server for handling the OAuth 2.0 callback
	server := &http.Server{Addr: fmt.Sprintf(":%d", port), ReadHeaderTimeout: 5 * time.Second}
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		code := r.URL.Query().Get("code")
		fmt.Fprintf(w, "Authorization code received. You can close this tab now.")
		authCode = code
		go func() {
			// Shutdown the HTTP server once the callback is received
			if err := server.Shutdown(ctx); err != nil {
				log.Printf("Failed to shutdown HTTP server: %s", err)
			}
		}()
	})

	// Start the HTTP server.
	if err := server.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
		log.Fatalf("Failed to start HTTP server: %v", err)
	}

	// Exchange the authorization code for an access token
	token, err := config.Exchange(ctx, authCode)
	if err != nil {
		log.Printf("Failed to exchange authorization code for access token: %v\n", err)
		return "", err
	} else {
		return token.RefreshToken, nil
	}
}

func NewHubspotClient(ctx context.Context, privateAccessToken string) (*HubspotClient, error) {
	c := &HubspotClient{}
	c.token = privateAccessToken
	c.users = make(map[string]string, 0)
	c.companies = make(map[string]Company, 0)
	return c, nil
}

func NewHubspotOAuth2Client(ctx context.Context, clientID, clientSecret, refreshToken string) (*HubspotClient, error) {
	c := &HubspotClient{}

	config := &oauth2.Config{
		ClientID:     clientID,
		ClientSecret: clientSecret,
		Scopes: []string{
			admin.AdminDirectoryUserScope,
			admin.AdminDirectoryGroupScope,
		},
		Endpoint: google.Endpoint,
	}

	token := &oauth2.Token{
		RefreshToken: refreshToken,
	}

	tokenSource := config.TokenSource(ctx, token)

	accessToken, err := tokenSource.Token()
	if err != nil {
		return nil, err
	}

	c.token = accessToken.AccessToken
	c.users = make(map[string]string, 0)
	c.companies = make(map[string]Company, 0)
	return c, nil
}
