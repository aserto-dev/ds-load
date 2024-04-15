package azureclient

import (
	"context"
	"fmt"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	"github.com/Azure/azure-sdk-for-go/sdk/azcore"
	"github.com/Azure/azure-sdk-for-go/sdk/azidentity"
	msgraphsdk "github.com/aserto-dev/msgraph-sdk-go"
	adgroups "github.com/aserto-dev/msgraph-sdk-go/groups"
	"github.com/aserto-dev/msgraph-sdk-go/models"
	adusers "github.com/aserto-dev/msgraph-sdk-go/users"
	abstractions "github.com/microsoft/kiota-abstractions-go"
	graphcore "github.com/microsoftgraph/msgraph-sdk-go-core"

	auth "github.com/microsoft/kiota-authentication-azure-go"
	http "github.com/microsoft/kiota-http-go"
)

type AzureADClient struct {
	appClient      *msgraphsdk.Msgraph
	requestAdaptor abstractions.RequestAdapter
}

func NewAzureADClient(ctx context.Context, tenant, clientID, clientSecret string) (*AzureADClient, error) {
	c := &AzureADClient{}

	credential, err := azidentity.NewClientSecretCredential(tenant, clientID, clientSecret, nil)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "failed to create an Azure secret credential: %s", err.Error())
	}

	err = c.initClient(credential)
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

	err = c.initClient(credential)
	if err != nil {
		return nil, err
	}
	return c, nil
}

func (c *AzureADClient) ListUsers(ctx context.Context, groups bool) ([]models.Userable, error) {
	return c.listUsers(ctx, "", groups)
}

func (c *AzureADClient) GetUserByID(ctx context.Context, id string, groups bool) ([]models.Userable, error) {
	filter := fmt.Sprintf("id eq '%s'", id)
	return c.listUsers(ctx, filter, groups)
}

func (c *AzureADClient) GetUserByEmail(ctx context.Context, email string, groups bool) ([]models.Userable, error) {
	filter := fmt.Sprintf("mail eq '%s'", email)

	azureadUsers, err := c.listUsers(ctx, filter, groups)
	if err != nil {
		return azureadUsers, err
	}

	if len(azureadUsers) < 1 {
		filter := fmt.Sprintf("userPrincipalName eq '%s'", email)
		return c.listUsers(ctx, filter, groups)
	}
	return azureadUsers, err
}

func (c *AzureADClient) ListGroups(ctx context.Context) ([]models.Groupable, error) {
	result := make([]models.Groupable, 0)

	groupsResp, err := c.appClient.Groups().
		Get(ctx,
			&adgroups.GroupsRequestBuilderGetRequestConfiguration{
				QueryParameters: &adgroups.GroupsRequestBuilderGetQueryParameters{
					Select: []string{"displayName", "id", "mail", "createdDateTime", "mailNickname"},
					Expand: []string{"memberOf"},
				},
			})
	if err != nil {
		return nil, err
	}

	pageIterator, err := graphcore.NewPageIterator[*models.Group](
		groupsResp,
		c.requestAdaptor,
		models.CreateGroupCollectionResponseFromDiscriminatorValue)
	if err != nil {
		return nil, err
	}

	// Iterate over all pages
	err = pageIterator.Iterate(
		ctx,
		func(group *models.Group) bool {
			result = append(result, group)
			return true
		})
	if err != nil {
		return nil, err
	}

	return result, nil
}

func (c *AzureADClient) listUsers(ctx context.Context, filter string, groups bool) ([]models.Userable, error) {
	query := adusers.UsersRequestBuilderGetQueryParameters{
		Select: []string{"displayName", "id", "mail", "createdDateTime", "mobilePhone", "userPrincipalName", "accountEnabled"},
		Filter: &filter,
	}

	if groups {
		query.Expand = []string{"memberOf"}
	}

	result := make([]models.Userable, 0)

	usersResp, err := c.appClient.Users().
		Get(ctx,
			&adusers.UsersRequestBuilderGetRequestConfiguration{
				QueryParameters: &query,
			})
	if err != nil {
		return nil, err
	}

	pageIterator, err := graphcore.NewPageIterator[*models.User](
		usersResp,
		c.requestAdaptor,
		models.CreateUserCollectionResponseFromDiscriminatorValue)
	if err != nil {
		return nil, err
	}

	// Iterate over all pages
	err = pageIterator.Iterate(
		ctx,
		func(user *models.User) bool {
			result = append(result, user)
			return true
		})
	if err != nil {
		return nil, err
	}

	return result, nil
}

func (c *AzureADClient) initClient(credential azcore.TokenCredential) error {
	authProvider, err := auth.NewAzureIdentityAuthenticationProviderWithScopes(credential, []string{
		"https://graph.microsoft.com/.default",
	})
	if err != nil {
		return status.Errorf(codes.Internal, "failed to create Azure identity provider: %s", err.Error())
	}

	// Create a request adapter using the auth provider
	adapter, err := http.NewNetHttpRequestAdapter(authProvider)
	if err != nil {
		return status.Errorf(codes.Internal, "failed to create Azure AD Graph request adapter: %s", err.Error())
	}

	// Create a Graph client using request adapter
	c.appClient = msgraphsdk.NewMsgraph(adapter)
	c.requestAdaptor = adapter
	return nil
}
