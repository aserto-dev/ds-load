package googleclient

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
	"google.golang.org/api/option"
)

type GoogleClient struct {
	googleClient *admin.Service
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

	// Create an HTTP server for handling the OAuth 2.0 callback
	server := &http.Server{Addr: fmt.Sprintf(":%d", port), ReadHeaderTimeout: 5 * time.Second}

	// Create a channel to receive the authorization code
	authCodeChan := make(chan string)

	// Start the HTTP server in a separate goroutine
	go func() {
		http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
			code := r.URL.Query().Get("code")
			fmt.Fprintf(w, "Authorization code received. You can close this tab now.")
			authCodeChan <- code
		})

		if err := server.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
			log.Fatalf("Failed to start HTTP server: %v", err)
		}
	}()

	// Generate the authorization URL with "access_type=offline" and "prompt=consent"
	authURL := config.AuthCodeURL("state-token", oauth2.AccessTypeOffline, oauth2.SetAuthURLParam("prompt", "consent"))

	fmt.Printf("Go to the following URL to authorize the application:\n\n%s\n\n", authURL)
	fmt.Println("Waiting for authorization code...")

	// Receive the authorization code from the channel, shut down the server
	authCode := <-authCodeChan

	go func() { server.Shutdown(ctx) }()

	// Exchange the authorization code for an access token
	token, err := config.Exchange(ctx, authCode)
	if err != nil {
		log.Printf("Failed to exchange authorization code for access token: %v\n", err)
		return "", err
	} else {
		return token.RefreshToken, nil
	}
}

func NewGoogleClient(ctx context.Context, clientID, clientSecret, refreshToken string) (*GoogleClient, error) {
	c := &GoogleClient{}

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

	svc, err := admin.NewService(ctx, option.WithTokenSource(config.TokenSource(ctx, token)))
	if err != nil {
		return nil, err
	}

	c.googleClient = svc
	return c, nil
}

func (c *GoogleClient) ListUsers() ([]*admin.User, error) {
	users := make([]*admin.User, 0)
	pageToken := ""

	for {
		response, err := c.googleClient.Users.List().Customer("my_customer").PageToken(pageToken).Do()
		if err != nil {
			return nil, err
		}
		users = append(users, response.Users...)

		if response.NextPageToken == "" {
			break
		}

		pageToken = response.NextPageToken
	}

	return users, nil
}

func (c *GoogleClient) ListGroups() ([]*admin.Group, error) {
	groups := make([]*admin.Group, 0)
	pageToken := ""

	for {
		response, err := c.googleClient.Groups.List().Customer("my_customer").PageToken(pageToken).Do()
		if err != nil {
			return nil, err
		}
		groups = append(groups, response.Groups...)

		if response.NextPageToken == "" {
			break
		}

		pageToken = response.NextPageToken
	}

	return groups, nil
}

func (c *GoogleClient) GetUsersInGroup(group string) ([]*admin.Member, error) {
	members := make([]*admin.Member, 0)
	pageToken := ""

	for {
		response, err := c.googleClient.Members.List(group).PageToken(pageToken).Do()
		if err != nil {
			return nil, err
		}
		members = append(members, response.Members...)

		if response.NextPageToken == "" {
			break
		}

		pageToken = response.NextPageToken
	}

	return members, nil
}

func (c *GoogleClient) GetUserByID(id string) (*admin.User, error) {
	return c.listUser(id)
}

func (c *GoogleClient) GetUserByEmail(email string) (*admin.User, error) {
	return c.listUser(email)
}

func (c *GoogleClient) listUser(user string) (*admin.User, error) {
	return c.googleClient.Users.Get(user).Do()
}
