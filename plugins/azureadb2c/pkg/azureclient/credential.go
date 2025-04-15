package azureclient

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strings"
	"time"

	"github.com/Azure/azure-sdk-for-go/sdk/azcore"
	"github.com/Azure/azure-sdk-for-go/sdk/azcore/policy"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type RefreshTokenCredential struct {
	clientID     string
	clientSecret string
	refreshToken string
	tenantID     string
}

func NewRefreshTokenCredential(ctx context.Context,
	tenantID, clientID, clientSecret, refreshToken string,
) (*RefreshTokenCredential, error) {
	c := &RefreshTokenCredential{
		clientID:     clientID,
		clientSecret: clientSecret,
		tenantID:     tenantID,
		refreshToken: refreshToken,
	}

	return c, nil
}

func (c *RefreshTokenCredential) GetToken(ctx context.Context, options policy.TokenRequestOptions) (azcore.AccessToken, error) {
	accessToken := azcore.AccessToken{}

	url := "https://login.microsoftonline.com/" + c.tenantID + "/oauth2/v2.0/token"
	data := fmt.Sprintf("grant_type=refresh_token&client_id=%s&client_secret=%s&refresh_token=%s",
		c.clientID, c.clientSecret, c.refreshToken)
	payload := strings.NewReader(data)

	// create the request and execute it
	req, err := http.NewRequestWithContext(ctx, http.MethodPost, url, payload)
	if err != nil {
		return accessToken, err
	}

	req.Header.Add("Content-Type", "application/x-www-form-urlencoded")

	res, err := http.DefaultClient.Do(req)
	if err != nil {
		return accessToken, err
	}

	// process the response
	defer res.Body.Close()

	var responseData map[string]interface{}

	body, _ := io.ReadAll(res.Body)

	// unmarshal the json into a string map
	if err := json.Unmarshal(body, &responseData); err != nil {
		return accessToken, err
	}

	// check for error
	if responseData["error_description"] != nil {
		errorMessage, ok := responseData["error_description"].(string)
		if !ok {
			return accessToken, status.Error(codes.InvalidArgument, "failed to get error description")
		}

		return accessToken, status.Error(codes.InvalidArgument, errorMessage)
	}

	// retrieve the access token and expiration
	token, ok := responseData["access_token"].(string)
	if !ok {
		return accessToken, status.Error(codes.InvalidArgument, "failed to cast access token to string")
	}

	accessToken.Token = token

	expiresIn, ok := responseData["expires_in"].(float64)
	if !ok {
		return accessToken, status.Error(codes.InvalidArgument, "failed to convert token expiration time")
	}

	accessToken.ExpiresOn = time.Now().Add(time.Second * time.Duration(int(expiresIn)))

	return accessToken, nil
}
