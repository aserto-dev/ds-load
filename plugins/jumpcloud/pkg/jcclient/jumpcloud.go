package jcclient

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

type JumpCloudClient struct {
	JumpCloudClient *admin.Service
	customer        string
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

func NewJumpCloudClient(ctx context.Context, clientID, clientSecret, refreshToken, customer string) (*JumpCloudClient, error) {
	c := &JumpCloudClient{}

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

	c.JumpCloudClient = svc
	c.customer = customer
	return c, nil
}

func (c *JumpCloudClient) ListUsers() ([]*admin.User, error) {
	users := make([]*admin.User, 0)
	pageToken := ""

	for {
		response, err := c.JumpCloudClient.Users.List().Customer(c.customer).PageToken(pageToken).Do()
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

func (c *JumpCloudClient) ListGroups() ([]*admin.Group, error) {
	groups := make([]*admin.Group, 0)
	pageToken := ""

	for {
		response, err := c.JumpCloudClient.Groups.List().Customer(c.customer).PageToken(pageToken).Do()
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

func (c *JumpCloudClient) GetUsersInGroup(group string) ([]*admin.Member, error) {
	members := make([]*admin.Member, 0)
	pageToken := ""

	for {
		response, err := c.JumpCloudClient.Members.List(group).PageToken(pageToken).Do()
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

func (c *JumpCloudClient) GetUserByID(id string) (*admin.User, error) {
	return c.listUser(id)
}

func (c *JumpCloudClient) GetUserByEmail(email string) (*admin.User, error) {
	return c.listUser(email)
}

func (c *JumpCloudClient) listUser(user string) (*admin.User, error) {
	return c.JumpCloudClient.Users.Get(user).Do()
}
