package azureclient

import (
	"context"
	"fmt"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	"github.com/Azure/azure-sdk-for-go/sdk/azcore"
	"github.com/Azure/azure-sdk-for-go/sdk/azidentity"
	msgraphsdk "github.com/aserto-dev/msgraph-sdk-go"
	"github.com/aserto-dev/msgraph-sdk-go/models"
	adusers "github.com/aserto-dev/msgraph-sdk-go/users"
	auth "github.com/microsoft/kiota-authentication-azure-go"
	http "github.com/microsoft/kiota-http-go"
)

type AzureADClient struct {
	appClient *msgraphsdk.Msgraph
}

func NewAzureADClient(ctx context.Context, tenant, clientID, clientSecret string) (*AzureADClient, error) {
	c := &AzureADClient{}

	credential, err := azidentity.NewClientSecretCredential(tenant, clientID, clientSecret, nil)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "failed to create an Azure secret credential: %s", err.Error())
	}

	c.appClient, err = getAppClient(credential)
	if err != nil {
		return nil, err
	}
	return c, nil
}

func NewAzureADClientWithRefreshToken(ctx context.Context, tenant, clientID, clientSecret, refreshToken string) (*AzureADClient, error) {
	c := &AzureADClient{}

	credential, err := NewRefreshTokenCredential(ctx, tenant, clientID, clientSecret, refreshToken)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "failed to create Refresh Token credential: %s", err.Error())
	}

	c.appClient, err = getAppClient(credential)
	if err != nil {
		return nil, err
	}
	return c, nil
}

func (c *AzureADClient) ListUsers(ctx context.Context) (models.UserCollectionResponseable, error) {
	return c.listUsers(ctx, "")
}

func (c *AzureADClient) GetUserByID(ctx context.Context, id string) (models.UserCollectionResponseable, error) {
	filter := fmt.Sprintf("id eq '%s'", id)
	return c.listUsers(ctx, filter)
}

func (c *AzureADClient) GetUserByEmail(ctx context.Context, email string) (models.UserCollectionResponseable, error) {
	filter := fmt.Sprintf("mail eq '%s'", email)

	aadUsers, err := c.listUsers(ctx, filter)
	if err != nil {
		return aadUsers, err
	}

	azureadUsers := aadUsers.GetValue()
	if len(azureadUsers) < 1 {
		filter := fmt.Sprintf("userPrincipalName eq '%s'", email)
		return c.listUsers(ctx, filter)
	}
	return aadUsers, err
}

func (c *AzureADClient) listUsers(ctx context.Context, filter string) (models.UserCollectionResponseable, error) {
	query := adusers.UsersRequestBuilderGetQueryParameters{
		Select: []string{"displayName", "id", "mail", "createdDateTime", "mobilePhone", "userPrincipalName", "accountEnabled"},
		Filter: &filter,
	}
	return c.appClient.Users().
		Get(ctx,
			&adusers.UsersRequestBuilderGetRequestConfiguration{
				QueryParameters: &query,
			})
}

func getAppClient(credential azcore.TokenCredential) (*msgraphsdk.Msgraph, error) {
	authProvider, err := auth.NewAzureIdentityAuthenticationProviderWithScopes(credential, []string{
		"https://graph.microsoft.com/.default",
	})
	if err != nil {
		return nil, status.Errorf(codes.Internal, "failed to create Azure identity provider: %s", err.Error())
	}

	// Create a request adapter using the auth provider
	adapter, err := http.NewNetHttpRequestAdapter(authProvider)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "failed to create Azure AD Graph request adapter: %s", err.Error())
	}

	// Create a Graph client using request adapter
	client := msgraphsdk.NewMsgraph(adapter)
	return client, nil
}
